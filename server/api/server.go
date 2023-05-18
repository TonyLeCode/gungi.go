package api

import (
	"database/sql"
	"log"

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
		log.Fatal(err)
		return err
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Postgres Connected...")

	dbs.PostgresDB = conn

	return nil
}
func (dbs *DBConn) RedisConnect() error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println("Redis Connected...")

	dbs.RedisClient = redisClient

	return nil
}
