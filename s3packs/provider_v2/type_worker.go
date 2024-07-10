package provider_v2

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
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

		log.Printf("Reading directory %s", w.path)

		files, err := objectify.Path(w.path, objectify.SetsAllNoChecksums())
		if err != nil {
			log.Printf("Error reading directory %s: %s", w.path, err)
			return
		}

		for i := range files {

			job := newJob(w.app, files[i], w.searchRoot)
			job.setStatus(JobStatusQueued, nil)
			jobs = append(jobs, job)

		}

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

					log.Printf("[%s] Checking if object exists: %s", w.uuid, jobs[i].Key)

					if w.app.Opts.Overwrite == conf.OverwriteNever {
						ex, err := w.oper.ObjectExists(jobs[i].Object)
						if ex && err != nil {
							_ = jobs[i].Object.Destroy()
							jobs[i].setStatus(JobStatusFailed, fmt.Errorf("Duplicate Object Check Failed: %s\n", err))
							continue
						}
						if ex {
							fmt.Println("Object already exists")
							_ = jobs[i].Object.Destroy()
							jobs[i].setStatus(JobStatusSkipped, fmt.Errorf("Object already exists"))
							continue
						}
					}

					log.Printf("[%s] starting upload: %s", w.uuid, jobs[i].Key)

					err = jobs[i].Object.Pre()
					if err != nil {
						_ = jobs[i].Object.Destroy()
						jobs[i].setStatus(JobStatusFailed, fmt.Errorf("could not initialize object: %s\n", err))
						continue
					}

					err = w.oper.ObjectUpload(jobs[i].Object)
					if err != nil {
						_ = jobs[i].Object.Destroy()
						log.Printf("ERROR: %s\n", err)
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
					fmt.Println("Waiting on job")
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
		fmt.Println("SKIPPING FILE")
	}

}
