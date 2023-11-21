package api

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConn struct {
	Conn *pgxpool.Pool
}

// Create a connection to database
func (db *DBConn) PostgresConnect(dbSource string) error {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err = conn.Ping(ctx); err != nil {
		return err
	}

	db.Conn = conn

	return nil
}
