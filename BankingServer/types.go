package main

import (
	"math/rand"
	"time"
)

// This is the body of the request, as they obviously won't hold IDs, account numbers or time stamps
type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	AccNumber int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

func NewAccount(firstName, lastName string) *Account {
	// Returns randomly generated account
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		AccNumber: int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}
