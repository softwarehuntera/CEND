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

func TestOverlappingDocuments(t *testing.T) {
	actualCollection := New("Test Collection")
	documents := []string{"apple", "apples"}

	for _, doc := range documents {
		actualCollection.DocumentAdd(doc)
	}

	expectedCollection := overlapCollection1()
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Overlapping collections do not match.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}

	actualCollection.DocumentRemove(documents[1])
	expectedCollection = overlapCollection2()
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Overlapping collections do not match after removal.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}

	expectedCollection = overlapCollection1()
	actualCollection.DocumentAdd(documents[1])
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Overlapping collections do not match after remove and add.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}

func overlapCollection1() (*Collection) {
	return &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
			"apple": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"apples": {
				count: 1,
				docIDs: map[int]struct{}{
					2: {},
				},
			},
			"app": {
				count: 2,
				docIDs: map[int]struct{}{
					1: {},
					2: {},
				},
			},
			"ppl": {
				count: 2,
				docIDs: map[int]struct{}{
					1: {},
					2: {},
				},
			},
			"ple": {
				count: 2,
				docIDs: map[int]struct{}{
					1: {},
					2: {},
				},
			},
			"les": {
				count: 1,
				docIDs: map[int]struct{}{
					2: {},
				},
			},
		},
		documents: &map[int]Document{
			1: {doc: "apple"},
			2: {doc: "apples"},
		},
	}
}

func overlapCollection2() (*Collection) {
	return &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
			"apple": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"app": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"ppl": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"ple": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
		},
		documents: &map[int]Document{
			1: {doc: "apple"},
		},
	}
}

func TestOverlappingDocumentsWithDuplicateTrigrams(t *testing.T) {
	actualCollection := New("Test Collection")
	documents := []string{"apple", "apples", "cargo cart"}

	for _, doc := range documents {
		actualCollection.DocumentAdd(doc)
	}

	// Define the expected collection state after adding "cargo cart"
	expectedCollection := overlapCollection3()
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Collections do not match after adding 'cargo cart'.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}

// overlapCollection3 defines the expected state after adding "apple", "apples", and "cargo cart"
func overlapCollection3() *Collection {
	return &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
			"apple": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"apples": {
				count: 1,
				docIDs: map[int]struct{}{
					2: {},
				},
			},
			"app": {
				count: 2,
				docIDs: map[int]struct{}{
					1: {},
					2: {},
				},
			},
			"ppl": {
				count: 2,
				docIDs: map[int]struct{}{
					1: {},
					2: {},
				},
			},
			"ple": {
				count: 2,
				docIDs: map[int]struct{}{
					1: {},
					2: {},
				},
			},
			"les": {
				count: 1,
				docIDs: map[int]struct{}{
					2: {},
				},
			},
			"car": {
				count: 1,
				docIDs: map[int]struct{}{
					3: {},
				},
			},
			"arg": {
				count: 1,
				docIDs: map[int]struct{}{
					3: {},
				},
			},
			"rgo": {
				count: 1,
				docIDs: map[int]struct{}{
					3: {},
				},
			},
			"go ": {
				count: 1,
				docIDs: map[int]struct{}{
					3: {},
				},
			},
			" ca": {
				count: 1,
				docIDs: map[int]struct{}{
					3: {},
				},
			},
			"art": {
				count: 1,
				docIDs: map[int]struct{}{
					3: {},
				},
			},
		},
		documents: &map[int]Document{
			1: {doc: "apple"},
			2: {doc: "apples"},
			3: {doc: "cargo cart"},
		},
	}
}


func TestAddRemoveSingleDocument(t *testing.T) {
	actualCollection := New("Test Collection")

	// Add a single document
	doc := "apple"
	actualCollection.DocumentAdd(doc)

	// Verify document exists in the collection
	if !actualCollection.DocumentExists(doc) {
		t.Errorf("Document %s should exist in the collection", doc)
	}

	// Remove the document
	actualCollection.DocumentRemove(doc)

	// Verify document no longer exists in the collection
	if actualCollection.DocumentExists(doc) {
		t.Errorf("Document %s should not exist in the collection after removal", doc)
	}
	expectedCollection := &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
		},
		documents: &map[int]Document{
		},
	}
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Collections do not match.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}


func TestAddDuplicateDocument(t *testing.T) {
	actualCollection := New("Test Collection")

	// Add the document twice
	doc := "banana"
	actualCollection.DocumentAdd(doc)
	actualCollection.DocumentAdd(doc)

	expectedCollection := &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
			"ban": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"ana": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"nan": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
		},
		documents: &map[int]Document{
			1: {doc: "banana"},
		},
	}

	// Ensure it is added only once
	if len(*actualCollection.documents) != 1 {
		t.Errorf("Expected only 1 document, got %d", len(*actualCollection.documents))
	}
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Collections do not match.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}


func TestAddLengthNDocument(t *testing.T) {
	actualCollection := New("Test Collection")

	// Add the document twice
	doc := "ban"
	actualCollection.DocumentAdd(doc)

	expectedCollection := &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
			"ban": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
		},
		documents: &map[int]Document{
			1: {doc: "ban"},
		},
	}

	// Ensure it is added only once
	if len(*actualCollection.documents) != 1 {
		t.Errorf("Expected only 1 document, got %d", len(*actualCollection.documents))
	}
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Collections do not match.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}

func TestAddShortDocument(t *testing.T) {
	actualCollection := New("Test Collection")

	// Add a document shorter than n-gram length
	doc := "hi" // shorter than the 3-character n-gram length
	actualCollection.DocumentAdd(doc)

	// Verify document is added without n-grams
	if !actualCollection.DocumentExists(doc) {
		t.Errorf("Document %s should exist in the collection", doc)
	}
	expectedCollection := &Collection{
		name: "Test Collection",
		ngram: 3,
		lookupTable: &map[string]*DocumentIDs{
			"hi": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
		},
		documents: &map[int]Document{
			1:  {doc: "hi"},
		},
	}
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Collections do not match.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}

func TestAddNormalizedDocument(t *testing.T) {
	actualCollection := New("Test Collection with Normalization")

	// Document with mixed case and accents
	doc := "HÉllo"
	normalizedDoc := "hello" // Expected normalized version (lowercased, accents removed)

	// Add document with normalization
	actualCollection.DocumentAdd(doc)

	// Verify normalized document is added
	if !actualCollection.DocumentExists(normalizedDoc) {
		t.Errorf("Document %s (normalized to %s) should exist in the collection", doc, normalizedDoc)
	}

	// Expected collection after adding normalized document
	expectedCollection := &Collection{
		name:    "Test Collection with Normalization",
		ngram:   3,
		lookupTable: &map[string]*DocumentIDs{
			"hello": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"hel": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"ell": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
			"llo": {
				count: 1,
				docIDs: map[int]struct{}{
					1: {},
				},
			},
		},
		documents: &map[int]Document{
			1: {doc: normalizedDoc},
		},
	}

	// Check that actual collection matches expected collection
	if !Equal(actualCollection, expectedCollection) {
		t.Errorf("Collections do not match.\nExpected: %+v\nGot: %+v", expectedCollection, actualCollection)
	}
}