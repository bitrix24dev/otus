package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	// Канал для передачи задач горутинам. Задачи отправляются в этот канал и извлекаются горутинами для выполнения
	tasksCh := make(chan Task)

	// Канал для уведомления горутин о том, что достигнут лимит ошибок, и они должны завершить выполнение
	errorsCh := make(chan struct{})

	// Канал для уведомления о завершении всех задач. Используется для корректного завершения горутин
	doneCh := make(chan struct{})

	var wg sync.WaitGroup

	//Счетчик ошибок, который увеличивается при возникновении ошибки в задаче.
	var errorsCount int32

	// Используется для гарантии, что `errorsCh` будет закрыт только один раз, даже если несколько горутин одновременно обнаружат превышение лимита ошибок.
	var once sync.Once

	/*
		Функция-работник:
		- Запускается `n` горутин, каждая из которых выполняет задачи из `tasksCh`.
		- Если задача возвращает ошибку, счетчик ошибок увеличивается. Если количество ошибок достигает `m`, канал `errorsCh` закрывается, что сигнализирует всем горутинам о необходимости завершения.
		- Горутина завершает выполнение, если канал `tasksCh` закрыт или если канал `errorsCh` или `doneCh` закрыт.
	*/
	worker := func() {
		defer wg.Done()
		for {
			select {
			case task, ok := <-tasksCh:
				if !ok {
					return
				}
				if err := task(); err != nil {
					if atomic.AddInt32(&errorsCount, 1) >= int32(m) {
						once.Do(func() { close(errorsCh) })
					}
				}
			case <-errorsCh:
				return
			case <-doneCh:
				return
			}
		}
	}

	// Start n workers
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}

	// Send tasks to workers
	/*
		- Задачи отправляются в `tasksCh` в отдельной горутине.
		- Если `errorsCh` закрыт, отправка задач прекращается.
	*/
	go func() {
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case tasksCh <- task:
			case <-errorsCh:
				return
			case <-doneCh:
				return
			}
		}
	}()

	/*
		- Используется `sync.WaitGroup` для ожидания завершения всех горутин-работников.
		- После завершения всех горутин `doneCh` закрывается, чтобы гарантировать, что все горутины завершились.
	*/
	wg.Wait()
	close(doneCh)

	/*
		- Если количество ошибок достигло `m`, функция возвращает `ErrErrorsLimitExceeded`.
		- В противном случае функция возвращает `nil`, что означает успешное выполнение всех задач без превышения лимита ошибок.
	*/

	if errorsCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func main() {
	// Example test case
	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}

	workersCount := 10
	maxErrorsCount := 23
	err := Run(tasks, workersCount, maxErrorsCount)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Run tasks count: %d\n", runTasksCount)
}
