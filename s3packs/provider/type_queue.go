package provider

import (
	"fmt"
	"sync"

	"github.com/orme292/s3packer/conf"
	sw "github.com/orme292/symwalker"
)

type queue struct {
	app   *conf.AppConfig
	oper  Operator
	stats *Stats

	objGenFn ObjectGenFunc

	workers []*worker
}

func newQueue(paths pathModeMap, app *conf.AppConfig, oper Operator, objFn ObjectGenFunc) (*queue, error) {

	q := &queue{
		app:      app,
		oper:     oper,
		stats:    &Stats{},
		objGenFn: objFn,
	}

	for file, mode := range paths {

		if mode.IsDir() {

			msg := fmt.Sprintf("Walking %s", file)
			app.Tui.Info(msg)

			opts := sw.NewSymConf(file,
				sw.WithoutFiles(),
				sw.WithDepth(sw.INFINITE),
			)

			results, err := sw.SymWalker(opts)
			if err != nil {
				app.Tui.Error(msg)
				continue
			}

			for i := range results.Dirs {

				j := newWorker(app, results.Dirs[i].Path, file, true, false, JobStatusQueued, oper, objFn)
				q.workers = append(q.workers, j)

			}

		} else {

			msg := fmt.Sprintf("Reading %s", file)
			app.Tui.Info(msg)

			j := newWorker(app, file, EmptyPath, false, true, JobStatusQueued, oper, objFn)
			q.workers = append(q.workers, j)

		}

	}

	return q, nil

}

func (q *queue) start() {

	sem := make(chan struct{}, q.app.Opts.MaxUploads)
	var wg sync.WaitGroup
	wg.Add(len(q.workers))

	stats := &Stats{}
	var agg []*Stats

	for i := range q.workers {

		go func(worker *worker, app *conf.AppConfig) {

			sem <- struct{}{}
			defer func() { <-sem }()

			worker.scan()
			agg = append(agg, worker.stats)

			wg.Done()

		}(q.workers[i], q.app)

	}

	wg.Wait()

	for _, stat := range agg {
		if stat != nil {
			stats.Merge(stat)
		}
	}

	q.stats = stats

}
