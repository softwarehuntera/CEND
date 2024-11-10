package collection

import (
	"fmt"
	"log/slog"
)

type DocumentLocations struct {
	frequency int
	docIDs []int // TODO: speedup with map
}

type Collection struct {
	name         string
	ngram		int
	lookupTable  *map[string]DocumentLocations
	documents    *map[int]string
}

func New(name string) *Collection {
	return &Collection{
		name:        name,
		ngram:		3,
		lookupTable: &map[string]DocumentLocations{},
		documents:   &map[int]string{},
	}
}

func nGrams(document string, n int) []string {
	ngrams := []string{}
	documentLength := len(document)
	for i := 0; i <= documentLength-n; i++ {
		ngram := document[i : i+n]
		ngrams = append(ngrams, ngram)
	}
	return ngrams
}

func (dl *DocumentLocations) removeDocID(docID int) error {
	if dl.frequency <= 0 {
		return fmt.Errorf("frequency is already zero or invalid")
	}
	updatedDocIDs := []int{}
	found := false
	for _, id := range dl.docIDs {
		if id != docID {
			updatedDocIDs = append(updatedDocIDs, id)
		} else {
			found = true
		}
	}
	if found {
		dl.docIDs = updatedDocIDs
		dl.frequency--
		return nil
	}
	
	return fmt.Errorf("docID %d not found", docID)
}

func (dl *DocumentLocations) addDocID(docID int) error {
	for existingID, _ := range dl.docIDs {
		if existingID == docID {
			return fmt.Errorf("docID %d already exists", docID)
		}
	}
	dl.frequency++
	dl.docIDs = append(dl.docIDs, docID)
	return nil
}

func (c *Collection) tableAdd(token string, docID int) {
	if locations, exists := (*c.lookupTable)[token]; exists {
		locations.addDocID(docID)
	} else {
		(*c.lookupTable)[token] = DocumentLocations{
			frequency: 1,
			docIDs:    []int{docID},
		}
	}
}

func (c *Collection) tableRemove(token string, docID int) error {
	locations, exists := (*c.lookupTable)[token]
	if ; !exists {
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

func (c *Collection) DocumentID(document string) *int {
	if _, exists := (*c.lookupTable)[document]; !exists {
		return nil
	}
	locations := (*c.lookupTable)[document]
	for docID, _ := range locations.docIDs {
		actualDocument := (*c.documents)[docID]
		if actualDocument == document {
			return &docID
		}
	}
	return nil
}

func (c *Collection) DocumentExists(document string) bool {
	return c.DocumentID(document) != nil
}

func (c *Collection) DocumentAdd(document string) {
	docID := len(*c.documents) + 1

	if c.DocumentExists(document) {
		slog.Info("Cannot add document that already exists.", "document", document)
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

func (c *Collection) DocumentRemove(document string) {
	docIDptr := c.DocumentID(document)
	if docIDptr == nil {
		slog.Info("Cannot remove document that does not exist.", "document", document)
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