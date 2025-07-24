package poctestcontainers

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInsertTable(t *testing.T) {
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(t.Context())

	host, err := postgresContainer.Host(t.Context())
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(t.Context(), "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	db, err := newDB(host, port.Int(), "user", "password", "test")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	content := "Hello, Testcontainers!"
	err = insertPost(db, content)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}
}

func TestGetData(t *testing.T) {
	ctr, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mitchallen/random-server:latest",
			ExposedPorts: []string{"3100"},
			WaitingFor:   wait.ForLog("random-server:2.1.15 - listening on port 3100!"),
		},
		Started: true,
	})

	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	defer ctr.Terminate(t.Context())

	url, err := ctr.Endpoint(t.Context(), "http")
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	err = getData(url)
	if err != nil {
		t.Fatalf("Failed to get data: %v", err)
	}
}
