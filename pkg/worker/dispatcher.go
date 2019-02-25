package worker

import (
	"log"
)

// Dispatcher -
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maxWorkers int
	WorkerPool chan chan Job
	JobQueue   chan Job
}

// NewDispatcher -
func NewDispatcher(maxWorkers int, maxQueue int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	queue := make(chan Job, maxQueue)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers, JobQueue: queue}
}

// Run -
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	log.Print("Worker queue dispatcher started...")
	for {
		select {
		case job := <-d.JobQueue:
			log.Printf("Dispatcher request received")
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
