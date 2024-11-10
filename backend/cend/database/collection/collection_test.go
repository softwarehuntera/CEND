package collection

import (
	"testing"
)

// TestDocumentAdd tests that documents are correctly added to the collection,
// including adding n-grams to the lookup table.
func TestDocumentAdd(t *testing.T) {
	collection := New("Test Collection")

	// Add a document
	document := "example"
	collection.DocumentAdd(document)

	// Check if the document exists
	if !collection.DocumentExists(document) {
		t.Errorf("Expected document '%s' to exist in the collection after addition", document)
	}

	// Verify document ID retrieval
	docID := collection.DocumentID(document)
	if docID == nil {
		t.Errorf("Expected a valid document ID for document '%s'", document)
	}

	// Verify n-grams are stored in lookup table
	ngrams := nGrams(document, collection.ngram)
	for _, ngram := range ngrams {
		if _, exists := (*collection.lookupTable)[ngram]; !exists {
			t.Errorf("Expected n-gram '%s' to be in the lookup table after adding document '%s'", ngram, document)
		}
	}
}

// TestDocumentRemove tests that documents are correctly removed from the collection,
// including removing associated tokens from the lookup table.
func TestDocumentRemove(t *testing.T) {
	collection := New("Test Collection")

	// Add and then remove a document
	document := "example"
	collection.DocumentAdd(document)
	collection.DocumentRemove(document)

	// Check if the document still exists
	if collection.DocumentExists(document) {
		t.Errorf("Expected document '%s' to be removed from the collection", document)
	}

	// Verify document ID retrieval returns nil
	docID := collection.DocumentID(document)
	if docID != nil {
		t.Errorf("Expected no document ID for document '%s' after removal", document)
	}

	// Verify n-grams are removed from lookup table
	ngrams := nGrams(document, collection.ngram)
	for _, ngram := range ngrams {
		if _, exists := (*collection.lookupTable)[ngram]; exists {
			t.Errorf("Expected n-gram '%s' to be removed from the lookup table after removing document '%s'", ngram, document)
		}
	}
}
