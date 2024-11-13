package collection

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// TODOs:
// Test add & delete operations
// Add lower-case filtration before term storage
// Add search functionality

// DocumentLocations stores the locations of documents containing a specific token.
type DocumentLocations struct {
	frequency int // number of times the token is found across documents
	docIDs map[int]struct{} // set of document IDs where the token appears.
}

// Collection represents a collection of documents and provides methods
// for managing tokenized entries and tracking document locations.
type Collection struct {
	name         string
	ngram		int
	lookupTable  *map[string]*DocumentLocations
	documents    *map[int]Document
}

type Document struct {
	doc string
	termFrequency map[string]int
}

// New creates and returns a new Collection with the specified name.
func New(name string) *Collection {
	return &Collection{
		name:        name,
		ngram:		3,
		lookupTable: &map[string]*DocumentLocations{},
		documents:   &map[int]Document{},
	}
}

func stringNormalize(s string) string {
	// TODO: Remove stop words
	
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Normalize to remove accents (NFD form splits characters from accents)
	s = norm.NFD.String(s)
	
	// Remove diacritics
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // Mn category is for non-spacing marks
			return -1
		}
		return r
	}, s)
	
	// Remove punctuation and special characters, retain spaces and alphanumerics
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	s = b.String()
	
	// Trim spaces
	s = strings.TrimSpace(s)
	
	// Replace multiple spaces with a single space
	s = strings.Join(strings.Fields(s), " ")
	
	return s
}

// nGrams generates a slice of n-grams from the provided document string.
func nGrams(document string, n int) []string {
	document = stringNormalize(document)
	ngrams := []string{}
	documentLength := len(document)
	for i := 0; i <= documentLength-n; i++ {
		ngram := document[i : i+n]
		ngrams = append(ngrams, ngram)
	}
	return ngrams
}

func nGramSet(document string, n int) map[string]struct{} {
	nGrams := nGrams(document, n)
	nGramSet := make(map[string]struct{})
	for _, nGram := range nGrams {
		nGramSet[nGram] = struct{}{}
	}
	return nGramSet
}

func nGramFrequency(document string, n int) map[string]int {
	nGrams := nGrams(document, n)
	nGramFrequency := make(map[string]int)
	for _, nGram := range nGrams {
		nGramFrequency[nGram]++
	}
	return nGramFrequency
}

// removeDocID removes the specified docID from DocumentLocations.
func (dl *DocumentLocations) removeDocID(docID int) error {

	// Check if docID exists and delete it if found
	if _, exists := dl.docIDs[docID]; exists {
		delete(dl.docIDs, docID)
		dl.frequency--
		
		if dl.frequency < 0 {
			return fmt.Errorf("frequency is negative")
		}
		return nil
	}

	return fmt.Errorf("docID %d not found", docID)
}

// addDocID adds a new docID to DocumentLocations.
func (dl *DocumentLocations) addDocID(docID int) error {
	if _, exists := dl.docIDs[docID]; exists {
		return fmt.Errorf("docID %d already exists", docID)
	}
	dl.frequency++
	dl.docIDs[docID] = struct{}{}
	return nil
}

// tableAdd adds the token and associated docID to the collection’s
// lookupTable.
func (c *Collection) tableAdd(token string, docID int) {
	if locations, exists := (*c.lookupTable)[token]; exists {
		locations.addDocID(docID)
	} else {
		(*c.lookupTable)[token] = &DocumentLocations{
			frequency: 1,
			docIDs:   make(map[int]struct{}),
		}
		(*c.lookupTable)[token].docIDs[docID] = struct{}{}
	}
}

// tableRemove removes the specified docID from the DocumentLocations
// entry for the token in the lookupTable.
func (c *Collection) tableRemove(token string, docID int) error {
	locations, exists := (*c.lookupTable)[token]
	if !exists {
		return fmt.Errorf("token not found")
	}
	locations.removeDocID(docID)

	if len(locations.docIDs) == 0 {
		delete(*c.lookupTable, token)
	} else {
		(*c.lookupTable)[token] = locations
	}
	return nil
	
}

// DocumentID retrieves the ID of a document if it exists in the
// collection.
func (c *Collection) DocumentID(document string) *int {
	if _, exists := (*c.lookupTable)[document]; !exists {
		LogInfo("Document not found in lookup table.")
		return nil
	}
	locations := (*c.lookupTable)[document]
	LogInfo(fmt.Sprintf("locations: %v", locations))
	LogInfo(fmt.Sprintf("documents: %v", *c.documents))
	for docID := range locations.docIDs {
		actualDocument := (*c.documents)[docID]
		LogInfo(fmt.Sprintf("Actual document %s", actualDocument))
		if actualDocument.doc == document {
			return &docID
		}
	}
	return nil
}

// DocumentExists returns true if the document exists, otherwise false.
func (c *Collection) DocumentExists(document string) bool {
	return c.DocumentID(document) != nil
}

// DocumentAdd adds a document; its n-grams are tokenized and stored in the lookupTable.
func (c *Collection) DocumentAdd(document string) {
	docID := len(*c.documents) + 1

	if c.DocumentExists(document) {
		LogInfo(fmt.Sprintf("Cannot add document that already exists. document=%s", document))
		return // do not add duplicate documents because this is a database
	}

	normalizedDocument := stringNormalize(document)
	(*c.documents)[docID] = Document{doc: normalizedDocument, termFrequency: nGramFrequency(normalizedDocument, c.ngram)}
	if len(normalizedDocument) != c.ngram {
		c.tableAdd(normalizedDocument, docID)
	}
	ngrams := nGramSet(document, c.ngram)

    for ngram := range ngrams {
        c.tableAdd(ngram, docID)
    }
}

// DocumentRemove removes a document from the collection. If the
// document exists, it is removed from documents and its associated
// tokens are removed from the lookupTable.
func (c *Collection) DocumentRemove(document string) {
	docIDptr := c.DocumentID(document)
	if docIDptr == nil {
		LogInfo(fmt.Sprintf("Cannot remove document that does not exist. document=%s", document))
		return
	}
	docID := *docIDptr

	delete(*c.documents, docID)
	
	if len(document) != c.ngram {
		c.tableRemove(document, docID)
	}
	ngrams := nGramSet(document, c.ngram)
	for ngram := range ngrams {
		c.tableRemove(ngram, docID)
	}
}

func (c *Collection) IDF(token string) float64 {
	docCount := len(*c.documents)
	if docCount == 0 {
		return 0
	}
	locations, exists := (*c.lookupTable)[token]
	if !exists || len(locations.docIDs) == 0 {
		return 0
	}
	docFrequency := len(locations.docIDs)
	return math.Log(float64(docCount) / float64(docFrequency))
}

// DocumentSearch searches for documents that contain n-grams from the provided document.
// It returns a list of document IDs ordered by TF-IDF relevance scores.
func (c *Collection) DocumentSearch(document string) []int {
	normalizedDocument := stringNormalize(document)
	ngrams := nGramSet(normalizedDocument, c.ngram)

	// Create a map to store the TF-IDF scores for each document.
	docIDScores := make(map[int]float64)

	// Calculate TF-IDF score for each document that contains n-grams from the query
	for ngram := range ngrams {
		idf := c.IDF(ngram)
		if idf == 0 {
			continue
		}
		if locations, exists := (*c.lookupTable)[ngram]; exists {
			for docID := range locations.docIDs {
				// Get the term frequency for this ngram in the document
				documentText := (*c.documents)[docID]
				nGramFreq := documentText.termFrequency
				tf := float64(nGramFreq[ngram])

				// Calculate TF-IDF score for this document and n-gram
				tfIdfScore := tf * idf
				docIDScores[docID] += tfIdfScore
			}
		}
	}

	// Convert the map of document scores to a slice of document IDs
	docIDList := make([]int, 0, len(docIDScores))
	for docID := range docIDScores {
		docIDList = append(docIDList, docID)
	}

	// Sort documents by their TF-IDF scores in descending order
	sort.Slice(docIDList, func(i, j int) bool {
		return docIDScores[docIDList[i]] > docIDScores[docIDList[j]]
	})
	return docIDList
}