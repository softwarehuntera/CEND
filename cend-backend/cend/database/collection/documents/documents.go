package documents

import (
	"fmt"
	"log"
)

type Document struct {
	doc                string
	id                 int
	fields             *map[string]string // Optional
	tokenFrequency     *map[string]int    // Optional
	isPreferred        *bool              // Optional
	preferredDocuments *[]int             // Optional
}

type DocumentCollection struct {
	documents map[int]*Document
}

func NewDocumentCollection() *DocumentCollection {
	return &DocumentCollection{
		documents: make(map[int]*Document),
	}
}

func (dc *DocumentCollection) DocumentList() []string {
	documents := []string{}
	for _, value := range dc.documents {
		documents = append(documents, value.doc)
	}
	return documents
}

func (dc *DocumentCollection) Documents() map[int]*Document {
	return dc.documents
}

func (dc *DocumentCollection) GetDocuments(min, max int) ([]*Document, error) {
	docs := []*Document{}
	dcLen := dc.Length()

	// ensure valid range
	if min < 0 {
		return nil, fmt.Errorf("invalid min: %v", min)
	}
	if max >= dcLen {
		return nil, fmt.Errorf("invalid max: %v", max)
	}

	for _, id := range make([]int, max-min+1) {
		docs = append(docs, dc.Get(id))
	}
	return docs, nil
}

func (dc *DocumentCollection) AddDocument(doc Document) int {
	docID := len(dc.documents) + 1
	dc.documents[docID] = &doc
	if doc.ID() != docID {
		log.Fatalf("Cannot add document with invalid ID");
	}
	return docID
}

func (dc *DocumentCollection) AddDocumentFromStr(doc string) int {
	docID := len(dc.documents) + 1
	dc.documents[docID] = NewDocument(doc, docID, nil, nil, nil, nil)
	return docID
}

func (dc *DocumentCollection) RemoveDocument(id int) {
	delete(dc.documents, id)
}

func (dc *DocumentCollection) Get(id int) *Document {
	doc, exists := dc.documents[id]
	if !exists {
		return nil
	}
	return doc
}

func (dc *DocumentCollection) Length() int {
	return len(dc.documents)
}

func NewDocument(
	doc string,
	id int,
	tokenFrequency *map[string]int,
	isPreferred *bool,
	fields *map[string]string,
	preferredDocuments *[]int,
) *Document {
	return &Document{
		doc:                doc,
		id:                 id,
		tokenFrequency:     tokenFrequency,
		isPreferred:        isPreferred,
		preferredDocuments: preferredDocuments,
	}
}

func (d *Document) String() string {
	return d.doc
}

func (d *Document) TokenFrequency() *map[string]int {
	return d.tokenFrequency
}

func (d *Document) SetTokenFrequency(tf *map[string]int) {
	d.tokenFrequency = tf
}

func (d *Document) ID() int {
	return d.id
}

func (d *Document) Preferred() bool {
	return *d.isPreferred
}

func (d *Document) SetPreferred(b bool) {
	*d.isPreferred = b
}


func (d *Document) PreferredDocuments() []int {
	return *d.preferredDocuments
}

func (d *Document) SetPreferredDocuments(preferredDocs []int) {
	d.preferredDocuments = &preferredDocs
}

func (d *Document) Fields() *map[string]string {
	return d.fields
}

func (d *Document) SetFields(newFields *map[string]string) error {
	if newFields == nil {
		return fmt.Errorf("cannot set fields to nil")
	}
	d.fields = newFields
	return nil
}

func (d *Document) AddFields(newFields *map[string]string) error {
	if newFields == nil {
		return fmt.Errorf("cannot add nil fields")
	}

	// If current fields are nil, initialize them
	if d.fields == nil {
		d.fields = &map[string]string{}
	}

	// Add each new field to the existing map
	for key, value := range *newFields {
		(*d.fields)[key] = value
	}

	return nil
}
