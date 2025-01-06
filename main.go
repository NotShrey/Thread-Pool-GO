package main

import (
	"fmt"
	"sync"
	"time"
)

// Job type representing a unit of work
type Job func()

// Pool struct representing a worker pool
type Pool struct {
	workQueue chan Job
	wg        sync.WaitGroup
}

// NewPool initializes a new worker pool with the specified number of workers
func NewPool(workerCount int) *Pool {
	p := &Pool{
		workQueue: make(chan Job),
	}

	// Start the workers
	for i := 0; i < workerCount; i++ {
		go p.worker()
	}

	return p
}

// worker processes jobs from the work queue
func (p *Pool) worker() {
	for job := range p.workQueue {
		job()
		p.wg.Done() // Mark job as done
	}
}

// AddJob adds a job to the pool
func (p *Pool) AddJob(job Job) {
	p.wg.Add(1) // Increment the WaitGroup counter
	p.workQueue <- job
}

// Wait waits for all jobs to complete
func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.workQueue) // Close the work queue when all jobs are done
}

func main() {
	pool := NewPool(5)

	for i := 0; i < 30; i++ {
		jobID := i // Capture loop variable
		job := func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("Job %d: completed\n", jobID)
		}
		pool.AddJob(job)
	}

	pool.Wait()
}
