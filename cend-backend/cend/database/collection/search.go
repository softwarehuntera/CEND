package collection

import (
	"math"
	"sort"
)

// DocumentSearch finds similar documents
func (c *Collection) DocumentSearch(searchDoc string) []docScore {
	searchVector := c.vectorTFIDF(searchDoc)
   
	similarities := make(map[int]float64)
	for docID := range c.RelevantDocumentIDs(searchDoc) {
		matchDoc := (*c.documents)[docID]
		matchVector := c.vectorTFIDF(matchDoc.doc)
		similarities[docID] = dotProduct(searchVector, matchVector)
	}
	matches := sortScores(similarities)
	return matches
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

func sortScores(scores map[int]float64) []docScore {
	
	var sorted []docScore
	for id, score := range scores {
		sorted = append(sorted, docScore{id, score})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})
	return sorted
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