package collection

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"

	"cend/database/collection/documents"
)

// TODOs:
// Test add & delete operations
// Add lower-case filtration before token storage
// Add search functionality

// DocumentIDs stores the locations of documents containing a specific token.
type DocumentIDs struct {
	count int // number of times the token is found across documents
	docIDs map[int]struct{} // set of document IDs where the token appears.
}

// Collection represents a collection of documents and provides methods
// for managing tokenized entries and tracking document locations.
type Collection struct {
	Path		 string
	name         string
	ngram		int
	lookupTable  *map[string]*DocumentIDs
	documents    *documents.DocumentCollection
}



type SearchResultScore struct {
	ID int `json:"id"`
	Document string  `json:"document"`
	Score    float64 `json:"score"`
}

// New creates and returns a new Collection with the specified name.
func New(name, path string) *Collection {

	return &Collection{
		Path:		path,
		name:        name,
		ngram:		3,
		lookupTable: &map[string]*DocumentIDs{},
		documents:   documents.NewDocumentCollection(),
	}
}

func (c *Collection) GetDocumentCollection() *documents.DocumentCollection {
	return c.documents
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

// removeDocID removes the specified docID from DocumentIDs.
func (IDs *DocumentIDs) removeDocID(docID int) error {

	// Check if docID exists and delete it if found
	if _, exists := IDs.docIDs[docID]; exists {
		delete(IDs.docIDs, docID)
		IDs.count--
		
		if IDs.count < 0 {
			return fmt.Errorf("docCount is negative")
		}
		return nil
	}

	return fmt.Errorf("docID %d not found", docID)
}

// addDocID adds a new docID to DocumentIDs.
func (IDs *DocumentIDs) addDocID(docID int) error {
	if _, exists := IDs.docIDs[docID]; exists {
		return fmt.Errorf("docID %d already exists", docID)
	}
	IDs.count++
	IDs.docIDs[docID] = struct{}{}
	return nil
}

// tableAdd adds the token and associated docID to the collection’s
// lookupTable.
func (c *Collection) tableAdd(token string, docID int) {
	if ids, exists := (*c.lookupTable)[token]; exists {
		ids.addDocID(docID)
	} else {
		(*c.lookupTable)[token] = &DocumentIDs{
			count: 1,
			docIDs:   make(map[int]struct{}),
		}
		(*c.lookupTable)[token].docIDs[docID] = struct{}{}
	}
}

// tableRemove removes the specified docID from the DocumentIDs
// entry for the token in the lookupTable.
func (c *Collection) tableRemove(token string, docID int) error {
	ids, exists := (*c.lookupTable)[token]
	if !exists {
		return fmt.Errorf("token not found")
	}
	ids.removeDocID(docID)

	if len(ids.docIDs) == 0 {
		delete(*c.lookupTable, token)
	} else {
		(*c.lookupTable)[token] = ids
	}
	return nil
	
}

// DocumentID retrieves the ID of a document if it exists in the
// collection.
func (c *Collection) DocumentID(document string) *int {
	if _, exists := (*c.lookupTable)[document]; !exists {
		return nil
	}
	ids := (*c.lookupTable)[document]
	for docID := range ids.docIDs {
		actualDocument := (*c.documents).Get(docID)
		if actualDocument.String() == document {
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
func (c *Collection) DocumentAdd(document string) error {

	if c.DocumentExists(document) {
		return fmt.Errorf("cannot add document that already exists: document=%s", document)
	}
	normalizedDocument := stringNormalize(document)
	x := nGramFrequency(normalizedDocument, c.ngram)

	docID := c.documents.AddDocumentFromStr(normalizedDocument)
	d := c.documents.Get(docID)
	d.SetTokenFrequency(&x)

	if len(normalizedDocument) != c.ngram {
		c.tableAdd(normalizedDocument, docID)
	}
	ngrams := nGramSet(document, c.ngram)

    for ngram := range ngrams {
        c.tableAdd(ngram, docID)
    }
	return nil
}



// DocumentRemove removes a document from the collection. If the
// document exists, it is removed from documents and its associated
// tokens are removed from the lookupTable.
func (c *Collection) DocumentRemove(docId int) error {

	doc := c.documents.Get(docId)
	docStr := doc.String()
	c.documents.RemoveDocument(docId)
	
	if len(docStr) != c.ngram {
		c.tableRemove(docStr, docId)
	}
	ngrams := nGramSet(docStr, c.ngram)

	for ngram := range ngrams {
		c.tableRemove(ngram, docId)
	}
	return nil
}

// DocumentList retrieves a list of documents from the collection.
func (c *Collection) DocumentList() []string {
	return c.documents.DocumentList()
}

// RelevantDocumentIDs returns a set of document IDs that contain at least one n-gram from the provided document.
func (c *Collection) RelevantDocumentIDs(document string) map[int]struct{} {
	documentIDs := make(map[int]struct{})
	ngrams := nGramSet(document, c.ngram)
	for ngram := range ngrams {
		if ids, exists := (*c.lookupTable)[ngram]; exists {
			for docID := range ids.docIDs {
				documentIDs[docID] = struct{}{}
			}
		}
	}
	return documentIDs
}

