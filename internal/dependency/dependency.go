// Package dependency manages
package dependency

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yael-castro/bender-beer/internal/business"
	"github.com/yael-castro/bender-beer/internal/handler"
	"github.com/yael-castro/bender-beer/internal/repository"
)

// Profile defines options of dependency injection
type Profile uint

// Supported profiles for dependency injection
const (
	// Default defines the production profile
	Default Profile = iota
	// Testing defines the testing profile used to make a unit and integration tests
	Testing
)

// Injector defines a dependency injector
type Injector interface {
	// Inject takes any data type and fill of required dependencies (dependency injection)
	Inject(interface{}) error
}

// InjectorFunc function that implements the Injector interface
type InjectorFunc func(interface{}) error

func (f InjectorFunc) Inject(i interface{}) error {
	return f(i)
}

// NewInjector is an abstract factory to Injector, it builds a instance of Injector interface based on the Profile based as parameter
//
// Supported profiles: Default and Testing
//
// If pass a parameter an invalid profile it panics
func NewInjector(p Profile) Injector {
	switch p {
	case Default:
		return InjectorFunc(handlerDefault)
	}

	panic(fmt.Sprintf(`invalid profile: "%d" is not supported by NewInjector`, p))
}

// handlerDefault InjectorFunc for *handler.Handler that uses a Default Profile
func handlerDefault(i interface{}) error {
	h, ok := i.(*handler.Handler)
	if !ok {
		return fmt.Errorf("")
	}

	h.Beer = handler.BeerGroup{
		BeerStorage: business.BeerProvider{
			BeerStorage:      repository.NewBS(repository.SQL),
			CurrencyProvider: repository.NewCP(repository.ThirdPartyAPI),
		},
	}

	// Inject instance of echo
	h.Echo = echo.New()

	return nil
}