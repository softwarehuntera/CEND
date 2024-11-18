package collection

import (
	"fmt"
	"testing"
)

func TestDocumentSearch(t *testing.T) {
	collection := New("Test Collection", "./test-data/test-collection")

	// Add documents with varying degrees of similarity
	docs := []string{
		"the quick brown fox jumps over the lazy dog",
		"a quick brown fox leaps over a sleepy dog",
		"the slow green turtle crawls under the lazy cat",
		"a fast red fox runs past the lazy dog",
		"this is a webpage about cats and dogs",
		"I cant stop quoting memes",
	}

	for _, doc := range docs {
		collection.DocumentAdd(doc)
	}

	// Search for a document similar to docs[0]
	searchDoc := "sophie keeps saying memes please help"
	results := collection.DocumentSearch(searchDoc)

	// Log the search results with document content
	var logMsg string
	for i, result := range results {
		doc := (*collection.documents)[result.ID].doc
		logMsg += fmt.Sprintf("\nRank %d (Score: %.4f): %s", i+1, result.Score, doc)
	}
	LogInfo("Search Results for: " + searchDoc + logMsg)
	t.Errorf("Test")
}

func TestTechnicalDocumentSearch(t *testing.T) {
	collection := New("Technical Documentation Collection", "./test-data/test-collection")

	// Add technical documentation style content
	docs := []string{
		"To install Docker, first update your package manager and install required dependencies",
		"Docker installation requires updating packages and installing system dependencies",
		"Kubernetes requires Docker to be installed and properly configured on all nodes",
		"Before installing any packages, ensure your system is up to date",
		"Configure Docker daemon settings in the daemon.json configuration file",
	}

	for _, doc := range docs {
		collection.DocumentAdd(doc)
	}

	// Search for installation-related content
	searchDoc := "how to install docker and update packages"
	results := collection.DocumentSearch(searchDoc)

	// Log the search results with document content
	var logMsg string
	for i, result := range results {
		doc := (*collection.documents)[result.ID].doc
		logMsg += fmt.Sprintf("\nRank %d (Score: %.4f): %s", i+1, result.Score, doc)
	}
	LogInfo("Search Results for: " + searchDoc + logMsg)
	t.Errorf("Test")
}