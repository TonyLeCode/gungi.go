package api

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DBConn struct {
	PostgresDB  *sql.DB
	RedisClient *redis.Client
}

// Create a connection to database
func (dbs *DBConn) PostgresConnect(dbSource string) error {
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	dbs.PostgresDB = conn

	return nil
}
func (dbs *DBConn) RedisConnect(dbSource string) error {
	opt, err := redis.ParseURL(dbSource)
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(opt)

	dbs.RedisClient = redisClient

	return nil
}
