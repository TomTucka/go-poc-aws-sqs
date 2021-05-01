package pool

import (
	"go-elastic-reindex/internal/person"
	"log"
)

var WorkerChannel = make(chan chan *person.Person)

type Collector struct {
	Work chan *person.Person // receives jobs to send to workers
	End  chan bool           // when receives bool stops workers
}

func StartDispatcher(workerCount int, l *log.Logger) Collector {
	var i int
	var workers []Worker
	input := make(chan *person.Person) // channel to recieve work
	end := make(chan bool)             // channel to spin down workers
	collector := Collector{Work: input, End: end}

	for i < workerCount {
		i++
		l.Println("starting worker: ", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan *person.Person),
			WorkerChannel: WorkerChannel,
			End:           make(chan bool)}
		worker.Start(l)
		workers = append(workers, worker) // stores worker
	}

	// start collector
	go func() {
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop() // stop worker
				}
				return
			case work := <-input:
				worker := <-WorkerChannel // wait for available channel
				worker <- work            // dispatch work to worker
			}
		}
	}()

	return collector
}
