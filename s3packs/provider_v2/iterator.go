package provider_v2

import (
	"errors"
	"fmt"
	"log"

	"github.com/orme292/objectify"
	sw "github.com/orme292/symwalker"
)

type objJob struct {
	object *objectify.FileObj
	status int
	err    error
}

func (h *Handler) buildJobs() error {

	sets := objectify.Sets{
		Modes: true,
	}

	for file, mode := range h.paths {

		if mode.IsDir() {

			log.Printf("Processing Directory: %s", file)

			opts := sw.NewSymConf(file,
				sw.WithFollowedSymLinks(),
				sw.WithFileData(),
				sw.WithLogging())

			results, err := sw.SymWalker(opts)
			if err != nil {
				return err
			}

			for i := range results.Files {

				job := &objJob{}

				job.object = results.Files[i].FileObj
				job.status = ObjStatusWaiting
				job.err = nil

				h.jobs = append(h.jobs, job)

			}

		}

		if mode.IsRegular() {

			log.Printf("Processing File: %s", file)

			job := &objJob{}

			f, err := objectify.File(file, sets)
			if err != nil {
				fmt.Printf("Skipping %s, unable to scan.", file)
				continue
			}

			job.object = f
			job.err = nil
			job.status = ObjStatusWaiting

			if f.Mode != objectify.EntModeRegular &&
				f.Mode != objectify.EntModeDir {
				job.status = ObjStatusSkipped
				job.err = errors.New("Unsupported mode: " + f.Mode.String())
			}

			h.jobs = append(h.jobs, job)

		}

	}

	return nil

}
