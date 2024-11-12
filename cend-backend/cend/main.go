package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func blogHandler(w http.ResponseWriter, r *http.Request) {
	titles := []string{"A", "B", "C", "D"}
	json.NewEncoder(w).Encode(titles)
}

func main() {
	log.Print("Prepare db...")

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", blogHandler)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}