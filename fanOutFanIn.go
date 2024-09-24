// A small piece of marvel from ChatGPT4o mini. For how short it is, it holds a great deal of knowledge
package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker function that processes a number
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate some work
		time.Sleep(time.Millisecond * 100)
		result := job * 2 // Example processing: double the number
		fmt.Printf("Worker %d processed job %d, result: %d\n", id, job, result)
		results <- result // Send the result to the results channel
	}
}

// Fan-out function to distribute jobs to workers
func fanOut(jobs <-chan int, numWorkers int) (chan int, *sync.WaitGroup) {
	results := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg) // Start a worker goroutine
	}

	return results, &wg
}

// Fan-in function to merge results from multiple channels
func fanIn(channels ...chan int) chan int {
	out := make(chan int)

	go func() {
		for _, ch := range channels {
			for result := range ch {
				out <- result // Send results to the output channel
			}
		}
		close(out) // Close the output channel when done
	}()

	return out
}

func notMain() {
	jobs := make(chan int, 10) // Channel for jobs
	numWorkers := 3            // Number of worker goroutines

	// Start fan-out
	results, wg := fanOut(jobs, numWorkers)

	// Send jobs to the jobs channel
	for i := 1; i <= 10; i++ {
		jobs <- i // Send job
	}
	close(jobs) // Close the jobs channel to signal no more jobs

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results) // Close results channel when all workers are done
	}()

	// Collect results using fan-in
	finalResults := fanIn(results)

	// Print final results
	for result := range finalResults {
		fmt.Println("Final result:", result)
	}
}

