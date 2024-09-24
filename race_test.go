package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestDataRaceConditions(t *testing.T) {
		// Data races are the enemy of mankind and concurrency.
		// They're your free ticket to flipping burgers if you ever enter a 
		// company as a not-CEO.
		var state int32
		var wg sync.WaitGroup	
		// var mu sync.Mutex
		// Use mutexes sparingly. They're not flexible, and they will slow down your code
		// Use them only when the situation is simple enough, or when you don't have alternatives.
		// Otherwise, prefer coding patterns with channels. main.go shows one such pattern 
		
		// When operation is simple enough like here, you can use 
		// atomic operations or atomic values from the sync package

		// The sync/atomic package itself warns about using atomicity. They're even more delicate
		// to handle than locks, and channels should be preferred when possible.

		for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(i int) {
						defer wg.Done()
						// mu.Lock()
						// defer mu.Unlock()
						fmt.Println("Goroutine number", i)
						// state += int32(i)
						atomic.AddInt32(&state, int32(i))
						fmt.Println(state)
				}(i)
		} 
		wg.Wait()
		fmt.Println("Done")
}
