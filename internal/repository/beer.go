package repository

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/yael-castro/bender-beer/internal/model"
)

// BeerStorage defines the storage interface for beer data
type BeerStorage interface {
	// CreateBeer use to create a beer
	CreateBeer(*model.Beer) error
	// GetAllBeers use to list all available beers
	GetAllBeers() (model.Beers, error)
	// GetBeerById use to get an beer by id
	GetBeerById(int) (model.Beer, error)
}

// NewBeerStorage abstract factory to BeerStorage based on a Type passed as parameter
//
// Returns an error if exists an error with BeerStorage implementation or a invalid Type is passed as parameter
//
// Supported types: Memory
func NewBeerStorage(t Type) (BeerStorage, error) {
	switch t {
	case Memory:
		return &memoryBeerStorage{
			data: beerData,
		}, nil

	case SQL:
		db, err := NewSQLDB()

		return &postgresqlBeerStorage{
			DB: db,
		}, err
	}

	return nil, fmt.Errorf(`type "%d" is not supported by repository.NewUserProvider`, t)
}

// NewBS works at the same that the NewCurrencyProvider but does not return an error instead panics
func NewBS(t Type) BeerStorage {
	storage, err := NewBeerStorage(t)
	if err != nil {
		panic(err)
	}

	return storage
}

// memoryBeerStorage BeerStorage that uses a hash map of model.Beer as storage
type memoryBeerStorage struct {
	data    map[int]model.Beer
	counter int
}

// CreateBeer creates a beer in the memory storage (hash map)
func (m *memoryBeerStorage) CreateBeer(beer *model.Beer) error {
	m.counter++

	beer.Id = m.counter
	_, ok := m.data[beer.Id]
	if ok {
		return fmt.Errorf(`beer id '%d' already exists`, beer.Id)
	}

	m.data[beer.Id] = *beer
	return nil
}

// GetBeerById search a beer by id in the m.data
//
// Search complexity: O(1)
func (m *memoryBeerStorage) GetBeerById(id int) (beer model.Beer, err error) {
	beer, ok := m.data[id]
	if !ok {
		err = model.NotFound(fmt.Sprintf(`not found beer with id '%d'`, id))
	}

	return
}

// GetAllBeers parse the m.data (hash map) to model.Beers
func (m memoryBeerStorage) GetAllBeers() (beers model.Beers, err error) {
	var keys []int

	for k := range m.data {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for k := range m.data {
		beers = append(beers, m.data[k])
	}

	return
}

// postgresqlBeerStorage
type postgresqlBeerStorage struct {
	*sql.DB
}

func (p postgresqlBeerStorage) CreateBeer(beer *model.Beer) (err error) {
	err = p.QueryRow(insertBeer,
		&beer.Name,
		&beer.Brewery,
		&beer.Country,
		&beer.Currency,
		&beer.Price,
	).Scan(&beer.Id)
	return
}

func (p postgresqlBeerStorage) GetAllBeers() (beers model.Beers, err error) {
	rows, err := p.Query(selectBeers)
	if err != nil {
		return
	}

	for rows.Next() {
		var beer model.Beer

		err = rows.Scan(
			&beer.Id, &beer.Name, &beer.Brewery, &beer.Country, &beer.Currency, &beer.Price,
		)
		if err != nil {
			return
		}

		beers = append(beers, beer)
	}

	return
}

func (p postgresqlBeerStorage) GetBeerById(beerId int) (beer model.Beer, err error) {
	err = p.QueryRow(selectBeer, beerId).Scan(
		&beer.Name, &beer.Brewery, &beer.Country, &beer.Currency, &beer.Price,
	)

	return
}
