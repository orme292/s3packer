package provider_v2

import (
	"fmt"
	"log"
	"sync"

	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
	sw "github.com/orme292/symwalker"
)

type queue struct {
	app   *conf.AppConfig
	oper  Operator
	stats *Stats

	objGenFn ObjectGenFunc

	jobs []*QueueJob
}

func newQueue(paths pathModeMap, app *conf.AppConfig, oper Operator, objFn ObjectGenFunc) (*queue, error) {

	q := &queue{
		app:   app,
		oper:  oper,
		stats: &Stats{},

		objGenFn: objFn,
	}

	sets := objectify.Sets{
		Modes: true,
	}

	for file, mode := range paths {

		if mode.IsDir() {

			log.Printf("Processing Directory: %s", file)

			opts := sw.NewSymConf(file,
				sw.WithFollowedSymLinks(),
				sw.WithFileData(),
			)

			results, err := sw.SymWalker(opts)
			if err != nil {
				log.Printf("Error processing directory %s: %s", file, err)
				continue
			}

			for i := range results.Files {

				job := &QueueJob{}

				switch results.Files[i].FileObj.Mode {
				case objectify.EntModeRegular:
					job = newJob(results.Files[i].FileObj, app, file, JobStatusWaiting, nil)
				default:
					job = newJob(results.Files[i].FileObj, app, file, JobStatusSkipped,
						fmt.Errorf("unsupported mode: %s", results.Files[i].FileObj.Mode.String()),
					)
				}

				q.addJob(job)

			}

			continue

		}

		if mode.IsRegular() {

			log.Printf("Processing File: %s", file)

			f, err := objectify.File(file, sets)
			if err != nil {
				fmt.Printf("Couldn't process file %s: %s\n", file, err)
				continue
			}

			job := &QueueJob{}

			switch f.Mode {
			case objectify.EntModeRegular:
				job = newJob(f, app, EmptyPath, JobStatusWaiting, nil)
			default:
				job = newJob(f, app, EmptyPath, JobStatusSkipped,
					fmt.Errorf("unsupported mode: %s", f.Mode.String()),
				)
			}

			q.addJob(job)
			continue

		}

		if !mode.IsRegular() && !mode.IsDir() {

			log.Printf("Processing: %s", file)

			f, err := objectify.File(file, sets)
			if err != nil {
				fmt.Printf("Couldn't process file %s: %s\n", file, err)
				continue
			}

			job := newJob(f, app, EmptyPath, JobStatusFailed,
				fmt.Errorf("unsupported mode: %s", f.Mode.String()),
			)
			q.addJob(job)

		}

	}

	if len(q.jobs) <= 0 {
		return &queue{}, fmt.Errorf("no jobs to process")
	}

	return q, nil

}

func (q *queue) addJob(job *QueueJob) { q.jobs = append(q.jobs, job) }

func (q *queue) hasPending() bool {

	for _, job := range q.jobs {
		if job.status == JobStatusWaiting {
			return true
		}
	}

	return false

}

func (q *queue) start() error {

	var wg sync.WaitGroup

	for i := q.app.Opts.MaxUploads; i > 0; i-- {

		log.Printf("spawning Worker %d", i)

		wg.Add(1)
		go q.spawnWorker(&wg, q.app)

	}

	wg.Wait()

	return nil

}

func (q *queue) spawnWorker(wg *sync.WaitGroup, app *conf.AppConfig) {

	defer wg.Done()

	for {

		for i := range q.jobs {

			if q.jobs[i].status == JobStatusWaiting {

				q.jobs[i].mu.Lock()
				if q.jobs[i].status == JobStatusWaiting {
					q.jobs[i].mu.Unlock()
					q.jobs[i].updateStatus(JobStatusQueued, nil)
				} else {
					q.jobs[i].mu.Unlock()
					continue
				}

				object := q.objGenFn(q.jobs[i])
				if object == nil {
					err := fmt.Errorf("could not generate provider object")
					q.jobs[i].updateStatus(JobStatusFailed, err)
					fmt.Printf("Start Upload Failed: %s\n", err)
					continue
				}

				ex, err := q.oper.ObjectExists(object)
				if ex && err != nil {
					_ = object.Destroy()
					q.jobs[i].updateStatus(JobStatusFailed, err)
					fmt.Printf("Duplicate Object Check Failed: %s\n", err)
					continue
				}

				if ex {
					_ = object.Destroy()
					q.jobs[i].updateStatus(JobStatusSkipped, err)
					fmt.Printf("Object %q Exists\n", q.jobs[i].Key)
					continue
				}

				fmt.Printf("Uploading %s... \n", q.jobs[i].Key)

				err = q.oper.ObjectUpload(object)
				if err != nil {
					_ = object.Destroy()
					q.jobs[i].updateStatus(JobStatusFailed, err)
					fmt.Printf("Upload Failed: %s\n", err)
					continue
				}

				_ = object.Destroy()
				q.jobs[i].updateStatus(JobStatusDone, nil)
			}

		}

		if !q.hasPending() {
			break
		}

	}

	log.Printf("A worker finished.")

}
