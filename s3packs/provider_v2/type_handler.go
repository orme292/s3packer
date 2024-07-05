package provider_v2

import (
	"fmt"

	"github.com/orme292/s3packer/conf"
)

type Handler struct {
	app *conf.AppConfig

	objFn ObjectGenFunc

	oper  Operator
	queue *queue

	Stats    *Stats
	supports *Supports
}

func NewHandler(app *conf.AppConfig, operFn OperGenFunc, objFn ObjectGenFunc) (*Handler, error) {

	fmt.Printf("Starting s3packer...\n\n")

	oper, err := operFn(app)
	if err != nil {
		return nil, fmt.Errorf("error during oper gen: %w", err)
	}

	h := &Handler{
		app:      app,
		oper:     oper,
		Stats:    &Stats{},
		supports: oper.Support(),
	}

	err = h.verifyBucket()
	if err != nil {
		return nil, err
	}

	paths := combinePaths(app.Dirs, app.Files)
	h.queue, err = newQueue(paths, app, oper, objFn)
	if err != nil {
		return nil, fmt.Errorf("error building queue: %w", err)
	}

	// TODO: Remove Logging
	for i := range h.queue.jobs {
		fmt.Printf("Job Path: %s (%s)\n", h.queue.jobs[i].Details.FullPath(), h.queue.jobs[i].SearchRoot)
	}

	return h, nil

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

	return h.queue.start()

}
