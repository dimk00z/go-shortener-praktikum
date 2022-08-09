package worker

import (
	"context"
	"sync"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type IWorkerPool interface {
	Push(task func(ctx context.Context) error)
	Run(ctx context.Context)
	Close()
}
type WorkersPool struct {
	workersNumber int
	inputCh       chan func(ctx context.Context) error
	done          chan struct{}
	l             *logger.Logger
}

var (
	wp   IWorkerPool
	once sync.Once
)

func NewWorkersPool(workersNumber int, poolLength int, l *logger.Logger) *WorkersPool {
	return &WorkersPool{
		workersNumber: workersNumber,
		inputCh:       make(chan func(ctx context.Context) error, poolLength),
		done:          make(chan struct{}),
		l:             l,
	}
}

func (wp *WorkersPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}

func (wp *WorkersPool) doTasksByWorkers(ctx context.Context,
	workerIndex int,
	taskCh chan func(ctx context.Context) error) error {
	wp.l.Debug("worker_%v started", workerIndex)
workerLoop:
	for {
		select {
		case <-ctx.Done():
			wp.l.Debug("worker_%v got context.Done", workerIndex)
			break workerLoop
		case workerTask := <-taskCh:
			wp.l.Debug("worker_%v is busy", workerIndex)
			if err := workerTask(ctx); err != nil {
				wp.l.Debug("worker_%v got error:%s", workerIndex, err.Error())
				return err
			} else {
				wp.l.Debug("worker %v finished task correctly", workerIndex)
			}
		}
	}
	return nil
}

func (wp *WorkersPool) Run(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	for workerIndex := 0; workerIndex < wp.workersNumber; workerIndex++ {
		workerIndex := workerIndex
		g.Go(func() error {
			return wp.doTasksByWorkers(ctx, workerIndex, wp.inputCh)
		})
	}
	if err := g.Wait(); err != nil {
		wp.l.Debug(err)
	}
	close(wp.inputCh)
}

func (wp *WorkersPool) Close() {
	close(wp.done)
}

func GetWorkersPool(l *logger.Logger, wpConfig config.Workers) IWorkerPool {
	once.Do(func() {
		wp = NewWorkersPool(wpConfig.WorkersNumber, wpConfig.PoolLength, l)
	})
	return wp
}
