package collection

import (
	"math"
	"slices"
)

// DocumentSearch finds similar documents
func (c *Collection) DocumentSearch(searchDoc string) []SearchResult {
	searchVector := c.vectorTFIDF(searchDoc)
   
	searchResult := []SearchResult{}
	for docID := range c.RelevantDocumentIDs(searchDoc) {
		matchDoc := (*c.documents)[docID]
		matchVector := c.vectorTFIDF(matchDoc.doc)
		searchResult = append(searchResult, SearchResult{docID, matchDoc.doc, dotProduct(searchVector, matchVector)})
	}
	sortSearchResult(searchResult)
	return searchResult
}

func dotProduct(v1, v2 map[string]float64) float64 {
	dotProduct := 0.0
	for ngram, searchTFIDF := range v1 {
		if matchTFIDF, exists := v2[ngram]; exists {
			dotProduct += searchTFIDF * matchTFIDF
		}
	}
	return dotProduct
}



func sortSearchResult(searchResult []SearchResult) {
	compareByScoreDesc := func(a, b SearchResult) int {
		if a.Score > b.Score {
			return -1
		}
		if a.Score < b.Score {
			return 1
		}
		return 0
	}
	slices.SortFunc(searchResult, compareByScoreDesc)
}

func (c *Collection) IDF(token string) float64 {
	docCount := len(*c.documents)
	if docCount == 0 {
		return 0
	}
	ids, exists := (*c.lookupTable)[token]
	if !exists || ids.count == 0 {
		return 0
	}
	return math.Log(float64(docCount) / float64(ids.count))
}

func (c *Collection) vectorTFIDF(document string) map[string]float64 {
	docIDptr := c.DocumentID(document)
	var tokenFrequency map[string]int
	if docIDptr == nil {
		tokenFrequency = nGramFrequency(document, c.ngram)
	} else {
		docID := *docIDptr
		tokenFrequency = (*c.documents)[docID].tokenFrequency
	}

	vector := make(map[string]float64)
	var norm float64
	for token, tf := range tokenFrequency {
		idf := c.IDF(token)
		tokenTFIDF := float64(tf) * idf
		vector[token] = tokenTFIDF
		norm += tokenTFIDF * tokenTFIDF
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for token := range vector {
			vector[token] /= norm
		}
	}

	return vector
}