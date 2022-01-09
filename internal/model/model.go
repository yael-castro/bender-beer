package model

import "github.com/lib/pq"

type (
	// Map used to make a unique structure as unique http responses
	Map map[string]interface{}

	// Error response for bad http requests
	Error struct {
		// Info description about error
		Info string `json:"info,omitempty"`
	}

	// SQLError alias for pq.Error
	SQLError = pq.Error
)
