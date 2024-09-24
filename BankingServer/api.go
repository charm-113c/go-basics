package main

import (
	"encoding/json"
	"os"
	"strconv"

	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	// To simplify writing to clients
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

// A custom type to simplify things

// A type to show you how potentially complex things can be simplified
type apiError struct {
	ErrorMsg string
}

// And the decorator function
func httpHandlerDecorator(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// We need to handle the error, and for now we'll simply log and tell the client
			// Get the name of the culprit function
			funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			errMsg := fmt.Sprintf("Error on handler function %s: %v", funcName, err)
			log.Println(errMsg)
			// Then write to client
			if err = WriteJSON(w, http.StatusInternalServerError, "Something went wrong on our side"); err != nil {
				log.Println("Error writing to client:", err)
			}
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func newAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	// Before you listen and serve anything, we need at least one router
	// Ok, theoretically, we don't need it, but practically, we do

	router := mux.NewRouter()

	// But the below methods aren't http handlers: they return an error, which http handlers don't
	// So while we could simply handle the error internally, that creates what is essentially
	// ugly code.
	// Instead, we'll use the decorator pattern: we'll wrap these handlers inside a function (
	// with said function corresponding to the http.handler signature) and handle any error
	// there, once and for all
	router.HandleFunc("/account/{id}", withJWTAuth(httpHandlerDecorator(s.handleGetAccountByID))).Methods("GET")
	router.HandleFunc("/account", httpHandlerDecorator(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account", httpHandlerDecorator(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", httpHandlerDecorator(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/transfer", httpHandlerDecorator(s.handleTransfer)).Methods("POST")

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Println("Error starting up server:", err)
		return
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	// Let's test the error logging first
	// return errors.New("error handling account")

	id, err := readID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	newAccountBody := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(newAccountBody); err != nil {
		return err
	}

	newAccount := NewAccount(newAccountBody.FirstName, newAccountBody.LastName)
	id, err := s.store.CreateAccount(newAccount)
	if err != nil {
		return err
	}

	tokenString, err := createJWT(newAccount, id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, tokenString)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := readID(r)
	if err != nil {
		return err
	}

	if err = s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "Account deleted")
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusAccepted, transferReq)
}

// When you use the same code more than once, it's time to make a function for it
func readID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, fmt.Errorf("provided id: %s is invalid", idStr)
	}

	return id, nil
}

func createJWT(acc *Account, id int) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": acc.AccNumber,
		"sub": id,
	}

	secret := os.Getenv("JWT_TOKEN")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	// THIS LINE RIGHT HERE is the most important part
	// The secret has to be securely inputted from the env variables
	secret := os.Getenv("JWT_TOKEN")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// Let's implement JWTs
func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Calling JWT Auth middleware")

		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w, err)
			return
		}

		if !token.Valid {
			permissionDenied(w, err)
			return
		}

		// Validate JWT subject == requesting user
		claims := token.Claims.(jwt.MapClaims)

		id, err := readID(r)
		if err != nil {
			permissionDenied(w, err)
			return
		}
		
		if id != int(claims["sub"].(float64)) {
			permissionDenied(w, err)
			return
		}

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter, err error) {
	log.Println("Error validating token:", err)
	WriteJSON(w, http.StatusForbidden, apiError{ErrorMsg: "permission denied"})
}
