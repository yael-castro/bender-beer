package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultRedisExpirationTime = time.Second * 5
	// defaultRedisDatabase       = 0
)

// REDIS variables
var (
	onceRedis   sync.Once
	redisClient *redis.Client
)

// NewRedisClient returns a singleton of *client.Redis
//
// The following environment variables are required:
//
// REDIS_USER (string)
// REDIS_PASSWORD (string)
// REDIS_HOST (string)
// REDIS_PORT (integer)
// REDIS_DB (integer)
func NewRedisClient() (*redis.Client, error) {
	var err error

	onceRedis.Do(func() {
		redisDB, err := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
		if err != nil {
			return
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Username: os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       int(redisDB),
		})

		statusCmd := redisClient.Ping(context.TODO())
		err = statusCmd.Err()
	})

	return redisClient, err
}

// NewRedis is the same function that NewRedisClient but does not return an error instead it panics
func NewRedis() *redis.Client {
	client, err := NewRedisClient()
	if err != nil {
		panic(err)
	}

	return client
}

// SQL variables
var (
	onceSQL sync.Once
	sqlDB   *sql.DB
)

// NewSQLDB returns a singleton of *sql.DB (uses postgresql driver)
//
// Required environment variables:
//
// DB_USER (postgresql user)
// DB_PASSWORD (postgresql password)
// DB_HOST (postgresql host)
// DB_PORT (postgresql port)
// DB_NAME (postgresql database name)
// DB_SSL (postgresql ssl)
func NewSQLDB() (*sql.DB, error) {
	var err error

	onceSQL.Do(func() {
		driverName := os.Getenv("DB_DRIVER")

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL"),
		)

		sqlDB, err = sql.Open(driverName, dsn)
		if err != nil {
			return
		}

		err = sqlDB.Ping()
	})

	return sqlDB, err
}

// NewSQL is the same function that NewSQLDB but it does not return an error, instead it panics
func NewSQL() *sql.DB {
	connection, err := NewSQLDB()
	if err != nil {
		panic(err)
	}

	return connection
}
