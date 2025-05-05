package main

import (
	"log"
	"net/http"

	"people/internal/db"
	"people/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	database := db.Init()
	h := handlers.Handler{DB: database}

	r := mux.NewRouter()
	r.HandleFunc("/people", h.CreatePerson).Methods("POST")
	r.HandleFunc("/people", h.GetPeople).Methods("GET")
	r.HandleFunc("/people/{id}", h.UpdatePerson).Methods("PUT")
	r.HandleFunc("/people/{id}", h.DeletePerson).Methods("DELETE")

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
