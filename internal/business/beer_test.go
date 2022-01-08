package business

import (
	"errors"
	"github.com/yael-castro/bender-beer/internal/repository"
	"log"
	"reflect"
	"strconv"
	"testing"
)

var (
	// ttGetBeerBox test table for BeerStorage
	ttGetBeerBox = []struct {
		BeerStorage
		beerId          int
		paymentCurrency string
		beerQuantity    int
		expectedError   error
		expectedValue   float64
	}{
		{
			BeerStorage: BeerProvider{
				BeerStorage:      repository.NewBS(repository.Memory),
				CurrencyProvider: repository.NewCP(repository.Memory),
			},
			beerId:          0,
			paymentCurrency: "MXN",
			expectedValue:   20.38 * 6 * 2,
			beerQuantity:    6,
		},
		// This case prove the payment for a beer box with nine beers (USD to MXN)
		{
			BeerStorage: BeerProvider{
				BeerStorage:      repository.NewBS(repository.Memory),
				CurrencyProvider: repository.NewCP(repository.Memory),
			},
			beerId:          0,
			paymentCurrency: "MXN",
			expectedValue:   20.38 * 9.0 * 2.0,
			beerQuantity:    9,
		},
		// This case prove the payment for a beer box with nine beers (MXN to USD)
		{
			BeerStorage: BeerProvider{
				BeerStorage:      repository.NewBS(repository.Memory),
				CurrencyProvider: repository.NewCP(repository.Memory),
			},
			beerId:          1,
			paymentCurrency: "USD",
			expectedValue:   0.049 * float64(6.0) * 17.5,
			beerQuantity:    6,
		},
	}
)

func TestBeerProvider_GetBeerBox(t *testing.T) {
	for i, v := range ttGetBeerBox {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			box, err := v.GetBeerBox(v.beerId, v.beerQuantity, v.paymentCurrency)
			if reflect.DeepEqual(v.expectedError, err) {
				goto success
			}
			if !errors.Is(v.expectedError, err) {
				t.Fatalf(`expected error "%T:%+v" got "%T:%+v"`, v.expectedError, v.expectedError, err, err)
			}

			if err != nil {
				t.Skip(err)
			}

		success:
			if v.expectedValue != box.PriceTotal {
				t.Fatalf(`expected currency value "%.2f" got "%.2f"`, v.expectedValue, box.PriceTotal)
			}

			log.Println(box.PriceTotal)
		})

	}
}
