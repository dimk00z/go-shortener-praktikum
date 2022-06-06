package worker

import (
	"context"
	"log"
	"sync"
)

type WorkersPool struct {
	workersNumber int
	inputCh       chan func(ctx context.Context) error
	done          chan struct{}
}

func NewWorkersPool(workersNumber int, poolLegth int) *WorkersPool {
	return &WorkersPool{
		workersNumber: workersNumber,
		inputCh:       make(chan func(ctx context.Context) error, poolLegth),
		done:          make(chan struct{}),
	}
}

func (wp *WorkersPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}

func doWorkersTask(ctx context.Context,
	workerIndex int,
	wg *sync.WaitGroup,
	taskCh chan func(ctx context.Context) error) {
	log.Printf("worker %v started\n", workerIndex)
workerLoop:
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %v got context.Done\n", workerIndex)
			break workerLoop
		case workerTask := <-taskCh:
			if err := workerTask(ctx); err != nil {
				log.Printf("worker %v got error:%s", workerIndex, err.Error())
			}
		}
	}
	wg.Done()
}

func (wp *WorkersPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for workerIndex := 0; workerIndex < wp.workersNumber; workerIndex++ {
		wg.Add(1)
		go doWorkersTask(ctx, workerIndex, wg, wp.inputCh)
	}
	wg.Wait()
	close(wp.inputCh)
}

func (wp *WorkersPool) Close() {
	close(wp.done)
}
