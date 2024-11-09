package testutils

import (
	"datamanagment"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
)

func setupDockerPostgres() {
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
}
// setupDatabaseContainer starts a PostgreSQL container and returns the container ID and the database connection.
func setupDatabaseContainer(ctx context.Context, cli *client.Client, dbConfig DatabaseConfig) (string, *sql.DB, error) {
	// Pull the PostgreSQL image
	_, err := cli.ImagePull(ctx, "docker.io/library/postgres:latest", image.PullOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to pull postgres image: %v", err)
	}

	port := nat.Port(dbConfig.Port)
	// Create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "postgres",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", dbConfig.Name),
			fmt.Sprintf("POSTGRES_USER=%s", dbConfig.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbConfig.Password),
		},
		ExposedPorts: map[nat.Port]struct{}{port: {}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			port: {{HostIP: "0.0.0.0", HostPort: dbConfig.Port}},
		},
	}, nil, nil, "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", nil, fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for the database to be ready
	time.Sleep(5 * time.Second)

	// Connect to the database
	connStr := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
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

// teardownDatabaseContainer stops and removes the container specified by containerID.
func teardownDatabaseContainer(ctx context.Context, cli *client.Client, containerID string) error {
	// Stop the container if it's running.
	err := cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerID, err)
	}

	// Remove the container.
	err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         true,  // Force removal even if the container is running (stopped first in the previous step).
		RemoveVolumes: true, // Optionally remove volumes associated with the container.
	})
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}

	log.Printf("Successfully removed container %s", containerID)
	return nil
}