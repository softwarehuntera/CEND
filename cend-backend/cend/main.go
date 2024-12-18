package main

import (
	"bufio"
	"cend/database"
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
	Query      string `json:"query"`
	MaxResults int    `json:"maxResults"`
}

type DeleteRequest struct {
	Document string `json:"document"`
}

type AddRequest struct {
	Document string `json:"document"`
}

type GetDocsRequest struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SearchResult struct {
	Document string  `json:"document"`
	Score    float64 `json:"score"`
	Id 	 int     `json:"id"`
	Fields   map[string]string `json:"fields"`
	IsPreferred bool `json:"isPreferred"`
	PreferredDocuments []int `json:"preferredDocuments"`
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
		fmt.Printf("Incoming search request: %v", req)
		fmt.Printf("Query: %v", req.Query)
		// Get the docs collection
		docs, err := db.GetCollection("docs")
		if err != nil {
			http.Error(w, "Error getting collection", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Collection: %v", docs.DocumentList())
		searchResults := docs.DocumentSearch(req.Query)
		fmt.Printf("Search results: %v", searchResults)

		// Apply maxResults limit if specified and greater than 0
		if req.MaxResults > 0 && len(searchResults) > req.MaxResults {
			searchResults = searchResults[:req.MaxResults]
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(searchResults)
	}
}

func getDocsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req GetDocsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate range
		if req.Min < 0 || req.Max < req.Min {
			http.Error(w, "Invalid range: min must be >= 0 and max must be >= min", http.StatusBadRequest)
			return
		}

		// Get the docs collection
		docs, err := db.GetCollection("docs")
		if err != nil {
			http.Error(w, "Error getting collection", http.StatusInternalServerError)
			return
		}

		docList := docs.DocumentList()
		
		// Ensure max doesn't exceed document list length
		if req.Max >= len(docList) {
			req.Max = len(docList) - 1
		}

		// Return the slice of documents within the requested range
		rangedDocs := docList[req.Min:req.Max+1]

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rangedDocs)
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
			http.Error(w, fmt.Sprintf("Error removing document: %v", err), http.StatusInternalServerError)
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
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	basicInfo := map[string]string{
		"greeting": "Welcome to Cend!",
		"name":     "Cend API",
		"version":  "1.0.0",
	}
	json.NewEncoder(w).Encode(basicInfo)
}

func readTestData() []string {
	filePath := "../test-data/names.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return []string{}
	}
	defer file.Close()

	// Slice to store names
	var names []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Append each line to the names slice
		names = append(names, scanner.Text())
	}
	return names
}

func main() {
	log.Print("Preparing database...")

	// Create database and collection
	db := database.New("test-db")
	db.AddCollection("docs")
	collection, _ := db.GetCollection("docs")

	// Add sample documents
	// sampleDocs := []string{
	// 	"To install Docker, first update your package manager and install required dependencies",
	// 	"Docker installation requires updating packages and installing system dependencies",
	// 	"Kubernetes requires Docker to be installed and properly configured on all nodes",
	// 	"Before installing any packages, ensure your system is up to date",
	// 	"Configure Docker daemon settings in the daemon.json configuration file",
	// }
	bigSampleDocs := readTestData()
	log.Printf("Adding %d sample documents...", len(bigSampleDocs))
	for _, doc := range bigSampleDocs {
		collection.DocumentAdd(doc)
	}


	log.Print("Setting up routes...")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/search", searchHandler(db))
	r.HandleFunc("/add", addHandler(db))
	r.HandleFunc("/delete", removeHandler(db))
	r.HandleFunc("/query", queryHandler(db))
	r.HandleFunc("/get-docs", getDocsHandler(db))

	log.Print("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}
