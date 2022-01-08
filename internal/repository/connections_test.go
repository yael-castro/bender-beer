package repository

import (
	"testing"
)

// TestNewRedisClient health check for *redis.Client
func TestNewRedisClient(t *testing.T) {
	// os.Setenv("REDIS_HOST", "localhost")
	client, err := NewRedisClient()
	if err != nil {
		t.Fatal("failed redis connection", err)
	}

	client.Close()
}

// TestNewSQLDB health check for *sql.DB
func TestNewSQLDB(t *testing.T) {
	db, err := NewSQLDB()
	if err != nil {
		t.Fatal("failed sql connection", err)
	}

	db.Close()
}
