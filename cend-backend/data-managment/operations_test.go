package datamanagment

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/docker/docker/client"

)

const (
	dbName     = "testdb"
	dbUser     = "testuser"
	dbPassword = "testpass"
	dbPort     = "54322"
)






func strPointer(s string) (*string) {
	return &s
}
func tableData() ([]Name) {
	return []Name{
		{Name: "A", PrimaryName: strPointer("B")},
		{Name: "B", PrimaryName: strPointer("B")},
		{Name: "C", PrimaryName: strPointer("D")},
		{Name: "D", PrimaryName: strPointer("D")},
		{Name: "E", PrimaryName: strPointer("D")},

		{Name: "A", PrimaryName: strPointer("B")},
		{Name: "B", PrimaryName: strPointer("B")},
		{Name: "C", PrimaryName: strPointer("D")},
		{Name: "D", PrimaryName: strPointer("D")},
		{Name: "E", PrimaryName: strPointer("D")},
	}
}

func TestDatabaseContainer(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}

	dbConfig := DatabaseConfig{Name: dbName, User: dbUser, Password: dbPassword, Port: dbPort}
	containerID, db, err := setupDatabaseContainer(ctx, cli, dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up database container: %v", err)
	}
	defer teardownDatabaseContainer(ctx, cli, containerID)
	defer db.Close()
	defer cli.Close()

	// Create a table in the database
	createTableQuery := `
	CREATE TABLE test_names (
		Name VARCHAR(255) NOT NULL,
		PrimaryName VARCHAR(255),
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data (optional)
	_, err = db.Exec("INSERT INTO test_names (name, email) VALUES ($1, $2)", "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Verify the data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test_names").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 user in the database, got %d", count)
	}
}