package main

import (
	"bufio"
	"cend/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


func readTestData() []string {
	filePath := "../test-data/names.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return []string{}
	}
	defer file.Close()

	// Slice to store names
	var names []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Append each line to the names slice
		names = append(names, scanner.Text())
	}
	return names
}

func main() {
	log.Print("Preparing database...")

	// Create database and collection
	db := database.New("test-db")
	db.AddCollection("docs")
	collection, _ := db.GetCollection("docs")

	bigSampleDocs := readTestData()
	log.Printf("Adding %d sample documents...", len(bigSampleDocs))
	for _, doc := range bigSampleDocs {
		collection.DocumentAdd(doc)
	}


	log.Print("Setting up routes...")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/search", searchHandler(db))
	r.HandleFunc("/add", addHandler(db))
	r.HandleFunc("/delete", removeHandler(db))
	r.HandleFunc("/query", queryHandler(db))
	r.HandleFunc("/get", getHandler(db))

	log.Print("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}
