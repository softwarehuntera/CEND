package index

import (
    "fmt"
    "os"
)

// Index represents the document storage index with a path to the directory.
type Index struct {
    path string
}

// New creates a new Index instance with the given path and ensures the directory is created.
func New(path string) (*Index, error) {
    idx := &Index{
        path: path,
    }

    // Create the document storage directory
    err := idx.CreateDocStorageDir()
    if err != nil {
        return nil, err
    }

    return idx, nil
}

// CreateDocStorageDir creates a directory for document storage.
func (i *Index) CreateDocStorageDir() error {
    // Create the directory, including any necessary parents, with read-write-execute permissions for the owner, and read-execute for others
    err := os.MkdirAll(i.path, 0755)
    if err != nil {
        return fmt.Errorf("failed to create directory %s: %w", i.path, err)
    }

    fmt.Println("Directory created at:", i.path)
    return nil
}
