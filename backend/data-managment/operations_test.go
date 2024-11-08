package datamanagment

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	dbName     = "testdb"
	dbUser     = "testuser"
	dbPassword = "testpass"
	dbPort     = "54322"
)

// setupDatabaseContainer starts a PostgreSQL container and returns the container ID and the database connection.
func setupDatabaseContainer(ctx context.Context, cli *client.Client) (string, *sql.DB, error) {
	// Pull the PostgreSQL image
	_, err := cli.ImagePull(ctx, "docker.io/library/postgres:latest", types.ImagePullOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to pull postgres image: %v", err)
	}

	// Create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "postgres",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			fmt.Sprintf("POSTGRES_USER=%s", dbUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
		},
		ExposedPorts: map[string]struct{}{dbPort: {}},
	}, &container.HostConfig{
		PortBindings: map[string][]container.PortBinding{
			dbPort: {{HostIP: "0.0.0.0", HostPort: dbPort}},
		},
	}, nil, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for the database to be ready
	time.Sleep(5 * time.Second)

	// Connect to the database
	connStr := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return "", nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return resp.ID, db, nil
}

// teardownDatabaseContainer stops and removes the container.
func teardownDatabaseContainer(ctx context.Context, cli *client.Client, containerID string) {
	cli.ContainerStop(ctx, containerID, nil)
	cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

func TestDatabaseContainer(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	containerID, db, err := setupDatabaseContainer(ctx, cli)
	if err != nil {
		t.Fatalf("Failed to set up database container: %v", err)
	}
	defer teardownDatabaseContainer(ctx, cli, containerID)
	defer db.Close()

	// Create a table in the database
	createTableQuery := `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data (optional)
	_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Verify the data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 user in the database, got %d", count)
	}
}