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
func (dbs *DBConn) RedisConnect(dbSource string) error {
	opt, err := redis.ParseURL(dbSource)
	if err != nil {
		log.Println(err)
		return err
	}
	redisClient := redis.NewClient(opt)
	log.Println("Redis Connected...")

	dbs.RedisClient = redisClient

	return nil
}
