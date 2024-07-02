package provider_v2

import (
	"fmt"
	"log"

	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
	sw "github.com/orme292/symwalker"
)

type queue struct {
	app  *conf.AppConfig
	jobs []*queueJob
}

func newQueue(paths pathModeMap) (*queue, error) {

	q := &queue{}

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

				job := &queueJob{}

				switch results.Files[i].FileObj.Mode {
				case objectify.EntModeRegular:
					job = newJob(results.Files[i].FileObj, file, JobStatusWaiting, nil)
				default:
					job = newJob(results.Files[i].FileObj, file, JobStatusSkipped,
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

			job := &queueJob{}

			switch f.Mode {
			case objectify.EntModeRegular:
				job = newJob(f, EmptyPath, JobStatusWaiting, nil)
			default:
				job = newJob(f, EmptyPath, JobStatusSkipped,
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

			job := newJob(f, EmptyPath, JobStatusFailed,
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

func (q *queue) addJob(job *queueJob) { q.jobs = append(q.jobs, job) }

func (q *queue) start() error {

	// var wg sync.WaitGroup
	//
	// for {
	//
	// }

	return nil
}

func (q *queue) worker() error {

	// for _, job := range q.jobs {
	//
	// }

	return nil
}
