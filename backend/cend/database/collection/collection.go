package collection

import (
	"fmt"
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
	documents    *map[int]string
}

// New creates and returns a new Collection with the specified name.
func New(name string) *Collection {
	return &Collection{
		name:        name,
		ngram:		3,
		lookupTable: &map[string]*DocumentLocations{},
		documents:   &map[int]string{},
	}
}

// nGrams generates a slice of n-grams from the provided document string.
func nGrams(document string, n int) []string {
	ngrams := []string{}
	documentLength := len(document)
	for i := 0; i <= documentLength-n; i++ {
		ngram := document[i : i+n]
		ngrams = append(ngrams, ngram)
	}
	return ngrams
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
		if actualDocument == document {
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

	(*c.documents)[docID] = document
	if len(document) != c.ngram {
		c.tableAdd(document, docID)
	}
	ngrams := nGrams(document, c.ngram)

    for _, ngram := range ngrams {
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
	ngrams := nGrams(document, c.ngram)
	for _, ngram := range ngrams {
		c.tableRemove(ngram, docID)
	}
}