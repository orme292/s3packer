package provider_v2

import (
	"fmt"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/tuipack"
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

	app.Tui.ToScreenHeader("Running...")

	paths := combinePaths(app.Dirs, app.Files)
	h.queue, err = newQueue(paths, app, oper, objFn)
	if err != nil {
		return nil, fmt.Errorf("error building queue2: %w", err)
	}

	return h, nil

}

func (h *Handler) createBucket() error {

	if !h.supports.BucketCreate {
		return fmt.Errorf("bucket creation not supported")
	}

	h.app.Tui.ToScreenHeader("Creating bucket...")

	if h.supports.BucketCreate && h.app.Bucket.Create {

		h.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: "Creating bucket...", Mark: false},
			"Creating bucket", tuipack.INFO, true, true)

		err := h.oper.BucketCreate()
		if err != nil {
			return fmt.Errorf("unable to create bucket: %w", err)
		}

		exists, err := h.oper.BucketExists()
		if err != nil {
			return fmt.Errorf("unable to check for bucket: %w", err)
		}
		if !exists {
			return fmt.Errorf("created bucket but couldn't verify it exists")
		}

		h.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: "Bucket Created.", Mark: true},
			"Bucket Created", tuipack.INFO, true, true)

	}

	h.app.Tui.ResetHeader()

	return nil

}

func (h *Handler) verifyBucket() error {

	// Check if bucket exists. If it does not, create it.
	h.app.Tui.ToScreen("Looking for bucket...", false)

	exists, err := h.oper.BucketExists()
	if err != nil && err.Error() != "bucket not found" {
		h.app.Tui.SendOutput(tuipack.ScreenMsg{Msg: "Bucket not found.", Mark: false},
			"Bucket not found", tuipack.INFO, true, true)
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

func (h *Handler) Run() error {

	h.queue.start()
	return nil

}
