package main

import (
	"encoding/json"
	"cend/database"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
)

type SearchRequest struct {
	Query      string `json:"query"`
	MaxResults int    `json:"maxResults"`
}

type DeleteRequest struct {
	Document *string `json:"document"`
	Id *int `json:"id"`
}

type AddRequest struct {
	Document string `json:"document"`
	Fields   *map[string]string `json:"fields"`
	IsPreferred bool `json:"isPreferred"`
	PreferredDocuments []int `json:"preferredDocuments"`
}

type QueryRequest struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SearchResult struct {
	Document string  `json:"document"`
	Score    float64 `json:"score"`
	Id 	 int     `json:"id"`
	Fields   *map[string]string `json:"fields"`
	IsPreferred bool `json:"isPreferred"`
	PreferredDocuments []int `json:"preferredDocuments"`
}

type DeleteResult struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type AddResult struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type QueryResult struct {
	Document string  `json:"document"`
	Id 	 int     `json:"id"`
	Fields   *map[string]string `json:"fields"`
	IsPreferred bool `json:"isPreferred"`
	PreferredDocuments []int `json:"preferredDocuments"`
}

type GetRequest struct {
	Ids []int `json:"ids"`
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

		// Convert SearchResultScore to SearchResult
		formattedSearchResult := getSearchResult(docs, searchResults)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(formattedSearchResult)
	}
}

func queryHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            writeError(w, 
                http.StatusMethodNotAllowed,
                "METHOD_NOT_ALLOWED",
                "Only POST method is allowed",
                "Use POST to retrieve documents",
            )
            return
        }

		var req QueryRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w,
                http.StatusBadRequest,
                "INVALID_JSON",
                "Invalid request body",
                err.Error(),
            )
            return
        }
		fmt.Printf("Incoming search request: %v", req)
		fmt.Printf("Query: %v", req)

		if req.Min < 0 || req.Max < req.Min {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid range: min must be >= 0 and max must be >= min", "Error: Invalid Range")
			return
		}

		docs, err := db.GetCollection("docs")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting collection", err.Error())
			return
		}

		docCollec := docs.GetDocumentCollection()
		docList, err := docCollec.GetDocuments(req.Min, req.Max)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting documents", err.Error())
			return
		}
		fmt.Printf("Document List: %v", docList)

		// convert documents to query results because document fields are private and cannot be JSON encoded
		queryResults := make([]QueryResult, 0, len(docList))
		for _, doc := range docList {
			if doc == nil {
				writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting documents", "Document not found")
				return
			}
			queryResults = append(queryResults, documentToQueryResult(doc))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(queryResults)
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
		if req.Document == "" {
			http.Error(w, "Document was empty", http.StatusBadRequest)
			return
		}
		docs.DocumentAdd(req.Document)
		docID := docs.DocumentID(req.Document)
		if docID == nil {
			http.Error(w, "Error adding document", http.StatusInternalServerError)
			return
		}
		docCollec := docs.GetDocumentCollection()
		if docCollec == nil {
			writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting documents", "Document collection not found")
			return
		}

		doc := docCollec.Get(*docID)
		if err := doc.AddFields(req.Fields); err != nil {
			http.Error(w, "Fields invalid", http.StatusBadRequest)
			return
		}
		doc.SetPreferred(req.IsPreferred)
		doc.SetPreferredDocuments(req.PreferredDocuments)


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
		if req.Id == nil {
			if req.Document == nil {
				writeError(w, http.StatusBadRequest, "INPUT_ERROR", "Must specify an Id or Document string.", "Invalid input.")
				return
			}
			docId := docs.DocumentID(*req.Document)
			if docId == nil {
				writeError(w, http.StatusInternalServerError, "NOT_FOUND", "Document not found", "Not found.")
				return
			}
			req.Id = docId

		}
		err = docs.DocumentRemove(*req.Id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error removing document: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("200")
	}
}

func getHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            writeError(w, 
                http.StatusMethodNotAllowed,
                "METHOD_NOT_ALLOWED",
                "Only POST method is allowed",
                "Use POST to retrieve documents",
            )
            return
        }

		var req GetRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w,
                http.StatusBadRequest,
                "INVALID_JSON",
                "Invalid request body",
                err.Error(),
            )
            return
        }
		fmt.Printf("Incoming get request: %v", req)

		docs, err := db.GetCollection("docs")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting collection", err.Error())
			return
		}

		docCollec := docs.GetDocumentCollection()
		if docCollec == nil {
			writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting documents", "Document collection not found")
			return
		}
		
		queryResults := make([]QueryResult, 0, len(req.Ids))
		for _, reqId := range req.Ids {
			doc := docCollec.Get(reqId)
			if doc == nil {
				writeError(w, http.StatusInternalServerError, "INERNAL_ERROR", "Error getting documents", "Document not found")
				return
			}
			queryResults = append(queryResults, documentToQueryResult(doc))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(queryResults)
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