package main

import (
	"cend/database"
	"cend/database/collection"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type DeleteRequest struct {
	Document string `json:"document"`
}

type AddRequest struct {
	Document string `json:"document"`
}

type SearchResult struct {
	Document string  `json:"document"`
	Score    float64 `json:"score"`
}

func searchHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get the docs collection
		docs, err := db.GetCollection("docs")
		if err != nil {
			http.Error(w, "Error getting collection", http.StatusInternalServerError)
			return
		}
		searchResults := docs.DocumentSearch(req.Query)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(searchResults)
	}
}

func addHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Get the docs collection
		docs, err := db.GetCollection("docs")
		if err != nil {
			http.Error(w, "Error getting collection", http.StatusInternalServerError)
			return
		}
		err = docs.DocumentAdd(req.Document)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding document: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("200")
	}
}

func removeHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req DeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Get the docs collection
		docs, err := db.GetCollection("docs")
		if err != nil {
			http.Error(w, "Error getting collection", http.StatusInternalServerError)
			return
		}
		err = docs.DocumentRemove(req.Document)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding document: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("200")
	}
}

func queryHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// var req struct{}
		// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
		// 	return
		// }
		
		// Get the docs collection
		docs, err := db.GetCollection("docs")
		if err != nil {
			http.Error(w, "Error getting collection", http.StatusInternalServerError)
			return
		}
		dl := docs.DocumentList()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dl)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	basicInfo := map[string]string {
		"greeting": "Welcome to Cend!",
		"name": "Cend API",
		"version": "1.0.0",
	}
	json.NewEncoder(w).Encode(basicInfo)
}

func main() {
	log.Print("Preparing database...")
	
	// Create database and collection
	db := database.New("test-db")
	docs := collection.New("Documentation Collection")

	// Add sample documents
	sampleDocs := []string{
		"To install Docker, first update your package manager and install required dependencies",
		"Docker installation requires updating packages and installing system dependencies",
		"Kubernetes requires Docker to be installed and properly configured on all nodes",
		"Before installing any packages, ensure your system is up to date",
		"Configure Docker daemon settings in the daemon.json configuration file",
	}
	for _, doc := range sampleDocs {
		docs.DocumentAdd(doc)
	}

	// Add collection to database
	db.AddCollection("docs", docs)
	
	log.Print("Setting up routes...")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/search", searchHandler(db))
	r.HandleFunc("/add", addHandler(db))
	r.HandleFunc("/delete", removeHandler(db))
	r.HandleFunc("/query", queryHandler(db))

	log.Print("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}