package provider_v2

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/orme292/s3packer/conf"
)

type Handler struct {
	app *conf.AppConfig

	oper Operator
	iter Iterator

	jobs []*objJob

	paths map[string]fs.FileMode

	Stats    *Stats
	supports *Supports
}

func NewHandler(app *conf.AppConfig, operFn OperGenFunc, iterFn IterGenFunc) (*Handler, error) {

	fmt.Printf("Starting s3packer...\n\n")

	oper, err := operFn(app)
	if err != nil {
		return nil, fmt.Errorf("error during oper gen: %w", err)
	}

	iter, err := iterFn(app)
	if err != nil {
		return nil, fmt.Errorf("error during iter gen: %w", err)
	}

	h := &Handler{
		app:      app,
		oper:     oper,
		iter:     iter,
		supports: oper.Support(),
		Stats:    &Stats{},
	}

	err = h.verifyBucket()
	if err != nil {
		return nil, err
	}

	h.combinePaths(app.Dirs, app.Files)

	h.dropBrokenPaths()

	return h, nil
}

func (h *Handler) combinePaths(dirs []string, files []string) {

	h.paths = make(map[string]os.FileMode)

	for _, dir := range dirs {

		info, err := os.Stat(dir)
		if err != nil {
			h.paths[dir] = fs.ModeIrregular
			continue
		}

		h.paths[dir] = info.Mode()

	}

	for _, file := range files {

		info, err := os.Stat(file)
		if err != nil {
			h.paths[file] = fs.ModeIrregular
			continue
		}

		h.paths[file] = info.Mode()

	}

	// TODO: REMOVE
	for name, mode := range h.paths {
		log.Printf("Added Path: %s [%v]", name, mode.String())
	}

}

func (h *Handler) createBucket() error {

	if !h.supports.BucketCreate {
		return fmt.Errorf("bucket creation not supported")
	}

	if h.supports.BucketCreate && h.app.Bucket.Create {

		fmt.Printf("Creating bucket... ")

		err := h.oper.BucketCreate()
		if err != nil {
			fmt.Printf("failed.\n")
			return fmt.Errorf("unable to create bucket: %w", err)
		}

		exists, err := h.oper.BucketExists()
		if err != nil {
			fmt.Printf("failed.\n")
			return fmt.Errorf("unable to check for bucket: %w", err)
		}
		if !exists {
			fmt.Printf("failed.\n")
			return fmt.Errorf("created bucket but couldn't verify it exists")
		}

	}

	return nil

}

func (h *Handler) dropBrokenPaths() {

	for path, mode := range h.paths {

		if !mode.IsDir() && !mode.IsRegular() {
			fmt.Printf("Skipping inaccessible path: %s\n", path)
			delete(h.paths, path)
		}

	}

}

func (h *Handler) verifyBucket() error {

	// Check if bucket exists. If it does not, create it.
	fmt.Printf("Verifying bucket... ")

	exists, err := h.oper.BucketExists()
	if err != nil && err.Error() != "bucket not found" {
		fmt.Printf("could not verify.\n")
		return fmt.Errorf("error while checking for bucket: %w", err)
	}

	if !exists {

		fmt.Printf("not found.\n")

		err = h.createBucket()
		if err != nil {
			return err
		}

		fmt.Printf("bucket created\n")

	} else {

		fmt.Printf("OK\n")

	}

	return nil

}

func (h *Handler) Run() error {

	err := h.buildJobs()
	if err != nil {
		return err
	}

	for _, job := range h.jobs {
		fmt.Printf("File: %s\n\tMode: %s\n", job.object.Filename, job.object.Mode.String())
	}

	return nil

}
