package pool

import (
	"encoding/json"
	"errors"
	"go-elastic-reindex/internal/person"
	"log"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch"
)

type Worker struct {
	ID            int
	WorkerChannel chan chan *person.Person
	Channel       chan *person.Person
	End           chan bool
}

func (w *Worker) Start(l *log.Logger) {
	go func() {
		for {
			w.WorkerChannel <- w.Channel // when the worker is available place channel in queue
			select {
			case work_request := <-w.Channel: // worker has received job
				// err := indexer(l, *work_request)
				// if err != nil {
				// 	log.Println(err)
				// }
				l.Println("Person Indexed ", *&work_request.UID)

			case <-w.End:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}

func indexer(l *log.Logger, p person.Person) error {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return errors.New("Cannot create ES client")
	}

	data, _ := json.Marshal(p)

	res, err := es.Index(
		"persons",                       // Index name
		strings.NewReader(string(data)), // Document body
		es.Index.WithDocumentID(p.UID),  // Document ID
		es.Index.WithRefresh("true"),    // Refresh
	)

	log.Println(res)

	if err != nil {
		l.Fatal(err)
		return errors.New("Error indexing document")
	}
	defer res.Body.Close()
	if res.IsError() {
		l.Fatal(err)
		return errors.New("Error indexing document")
	}
	return nil
}
