package repository

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"testing"
)

var (
	// ttCurrency test table for CurrencyProvider
	ttCurrency = []struct {
		CurrencyProvider
		beerCurrency    string
		paymentCurrency string
		expectedError   error
		expectedValue   float64
	}{
		// This case prove to get a currency in real time
		{
			CurrencyProvider: NewCP(ThirdPartyAPI),
			beerCurrency:     "🧠",
			paymentCurrency:  "MXN",
			expectedError:    errors.New("You have supplied an invalid Source Currency. [Example: source=EUR]"),
		},
		// This case prove a get currency MXN to USD from data saved in memory
		{
			CurrencyProvider: NewCP(Memory),
			beerCurrency:     "USD",
			paymentCurrency:  "MXN",
			expectedValue:    20.38,
		},
		// This case prove a get currency MXN to USD from data saved in memory
		{
			CurrencyProvider: NewCP(Memory),
			beerCurrency:     "MXN",
			paymentCurrency:  "USD",
			expectedValue:    0.049,
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

			if err != nil {
				t.Skip(err)
			}

		success:
			if v.expectedValue != value {
				t.Fatalf(`expected currency value "%.2f" got "%.2f"`, v.expectedValue, value)
			}

			log.Println(value)
		})
	}
}
