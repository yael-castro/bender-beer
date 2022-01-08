// Package repository persistence layer.
// Here you will find everything related to data persistence like crud operations and repository conenctions
package repository

import "github.com/yael-castro/bender-beer/internal/model"

// Type defines the available storage types
type Type uint

// Repository types
const (
	// Memory will store the data in memory
	Memory Type = iota
	// Mock defines a mock storage (fake storage used to testing)
	Mock
	// ThirdPartyAPI defines a REST API as storage
	ThirdPartyAPI
	SQL
)

// beerData preload beer data used in storage of type Memory
var beerData = map[int]model.Beer{
	0: {
		Id:       0,
		Currency: "USD",
		Price:    2,
	},
	1: {
		Id:       1,
		Currency: "MXN",
		Price:    17.50,
	},
}

const (
	insertBeer = `INSERT INTO "beers"(name, brewery, country, currency, price) VALUES ($1, $2, $3, $4, $5) RETURNING beer_id`

	selectBeers = `SELECT beer_id, name, brewery, country, currency, price FROM "beers"`

	selectBeer = `SELECT name, brewery, country, currency, price FROM "beers" WHERE beer_id = $1`
)
