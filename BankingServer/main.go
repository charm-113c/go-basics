package main

import (
	"log"
)

func main() {
	log.Println("Shall we dance?")

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	if err = store.Init(); err != nil {
		log.Fatal("Could not initialize DB:", err)
	}

	server := newAPIServer(":3000", store)

	go server.Run()

	select {}
}
