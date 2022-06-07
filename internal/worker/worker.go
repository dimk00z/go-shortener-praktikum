package worker

import (
	"context"
	"log"
	"sync"

	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
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
}

var (
	wp   IWorkerPool
	once sync.Once
)

func NewWorkersPool(workersNumber int, poolLength int) *WorkersPool {
	return &WorkersPool{
		workersNumber: workersNumber,
		inputCh:       make(chan func(ctx context.Context) error, poolLength),
		done:          make(chan struct{}),
	}
}

func (wp *WorkersPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}

func doTasksByWorkers(ctx context.Context,
	workerIndex int,
	wg *sync.WaitGroup,
	taskCh chan func(ctx context.Context) error) {
	defer wg.Done()
	log.Printf("worker_%v started\n", workerIndex)
workerLoop:
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker_%v got context.Done\n", workerIndex)
			break workerLoop
		case workerTask := <-taskCh:
			log.Printf("worker_%v is busy\n", workerIndex)
			if err := workerTask(ctx); err != nil {
				log.Printf("worker_%v got error:%s", workerIndex, err.Error())
			} else {
				log.Printf("worker %v finished task correctly", workerIndex)
			}
		}
	}
}

func (wp *WorkersPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for workerIndex := 0; workerIndex < wp.workersNumber; workerIndex++ {
		wg.Add(1)
		go doTasksByWorkers(ctx, workerIndex, wg, wp.inputCh)
	}
	wg.Wait()
	close(wp.inputCh)
}

func (wp *WorkersPool) Close() {
	close(wp.done)
}

func GetWorkersPool(wpConfig settings.WorkersConfig) IWorkerPool {
	once.Do(func() {
		wp = NewWorkersPool(wpConfig.WorkersNumber, wpConfig.PoolLength)
	})
	return wp
}
