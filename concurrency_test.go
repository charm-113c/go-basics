package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
		// The basics of concurrency
		now := time.Now()
		userID := 10
		respCh := make(chan string, 3)
		// One of the most useful features we will use 
		wg := &sync.WaitGroup{}

		fmt.Println("Beginning concurrency test")
		// we need to run multiple operations
		go fetchUserData(userID, respCh, wg)
		go fetchUserRecommendations(userID, respCh, wg)	
		go fetchUserLikes(userID, respCh, wg)
		wg.Add(3)
		// We add all three goroutines to the wait group, although we could 
		// add them one by one too

		wg.Wait()
		// And then we wait for the slowest of them to finish before
		close(respCh)

		for resp := range respCh {
				// But, unless channels are closed, this will loop forever
				fmt.Println(resp)
				// In short, there's a deadlock, all goroutines are blocked
				// And yet, closing early means goroutines may never get to 
				// finish their jobs. So we have waitGroups
		}

		fmt.Println(time.Since(now))
}

func fetchUserData(userID int, respCh chan string, wg *sync.WaitGroup) {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("user data blocking everything")
		// Signal to the wait group that the work is done 
		respCh <- "user data"
		wg.Done()
}

func fetchUserRecommendations(userID int, respCh chan string, wg *sync.WaitGroup) {
		time.Sleep(120 * time.Millisecond)
		fmt.Println("user recs blocking everything")
		// Signal to the wait group that the work is done 
		respCh <- "user recommendations"
		wg.Done()
}

func fetchUserLikes(userID int, respCh chan string, wg *sync.WaitGroup) {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("user likes blocking everything")

		respCh <- "user likes"
		// Signal to the wait group that the work is done 
		wg.Done()
}
