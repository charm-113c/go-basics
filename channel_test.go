package main

import (
	"fmt"
	"testing"	
)

func TestAddUser(t *testing.T) {
		// In Go, some data structures are not concurreny-safe,
		// if multiple goroutines access them, data races can occur
		// in which case Go would panic, to prevent unexpected behaviour
		// Maps are one such data structure 
		server := NewServer()
		server.Start()

		for i := 0; i < 10; i ++ {
				// adding to the map here, like server.users[user] = user_i
				// Would result in a race condition, because concurrent access could mess with the map's hash buckets or any of its underlying data
				go func(i int) {
						server.userch <- fmt.Sprintf("user_%d", i)
				}(i)
		}
}

