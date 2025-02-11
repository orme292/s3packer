package provider

import (
	"fmt"
	"time"

	"s3p/internal/conf"
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

	oper, err := operFn(app)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		app:      app,
		oper:     oper,
		Stats:    &Stats{},
		supports: oper.Support(),
	}

	paths := combinePaths(app.Dirs, app.Files)
	h.queue, err = newQueue(paths, app, oper, objFn)
	if err != nil {
		return nil, err
	}

	return h, nil

}

func (h *Handler) Init() error {

	err := h.verifyBucket()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 2)

	h.queue.start()

	h.Stats.Merge(h.queue.stats)

	return nil

}

func (h *Handler) createBucket() error {

	if !h.supports.BucketCreate {
		return fmt.Errorf("bucket creation not supported")
	}

	if h.supports.BucketCreate && h.app.Bucket.Create {

		h.app.Log.Info("Creating bucket")

		err := h.oper.BucketCreate()
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}

		exists, err := h.oper.BucketExists()
		if err != nil {
			return fmt.Errorf("check bucket: %w", err)
		}
		if !exists {
			return fmt.Errorf("created bucket but couldn't verify it exists")
		}

		h.app.Log.Info("Bucket Created")

	}

	return nil

}

func (h *Handler) verifyBucket() error {

	// Check if bucket exists. If it does not, create it.
	exists, err := h.oper.BucketExists()
	if err != nil && err.Error() != "bucket not found" {
		h.app.Log.Warn("Bucket not found.")
		return fmt.Errorf("error while checking for bucket: %w", err)
	}

	if !exists {

		err = h.createBucket()
		if err != nil {
			return err
		}

	} else {

	}

	return nil

}
