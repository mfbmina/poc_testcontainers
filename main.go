package poctestcontainers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func main() {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "<password>"
	dbname := "<dbname>"

	db, err := newDB(host, port, user, password, dbname)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	insertPost(db, "Hello, world!")
	getData("http://example.com/api/data")
}

func newDB(host string, port int, user, password, dbname string) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Connected!")
	return db, nil
}

func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Table created successfully")
	return nil
}

func insertPost(db *sql.DB, content string) error {
	query := `INSERT INTO posts (content) VALUES ($1);`
	_, err := db.Exec(query, content)
	if err != nil {
		return fmt.Errorf("error inserting post: %w", err)
	}

	log.Println("Post inserted successfully")
	return nil
}

func getData(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error fetching: %v", resp.Status)
	}

	// do something with the response

	return nil
}
