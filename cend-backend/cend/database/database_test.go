package database

import (
	"fmt"
	"testing"
)

func TestDatabaseImplementation(t *testing.T) {
	// Step 1: Create a test database
	db := New("test-db")

	// Step 2: Add a collection
	db.AddCollection("docs")

	// Step 3: Retrieve the collection
	collection, err := db.GetCollection("docs")
	if err != nil {
		t.Errorf("Error getting collection: %s", err)
	}

	// Step 4: Add sample documents to the collection
	sampleDocs := []string{
		"To install Docker, first update your package manager and install required dependencies",
		"Docker installation requires updating packages and installing system dependencies",
		"Kubernetes requires Docker to be installed and properly configured on all nodes",
		"Before installing any packages, ensure your system is up to date",
		"Configure Docker daemon settings in the daemon.json configuration file",
	}
	for _, doc := range sampleDocs {
		err := collection.DocumentAdd(doc)
		if err != nil {
			t.Errorf("Error adding document: %s", err)
		}
	}

	// Step 5: Retrieve all documents and verify
	// allDocs := collection.DocumentList()

	for _, searchRes := range collection.DocumentSearch("docker") {
		fmt.Printf("Search result: [%v] score: [%v]", searchRes.Document, searchRes.Score)
	}

}
