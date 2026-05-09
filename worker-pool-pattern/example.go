package main

import (
	"fmt"
	"sync"
	"time"
)

// Job represents work to be done
type Job struct {
	ID    int
	Value int
}

// Result represents job outcome
type Result struct {
	JobID  int
	Output int
	Worker int
}

// worker processes jobs from channel
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Simulate processing
		time.Sleep(100 * time.Millisecond)
		output := job.Value * 2

		results <- Result{
			JobID:  job.ID,
			Output: output,
			Worker: id,
		}
	}
}

// BasicWorkerPool demonstrates simple worker pool
func BasicWorkerPool(numWorkers int, jobs []Job) []Result {
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobChan, resultChan, &wg)
	}

	// Close results when all workers done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Send jobs
	go func() {
		for _, job := range jobs {
			jobChan <- job
		}
		close(jobChan)
	}()

	// Collect results
	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// WorkerPool provides reusable worker pool
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
}

// NewWorkerPool creates pool with fixed workers
func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
	pool := &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, queueSize),
		results:    make(chan Result, queueSize),
	}

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		pool.wg.Add(1)
		go worker(w, pool.jobs, pool.results, &pool.wg)
	}

	return pool
}

// Submit sends job to pool
func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

// Close stops accepting jobs and waits for completion
func (p *WorkerPool) Close() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

// Results returns result channel
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

func main() {
	// Example 1: Basic worker pool
	fmt.Println("=== Example 1: Basic Worker Pool ===")
	jobs := make([]Job, 10)
	for i := range jobs {
		jobs[i] = Job{ID: i, Value: i + 1}
	}

	start := time.Now()
	results := BasicWorkerPool(3, jobs)
	duration := time.Since(start)

	fmt.Printf("Processed %d jobs with 3 workers in %v\n", len(results), duration)
	for _, r := range results {
		fmt.Printf("  Job %d: %d (by worker %d)\n", r.JobID, r.Output, r.Worker)
	}

	// Example 2: Reusable worker pool
	fmt.Println("\n=== Example 2: Reusable Pool ===")
	pool := NewWorkerPool(2, 5)

	// Submit jobs
	go func() {
		for i := 0; i < 6; i++ {
			pool.Submit(Job{ID: i, Value: (i + 1) * 10})
		}
		pool.Close()
	}()

	// Collect results
	fmt.Println("Results:")
	for result := range pool.Results() {
		fmt.Printf("  Job %d: %d (worker %d)\n", result.JobID, result.Output, result.Worker)
	}

	// Example 3: Compare sequential vs parallel
	fmt.Println("\n=== Example 3: Performance Comparison ===")

	// Sequential
	start = time.Now()
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
	}
	sequential := time.Since(start)

	// Parallel with 5 workers
	start = time.Now()
	parallelJobs := make([]Job, 10)
	for i := range parallelJobs {
		parallelJobs[i] = Job{ID: i, Value: i}
	}
	BasicWorkerPool(5, parallelJobs)
	parallel := time.Since(start)

	fmt.Printf("Sequential: %v\n", sequential)
	fmt.Printf("Parallel (5 workers): %v\n", parallel)
	fmt.Printf("Speedup: %.2fx\n", float64(sequential)/float64(parallel))
}
