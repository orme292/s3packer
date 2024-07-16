package provider_v2

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/tuipack"
)

type worker struct {
	uuid string

	path       string
	searchRoot string
	isDir      bool
	isFile     bool

	stats *Stats

	app *conf.AppConfig

	status int

	oper  Operator
	objFn ObjectGenFunc
}

func newWorker(app *conf.AppConfig, path, searchRoot string, d, f bool, stat int, oper Operator, objFn ObjectGenFunc) *worker {

	if d == f {
		d = true
	}

	return &worker{
		uuid:       uuid.New().String(),
		path:       path,
		searchRoot: searchRoot,
		isDir:      d,
		isFile:     f,
		app:        app,
		status:     stat,
		oper:       oper,
		objFn:      objFn,
		stats:      &Stats{},
	}

}

func (w *worker) scan() {

	if w.status == JobStatusQueued {
		w.status = JobStatusWaiting
	}

	var jobs []*Job

	uploadHandler := func(job *Job) (*Job, error) {

		if job.status == JobStatusQueued {

			job.setStatus(JobStatusWaiting, nil)
			job.Object = w.objFn(job)

			if job.Metadata.Mode != objectify.EntModeRegular {
				msg := fmt.Sprintf("Skipping %s [invalid file format %q]", job.Metadata.FullPath(), job.Metadata.Mode.String())
				w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfSkip,
					tuipack.WARN, msg))
				job.setStatus(JobStatusSkipped, fmt.Errorf("not a valid file: %s", job.Metadata.Mode.String()))
				return job, nil
			}

			err := job.Object.Generate()
			if err != nil {
				_ = job.Object.Destroy()
				msg := fmt.Sprintf("Failed on %s [could not build object]", job.Metadata.FullPath())
				w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfFailed, tuipack.WARN, msg))
				job.setStatus(JobStatusFailed, fmt.Errorf("unable to build data object: %s", err))
				return job, nil
			}

			if w.app.Opts.Overwrite == conf.OverwriteNever {
				ex, err := w.oper.ObjectExists(job.Object)
				if ex && err != nil {
					_ = job.Object.Destroy()
					msg := fmt.Sprintf("Existing object check failed for %s", job.Metadata.FullPath())
					w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfOperFailed, tuipack.WARN, msg))
					job.setStatus(JobStatusFailed, fmt.Errorf("unable to check if object exists: %s\n", err))
					return job, nil
				}
				if ex {
					_ = job.Object.Destroy()
					msg := fmt.Sprintf("Skipping %s [object already exists]", job.Metadata.FullPath())
					w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfSkip,
						tuipack.WARN, msg))
					job.setStatus(JobStatusSkipped, fmt.Errorf("object already exists"))
					return job, nil
				}
			}

			err = job.Object.Pre()
			if err != nil {
				_ = job.Object.Destroy()
				msg := fmt.Sprintf("Object prepare failed for %s", job.Metadata.FullPath())
				w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfFailed, tuipack.WARN, msg))
				job.setStatus(JobStatusFailed, fmt.Errorf("could not initialize object: %s\n", err))
				return job, nil
			}

			err = w.oper.ObjectUpload(job.Object)
			if err != nil {
				_ = job.Object.Destroy()
				msg := fmt.Sprintf("Upload Failed: %v", err)
				w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfUploadFailed,
					tuipack.ERROR, msg))
				job.setStatus(JobStatusFailed, fmt.Errorf("could not upload object: %s\n", err))
				return job, nil
			}

			job.setStatus(JobStatusDone, nil)

			err = job.Object.Post()
			if err != nil {
				_ = job.Object.Destroy()
				job.setStatus(job.status, fmt.Errorf("post failed: %s\n", err))
				return job, nil
			}

			_ = job.Object.Destroy()

		}

		return job, nil

	}

	sets := objectify.SetsAllNoChecksums()
	if w.app.TagOpts.ChecksumSHA256 {
		sets.ChecksumSHA256 = true
	}

	if w.isDir {

		msg := fmt.Sprintf("Reading directory %s...", w.path)
		w.app.Tui.SendOutput(tuipack.NewLogMsg(tuipack.EMPTY, tuipack.ScrnNone,
			tuipack.INFO, msg))

		files, err := objectify.Path(w.path, sets)
		if err != nil {
			if strings.Contains(err.Error(), "StartingPath has no non-directory entries") {
				return
			}
			msg = fmt.Sprintf("Error reading directory %s: %s", w.path, err.Error())
			w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfFailed,
				tuipack.ERROR, msg))
			return
		} else if len(files) == 0 {
			return // there are times when objectify returns no error and no file entries.
		}

		for i := range files {

			job := newJob(w.app, files[i], w.searchRoot)
			jobs = append(jobs, job)

		}

		msg = fmt.Sprintf("Uploading directory %s...", w.path)
		w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfDefault,
			tuipack.INFO, msg))

		for {

			for i := range jobs {

				if jobs[i].status == JobStatusQueued {

					jobs[i], _ = uploadHandler(jobs[i])

				}

			}

			var breakout bool
			for i := range jobs {
				if jobs[i].status != JobStatusQueued && jobs[i].status != JobStatusWaiting {
					breakout = true
				}
			}

			if breakout {

				for i := range jobs {
					if jobs[i].status == JobStatusDone {
						w.stats.IncObjects(1, jobs[i].Metadata.SizeBytes)
					}

					if jobs[i].status == JobStatusSkipped {
						w.stats.IncSkipped(1, jobs[i].Metadata.SizeBytes)
					}

					if jobs[i].status == JobStatusFailed {
						w.stats.IncFailed(1, jobs[i].Metadata.SizeBytes)
					}
				}

				break
			}

		}

		w.app.Tui.SendOutput(tuipack.NewLogMsg("", "", tuipack.INFO, fmt.Sprintf("Upload Complete [%s]", w.path)))

	}

	if w.isFile {

		file, err := objectify.File(w.path, sets)
		if err != nil {
			if strings.Contains(err.Error(), "StartingPath has no non-directory entries") {
				return
			}
			msg := fmt.Sprintf("Error reading directory %s: %s", w.path, err.Error())
			w.app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfFailed,
				tuipack.ERROR, msg))
			return
		}

		job := newJob(w.app, file, w.searchRoot)
		job.setStatus(JobStatusQueued, nil)

		job, _ = uploadHandler(job)

	}

}
