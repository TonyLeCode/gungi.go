package api

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Connection struct for receiver functions
type DBConn struct {
	conn *sql.DB
}

// Create a connection to database
func NewConnection(dbSource string) (*DBConn, error) {
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB Connected...")

	return &DBConn{conn: conn}, nil
}
