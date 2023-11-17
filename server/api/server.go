package api

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConn struct {
	Conn *sql.DB
}

// Create a connection to database
func (db *DBConn) PostgresConnect(dbSource string) error {
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	db.Conn = conn

	return nil
}
