package main

import (
	"cend/database/collection"
	"encoding/json"
	"net/http"
)

func getSearchResult(collec *collection.Collection, searchResults []collection.SearchResultScore) []SearchResult {
	dc := collec.GetDocumentCollection()
	results := make([]SearchResult, 0, len(searchResults))
	for _, res := range searchResults {
		doc := dc.Get(res.ID)
		fields := doc.Fields()

		if fields == nil {
			f := make(map[string]string)
			fields = &f
		}
		curSearchRes := SearchResult{
			Document: res.Document,
			Score: res.Score,
			Id: res.ID,
			Fields: fields,
			IsPreferred: doc.Preferred(),
			PreferredDocuments: doc.PreferredDocuments(),
		}

		results = append(results, curSearchRes)
	}
	return results
}

type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code string, message string, details string) {
    response := ErrorResponse{
        Status:  status,
        Message: message,
        Code:    code,
        Details: details,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}