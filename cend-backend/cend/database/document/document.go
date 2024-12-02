package document

import "fmt"

// TODO: Implement document index store


// DocumentIDs stores the locations of documents containing a specific token.
type DocumentIDs struct {
	Count int // number of times the token is found across documents
	Ids map[int]struct{} // set of document IDs where the token appears.
}

// Remove removes the specified docID from DocumentIDs.
func (d *DocumentIDs) Remove(docID int) error {
	// Check if docID exists and delete it if found
	if _, exists := d.Ids[docID]; exists {
		delete(d.Ids, docID)
		d.Count--
		
		if d.Count < 0 {
			return fmt.Errorf("docCount is negative")
		}
		return nil
	}

	return fmt.Errorf("docID %d not found", docID)
}

// Add adds a new docID to DocumentIDs.
func (IDs *DocumentIDs) Add(docID int) error {
	if _, exists := IDs.Ids[docID]; exists {
		return fmt.Errorf("docID %d already exists", docID)
	}
	IDs.Count++
	IDs.Ids[docID] = struct{}{}
	return nil
}


