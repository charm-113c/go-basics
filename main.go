package main

import (
	"fmt"
	// "time"
)

type Server struct {
		users map[string] string 
		// add channels to provide concurrency-safety
		userch chan string
		quitch chan struct{}
}

func NewServer() *Server {
		return &Server{
				users: make(map[string]string),
				userch: make(chan string),
				quitch: make(chan struct{}), // struct attrs take up no space until init
		}
}

func (s *Server) Start() {
		// The first part of the pattern to aid in concurrency-safety, among other things
		// (a sort of builder pattern)
		go s.loop()
}

func (s *Server) loop() {
		// This pattern avoids race conditions, all the while avoiding locking the users map
		// This is essentially our first use case of channels
		// Note, however, how it is unbuffered
		// And that's okay, because this down here concurrently reads from it, always
		out:
		for {
				// user := <- s.userch
				// s.users[user] = user
				
				// But in terms of loop, main loops are usually done this way:
				select {
				case msg := <- s.userch:
						fmt.Println(msg)
				case <- s.quitch:
						// two ways to gracefully terminate the server 
						// return
						// or 
						break out // that's right, you can name loops
						// and you can even continue/skip them 
						// but be careful, this is done on purpose as select is a loop of 
						// its own
				}
		}
}

func main() {
		userch := make(chan string, 1)

		userch <- "Bob" // blocking here, it expects a simulatneous read
		// Channels always block when their buffer is full, be extra careful about that
		// An unbuffered channel be considered as having size 0
		// Unless read from concurrently, it will always block; it's a design choice,
		// to make programmers think about concurrency, things are different for buffers >= 1
		// go func() {
		// 		time.Sleep(1 * time.Second)
		// 		userch <- "Bob"
		// }()
		// Now it won't block, goroutines don't block, no matter how long you wait

		user := <- userch

		fmt.Println("Channel 1:", user)

		// But if the channel isn't full
		otherch := make(chan string, 2)
		// It won't block
		otherch <- "Alice"
		// ["Alice", ""]
		otherch <- "Charles"
		// ["Alice", "Charles"]
		// It won't block here, because large buffers aren't as tightly coupled with 
		// consuming operations as small buffers (e.g. consumers might be slower)
		// otherch <- "David" would block, it would be a deadlock unless measures are taken
		user2 := <- otherch
		// ["", ""]
		fmt.Println("Channel 2:", user2)
}

func sendMessage(msgCh chan<- string) {
		// the above specifies that the chan is send only
		// it won't be able to read, it can only send 
		// msg := <- msgCh will result in compiler complaining
		msgCh <- "Hello!"
}

func readMessage(msgCh <-chan string) {
		// same concept as above, these mechanisms protect Channels
		msg := <- msgCh
		fmt.Println(msg)
}
