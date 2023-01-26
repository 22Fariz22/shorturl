package worker

import (
	"errors"
	"github.com/22Fariz22/shorturl/internal/repository"
	"log"
	"sync"
)

type Pool struct {
	wg         sync.WaitGroup
	once       sync.Once
	shutDown   chan struct{}
	mainCh     chan workerData
	repository repository.Repository
}

type workerData struct {
	urls   []string
	cookie string
}

func (w *Pool) AddJob(arr []string, cookies string) error {
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

func (w *Pool) Stop() {
	w.once.Do(func() {
		close(w.shutDown)
		close(w.mainCh)
	})
	w.wg.Wait()
}

func NewWorkerPool(repo repository.Repository) *Pool {
	return &Pool{
		wg:         sync.WaitGroup{},
		once:       sync.Once{},
		shutDown:   make(chan struct{}),
		mainCh:     make(chan workerData, 10),
		repository: repo,
	}
}
