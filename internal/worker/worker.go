// Package worker пакет для работы воркера
package worker

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/22Fariz22/shorturl/internal/usecase"
)

//Pool структура для воркера
type Pool struct {
	wg         sync.WaitGroup
	once       sync.Once
	shutDown   chan struct{}
	mainCh     chan workerData
	repository usecase.Repository
}

//NewWorkerPool создание воркера
func NewWorkerPool(repo usecase.Repository) *Pool {
	return &Pool{
		wg:         sync.WaitGroup{},
		once:       sync.Once{},
		shutDown:   make(chan struct{}),
		mainCh:     make(chan workerData, 10),
		repository: repo,
	}
}

//workerData структура содержания воркера
type workerData struct {
	urls   []string
	cookie string
}

//AddJob запуск в handler
func (w *Pool) AddJob(ctx context.Context, arr []string, cookies string) error {
	select {
	case <-w.shutDown:
		return errors.New("all channels are closed")
	case w.mainCh <- workerData{
		urls:   arr,
		cookie: cookies,
	}:
		return nil
	}
}

//RunWorkers запуск в App
func (w *Pool) RunWorkers(count int) {
	for i := 0; i < count; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case <-w.shutDown:
					return
				case urls, ok := <-w.mainCh:
					if !ok {
						return
					}

					err := w.repository.Delete(urls.urls, urls.cookie)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}()
	}
}

//Stop остановка воркера
func (w *Pool) Stop() {
	w.once.Do(func() {
		close(w.shutDown)
		close(w.mainCh)
	})
	w.wg.Wait()
}
