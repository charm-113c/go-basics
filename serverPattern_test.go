package main

import (
	"fmt"
	"testing"
	"time"
)

type Message struct {
		from string
		payload string
}

type TServer struct{
		msgCh chan Message
		quitCh chan struct{}	
}

func NewTServer() *TServer {
		return &TServer{
				make(chan Message),
				make(chan struct{}),
		}
}

func (s *TServer) startAndListen() {
		out:
		for {
				// An infinite for loop that listens
				// to messages
				select {
				case msg := <- s.msgCh:
						fmt.Printf("Received msg from: %s, content: %s\n", msg.from, msg.payload)
				case <-s.quitCh:
						fmt.Println("Server shutting down gracefully")
						break out
				case <- time.After(2 * time.Second):
						fmt.Println("Two seconds have passed")
				}
		}
}

func sendMsgToServer(msgCh chan<- Message) {
		fmt.Println("Sending message to test server")
		msg := Message{
				from: "User",
				payload: "Hello there",
		}
		msgCh <- msg
}

func shutdownGracefully(quitCh chan<- struct{}) {
		close(quitCh)
}

// The effective main() function
func TestServer(t *testing.T) {
		// One of the most common patterns to start up a server 
		// Basically, the server is a daemon
		s := NewTServer()

		// And here we start the daemon
		go s.startAndListen()

		go func() {
				time.Sleep(1* time.Second)
				sendMsgToServer(s.msgCh)
				// close(s.msgCh)
		}()

		time.Sleep(5 * time.Second)
		shutdownGracefully(s.quitCh)

		// select {
		// 		// An empty select is a way to put the main goroutine
		// 		// to sleep indefinitely. DO NOT do this in prod.
		// 		// Listen for errors or shutdown signals instead
		// }
}
