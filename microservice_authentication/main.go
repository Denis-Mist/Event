package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"handlers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	fmt.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
