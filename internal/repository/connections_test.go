package repository

import (
	"testing"
)

// TestNewRedisClient health check for *redis.Client
func TestNewRedisClient(t *testing.T) {
	// os.Setenv("REDIS_HOST", "localhost")

	_, err := NewRedisClient()
	if err != nil {
		t.Fatal("failed redis connection", err)
	}
}
