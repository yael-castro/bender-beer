package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/bender-beer/internal/model"
	"log"
	"net/http"
	"os"
)

// CurrencyProvider defines a currency storage
type CurrencyProvider interface {
	// ProvideCurrency use to get the currency value
	ProvideCurrency(string, string) (float64, error)
}

// NewCurrencyProvider abstract factory for CurrencyProvider
// It creates a CurrencyProvider based on the Type passed as parameter
//
// Supported types: ThirdPartyAPI
func NewCurrencyProvider(t Type) (CurrencyProvider, error) {
	switch t {
	case ThirdPartyAPI:
		return &currencyRealTime{
			URI: os.Getenv("CURRENCY_API_URL"),
		}, nil
	}

	return nil, fmt.Errorf(`type "%d" is not supported by NewCurrencyProvider`, t)
}

// NewCP works at the same that the NewCurrencyProvider but does not return an error instead panics
func NewCP(t Type) CurrencyProvider {
	provider, err := NewCurrencyProvider(t)
	if err != nil {
		panic(err)
	}

	return provider
}

// currencyRealTime access to third party repository for consult the currency data in real time
type currencyRealTime struct {
	URI string
}

// ProvideCurrency sends a GET request from third party api for get the currency data in real time
//
// If the currency data is obtained successfully but the paymentCurrency does not exist
// the type of error returned is model.NotFound
func (c currencyRealTime) ProvideCurrency(beerCurrency, paymentCurrency string) (value float64, err error) {
	request, err := http.NewRequest(http.MethodGet, c.URI, nil)
	if err != nil {
		return
	}

	params := request.URL.Query()

	params.Set("source", beerCurrency)
	params.Set("currencies", paymentCurrency)

	request.URL.RawQuery = params.Encode()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	if code := response.StatusCode; code != http.StatusOK {
		err = fmt.Errorf(`currency api: unexpected code '%d'`, code)
		return
	}

	info := model.CurrencyInfo{}

	err = json.NewDecoder(response.Body).Decode(&info)
	if err != nil {
		return
	}

	if !info.Success {
		err = errors.New(info.Info) // no mismatch error
		return
	}

	log.Printf("%d\n", response.StatusCode)
	log.Printf("%+v\n", info)

	filter := beerCurrency + paymentCurrency

	value, ok := info.Quotes[filter]
	if !ok {
		err = model.NotFound(fmt.Sprintf("not found a currency with the key '%s'", filter))
	}

	return
}

// CacheCurrencyProvider wraps a CurrencyProvider to caches the currency provider to access more fast to the currency data
type CacheCurrencyProvider struct {
	CurrencyProvider
	*redis.Client
}

// ProvideCurrency uses the embed CurrencyProvider to get the currency data and also uses redis to caches the got data
// Note: The expiration time for the cache is equal to defaultRedisExpirationTime
func (p CacheCurrencyProvider) ProvideCurrency(beerCurrency, paymentCurrency string) (value float64, err error) {
	key := "currency:" + beerCurrency + paymentCurrency

	value, err = p.Get(context.TODO(), key).Float64()
	if err == nil {
		goto end
	}

	value, err = p.CurrencyProvider.ProvideCurrency(beerCurrency, paymentCurrency)
	if err != nil {
		return
	}

end:
	p.Set(context.TODO(), key, value, defaultRedisExpirationTime)
	return
}
