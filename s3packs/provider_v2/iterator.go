package provider_v2

import (
	"errors"
	"fmt"

	objf "github.com/orme292/objectify"
)

type objJob struct {
	object *objf.FileObj
	status int
	err    error
}

func (h *Handler) buildJobs() error {

	sets := objf.Sets{
		Modes: true,
	}

	for file, mode := range h.paths {

		if mode.IsDir() {

			files, err := objf.Path(file, sets)
			if err != nil {
				fmt.Printf("Skipping %s, unable to scan.", file)
				continue
			}

			for i := range files {

				job := &objJob{}

				if files[i].Mode == objf.EntModeRegular {
					job.object = files[i]
					job.status = ObjStatusWaiting
					job.err = nil
				}

				if files[i].Mode != objf.EntModeRegular &&
					files[i].Mode != objf.EntModeDir {
					job.object = files[i]
					job.status = ObjStatusSkipped
					job.err = errors.New("Unsupported mode: " + files[i].Mode.String())
				}

				h.jobs = append(h.jobs, job)

			}

			if mode.IsRegular() {

				job := &objJob{}

				f, err := objf.File(file, sets)
				if err != nil {
					fmt.Printf("Skipping %s, unable to scan.", file)
					continue
				}

				job.object = f
				job.err = nil
				job.status = ObjStatusWaiting

				if f.Mode != objf.EntModeRegular &&
					f.Mode != objf.EntModeDir {
					job.status = ObjStatusSkipped
					job.err = errors.New("Unsupported mode: " + f.Mode.String())
				}

			}

		}

	}

	return nil

}
