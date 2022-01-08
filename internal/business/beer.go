// Package business orchestrate business logic
package business

import (
	"github.com/yael-castro/bender-beer/internal/model"
	"github.com/yael-castro/bender-beer/internal/repository"
	"log"
)

type BeerStorage interface {
	repository.BeerStorage
	GetBeerBox(int, int, string) (model.BeerBox, error)
}

// _ implements constraint for BeerProvider
var _ BeerStorage = BeerProvider{}

// BeerProvider beer provider
type BeerProvider struct {
	repository.BeerStorage
	repository.CurrencyProvider
}

func (p BeerProvider) CreateBeer(beer *model.Beer) error {
	if beer.Name == "" {
		return model.ValidationError("the name of beer cannot be empty")
	}

	if beer.Currency == "" {
		return model.ValidationError("the currency of beer cannot be empty")
	}

	_, err := p.ProvideCurrency(beer.Currency, beer.Currency)
	if err != nil {
		return err
	}

	return p.BeerStorage.CreateBeer(beer)
}

func (p BeerProvider) GetBeerBox(beerId, beerNumber int, currency string) (box model.BeerBox, err error) {
	beer, err := p.BeerStorage.GetBeerById(beerId)
	if err != nil {
		return
	}

	log.Println(beer.Currency, currency)
	if beer.Currency == currency {
		box.PriceTotal = beer.Price * float64(beerNumber)
		return
	}

	value, err := p.ProvideCurrency(beer.Currency, currency)
	if err != nil {
		return
	}

	box.PriceTotal = beer.Price * float64(beerNumber) * value
	return
}
