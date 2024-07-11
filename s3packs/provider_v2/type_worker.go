package provider_v2

import (
	"fmt"

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

	if w.isDir {

		msg := fmt.Sprintf("Reading directory %s...", w.path)
		w.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: msg, Mark: false}, msg, tuipack.INFO, true, true)

		files, err := objectify.Path(w.path, objectify.SetsAllNoChecksums())
		if err != nil {
			errMsg := fmt.Sprintf("Error reading directory %s: %s", w.path, err.Error())
			w.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: errMsg, Mark: false},
				errMsg, tuipack.ERROR, true, true)
			return
		}

		for i := range files {

			job := newJob(w.app, files[i], w.searchRoot)
			job.setStatus(JobStatusQueued, nil)
			jobs = append(jobs, job)

		}

		msg = fmt.Sprintf("Uploading directory %s...", w.path)
		w.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: msg, Mark: false}, msg, tuipack.INFO, true, true)

		for {

			for i := range jobs {

				if jobs[i].status == JobStatusQueued {

					jobs[i].setStatus(JobStatusWaiting, nil)

					jobs[i].Object = w.objFn(jobs[i])

					if jobs[i].Metadata.Mode != objectify.EntModeRegular {
						jobs[i].setStatus(JobStatusSkipped, fmt.Errorf("unsupported file mode: %s", files[i].Mode.String()))
						continue
					}

					err = jobs[i].Object.Generate()
					if err != nil {
						_ = jobs[i].Object.Destroy()
						jobs[i].setStatus(JobStatusFailed, fmt.Errorf("unable to build data object: %s", err))
						continue
					}

					if w.app.Opts.Overwrite == conf.OverwriteNever {
						ex, err := w.oper.ObjectExists(jobs[i].Object)
						if ex && err != nil {
							_ = jobs[i].Object.Destroy()
							jobs[i].setStatus(JobStatusFailed, fmt.Errorf("Duplicate Object Check Failed: %s\n", err))
							continue
						}
						if ex {
							// w.app.ScreenSend("Object Exists", "", true)
							_ = jobs[i].Object.Destroy()
							jobs[i].setStatus(JobStatusSkipped, fmt.Errorf("Object already exists"))
							continue
						}
					}

					err = jobs[i].Object.Pre()
					if err != nil {
						_ = jobs[i].Object.Destroy()
						jobs[i].setStatus(JobStatusFailed, fmt.Errorf("could not initialize object: %s\n", err))
						continue
					}

					err = w.oper.ObjectUpload(jobs[i].Object)
					if err != nil {
						_ = jobs[i].Object.Destroy()
						msg = fmt.Sprintf("Upload Failed: %v", err)
						w.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: msg, Mark: false}, msg, tuipack.ERROR, true, true)
						jobs[i].setStatus(JobStatusFailed, fmt.Errorf("could not upload object: %s\n", err))
						continue
					}

					jobs[i].setStatus(JobStatusDone, nil)

					err = jobs[i].Object.Post()
					if err != nil {
						_ = jobs[i].Object.Destroy()
						jobs[i].setStatus(jobs[i].status, fmt.Errorf("post failed: %s\n", err))
						continue
					}

					_ = jobs[i].Object.Destroy()

				}

			}

			breakout := true
			for i := range jobs {
				if jobs[i].status == JobStatusQueued || jobs[i].status == JobStatusWaiting {
					breakout = false
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

				return

			}

		}

	}

	if w.isFile {
		w.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: "File Skipped", Mark: true}, "File Skipped", tuipack.WARN, true, true)
	}

}
