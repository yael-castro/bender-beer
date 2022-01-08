package repository

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"testing"
)

var (
	ttCurrency = []struct {
		CurrencyProvider
		beerCurrency    string
		paymentCurrency string
		expectedError   error
	}{
		{
			CurrencyProvider: NewCP(ThirdPartyAPI),
			beerCurrency:     "🧠",
			paymentCurrency:  "MXN",
			expectedError:    errors.New("You have supplied an invalid Source Currency. [Example: source=EUR]"),
		},
		{
			CurrencyProvider: NewCP(ThirdPartyAPI),
			beerCurrency:     "USD",
			paymentCurrency:  "MXN",
		},
		{
			CurrencyProvider: NewCP(ThirdPartyAPI),
			beerCurrency:     "MXN",
			paymentCurrency:  "USD",
		},
	}
)

func TestCurrencyProvider_ProvideCurrency(t *testing.T) {
	for i, v := range ttCurrency {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			value, err := v.CurrencyProvider.ProvideCurrency(v.beerCurrency, v.paymentCurrency)
			if reflect.DeepEqual(v.expectedError, err) {
				goto success
			}
			if !errors.Is(v.expectedError, err) {
				t.Fatalf(`expected error "%T:%+v" got "%T:%+v"`, v.expectedError, v.expectedError, err, err)
			}
		success:

			log.Println(value)
		})
	}
}
