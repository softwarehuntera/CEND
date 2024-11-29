package collection

import (
	"fmt"
	"cend/database/tokenizer"
	"cend/database/document"
)

// TODOs:
// Test add & delete operations
// Add lower-case filtration before token storage
// Add search functionality



// Collection represents a collection of documents and provides methods
// for managing tokenized entries and tracking document locations.
type Collection struct {
	Path		 string
	name         string
	ngram		int
	lookupTable  *map[string]*document.DocumentIDs
	documents    *map[int]Document
}

type Document struct {
	doc string
	tokenFrequency map[string]int
}

type SearchResult struct {
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
		lookupTable: &map[string]*document.DocumentIDs{},
		documents:   &map[int]Document{},
	}
}



// tableAdd adds the token and associated docID to the collection’s
// lookupTable.
func (c *Collection) tableAdd(token string, docID int) {
	if ids, exists := (*c.lookupTable)[token]; exists {
		ids.Add(docID)
	} else {
		(*c.lookupTable)[token] = &document.DocumentIDs{
			Count: 1,
			Ids:   make(map[int]struct{}),
		}
		(*c.lookupTable)[token].Ids[docID] = struct{}{}
	}
}

// tableRemove removes the specified docID from the DocumentIDs
// entry for the token in the lookupTable.
func (c *Collection) tableRemove(token string, docID int) error {
	ids, exists := (*c.lookupTable)[token]
	if !exists {
		return fmt.Errorf("token not found")
	}
	ids.Remove(docID)

	if len(ids.Ids) == 0 {
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
	for docID := range ids.Ids {
		actualDocument := (*c.documents)[docID]
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
func (c *Collection) DocumentAdd(document string) error {
	docID := len(*c.documents) + 1

	if c.DocumentExists(document) {
		return fmt.Errorf("cannot add document that already exists: document=%s", document)
	}
	normalizedDocument := tokenizer.StringNormalize(document)
	(*c.documents)[docID] = Document{doc: normalizedDocument, tokenFrequency: tokenizer.NGramFrequency(normalizedDocument, c.ngram)}

	if len(normalizedDocument) != c.ngram {
		c.tableAdd(normalizedDocument, docID)
	}
	ngrams := tokenizer.NGramSet(document, c.ngram)

    for ngram := range ngrams {
        c.tableAdd(ngram, docID)
    }
	return nil
}

// DocumentRemove removes a document from the collection. If the
// document exists, it is removed from documents and its associated
// tokens are removed from the lookupTable.
func (c *Collection) DocumentRemove(document string) error {
	docIDptr := c.DocumentID(document)
	if docIDptr == nil {
		return fmt.Errorf("cannot remove document that does not exist. document=%s", document)
	}

	docID := *docIDptr

	delete(*c.documents, docID)
	
	if len(document) != c.ngram {
		c.tableRemove(document, docID)
	}
	ngrams := tokenizer.NGramSet(document, c.ngram)

	for ngram := range ngrams {
		c.tableRemove(ngram, docID)
	}
	return nil
}

// DocumentList retrieves a list of documents from the collection.
func (c *Collection) DocumentList() []string {
	documents := []string{}
	for _, value := range *c.documents {
		documents = append(documents, value.doc)
	}
	return documents
}

// RelevantDocumentIDs returns a set of document IDs that contain at least one n-gram from the provided document.
func (c *Collection) RelevantDocumentIDs(document string) map[int]struct{} {
	documentIDs := make(map[int]struct{})
	ngrams := tokenizer.NGramSet(document, c.ngram)
	for ngram := range ngrams {
		if ids, exists := (*c.lookupTable)[ngram]; exists {
			for docID := range ids.Ids {
				documentIDs[docID] = struct{}{}
			}
		}
	}
	return documentIDs
}

