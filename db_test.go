package poctestcontainers

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupDB(t *testing.T) (*sql.DB, error) {
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
		return nil, err
	}
	defer postgresContainer.Terminate(t.Context())

	host, err := postgresContainer.Host(t.Context())
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
		return nil, err
	}

	port, err := postgresContainer.MappedPort(t.Context(), "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
		return nil, err
	}

	db, err := newDB(host, port.Int(), "user", "password", "test")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
		return nil, err
	}

	return db, nil
}

func TestInsertTable(t *testing.T) {
	db, err := setupDB(t)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
		return
	}

	content := "Hello, Testcontainers!"
	err = insertPost(db, content)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}
}
