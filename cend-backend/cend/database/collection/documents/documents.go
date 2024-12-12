package documents

type Document struct {
	doc                string
	id                 int
	tokenFrequency     *map[string]int  // Optional
	isPreferred        *bool            // Optional
	preferredDocuments *[]int           // Optional
}

type DocumentCollection struct {
	documents map[int]*Document
}

func NewDocumentCollection() *DocumentCollection {
	return &DocumentCollection{
		documents: make(map[int]*Document),
	}
}

func (dc *DocumentCollection) AddDocument(doc Document) int {
	docID := len(dc.documents) + 1
	dc.documents[docID] = &doc
	return docID
}

func (dc *DocumentCollection) AddDocumentFromStr(doc string) int {
	docID := len(dc.documents) + 1
	dc.documents[docID] = NewDocument(doc, docID, nil, nil, nil)
	return docID
}

func (dc *DocumentCollection) Get(id int) *Document {
	return dc.documents[id]
}

func (dc *DocumentCollection) Length() int {
	return len(dc.documents)
}


func NewDocument(
	doc string,
	id int,
	tokenFrequency *map[string]int,
	isPreferred *bool,
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

func (d *Document) PreferredDocuments() []int {
	return *d.preferredDocuments
}

func (d *Document) SetPreferred(b bool) {
	*d.isPreferred = b
}

