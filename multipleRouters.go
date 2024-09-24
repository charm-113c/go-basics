// Another GPT example on how to use multiple routers
package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "log"
)

// UserRouter handles user-related routes
func UserRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/users", GetUsers).Methods("GET")
    r.HandleFunc("/users/{id}", GetUser).Methods("GET")
    r.HandleFunc("/users", CreateUser).Methods("POST")
    return r
}

// ProductRouter handles product-related routes
func ProductRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/products", GetProducts).Methods("GET")
    r.HandleFunc("/products/{id}", GetProduct).Methods("GET")
    r.HandleFunc("/products", CreateProduct).Methods("POST")
    return r
}

// MainRouter combines all subrouters
func MainRouter() *mux.Router {
    r := mux.NewRouter()
    
    // Subrouters
    r.PathPrefix("/api/users").Handler(UserRouter())
    r.PathPrefix("/api/products").Handler(ProductRouter())
    
    return r
}

// Example handler functions
func GetUsers(w http.ResponseWriter, r *http.Request) {
    // Implementation here
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    // Implementation here
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Implementation here
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
    // Implementation here
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
    // Implementation here
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
    // Implementation here
}

func routerMain() {
    r := MainRouter()
    log.Fatal(http.ListenAndServe(":8080", r))
}

