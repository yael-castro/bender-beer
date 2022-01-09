package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yael-castro/bender-beer/internal/business"
	"github.com/yael-castro/bender-beer/internal/model"
	"log"
	"mime"
	"net/http"
	"strconv"
)

// Beer defines the http handlers group used to handle all http requests related to beer storage
type Beer interface {
	// CreateBeer handles http requests made to create a new beer in the storage records
	CreateBeer(echo.Context) error
	// GetBeerById handles http requests made to get a beer by id
	GetBeerById(echo.Context) error
	// GetAllBeers handles http requests made to get all beer records from a storage
	GetAllBeers(echo.Context) error
	// GetBeerBox handles http requests made to get a price of beer box
	GetBeerBox(echo.Context) error
}

// BeerGroup container for all http handlers related to beer operations
type BeerGroup struct {
	// BeerStorage required storage to manages the beer records
	business.BeerStorage
}

func (b BeerGroup) CreateBeer(c echo.Context) error {
	beer := &model.Beer{}

	contentType, _, _ := mime.ParseMediaType(c.Request().Header.Get("Content-Type"))
	if contentType != "application/json" {
		return c.JSON(http.StatusBadRequest, model.Error{Info: fmt.Sprintf(`mime type '%v' is not supported`, contentType)})
	}

	err := c.Bind(beer)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Info: err.Error()})
	}

	err = b.BeerStorage.CreateBeer(beer)
	if _, ok := err.(model.ValidationError); ok {
		return c.JSON(http.StatusBadRequest, model.Error{Info: err.Error()})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Info: err.Error()})
	}

	return c.JSON(http.StatusCreated, beer)
}

func (b BeerGroup) GetBeerById(c echo.Context) error {
	id := c.Param("beerId")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.Error{Info: "missing path param 'id'"})
	}

	beerId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Info: "missing a numeric value for the path param: " + id})
	}

	beer, err := b.BeerStorage.GetBeerById(beerId)
	if _, ok := err.(model.NotFound); ok || errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, model.Error{Info: err.Error()})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Info: err.Error()})
	}

	c.Response().Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d", defaultMaxAge))
	c.Response().Header().Set("Vary", "User-Agent")

	return c.JSON(http.StatusOK, beer)
}

func (b BeerGroup) GetAllBeers(c echo.Context) error {
	beers, err := b.BeerStorage.GetAllBeers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Info: err.Error()})
	}

	c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", defaultMaxAge))
	c.Response().Header().Set("Vary", "User-Agent")

	return c.JSON(http.StatusOK, beers)
}

func (b BeerGroup) GetBeerBox(c echo.Context) error {
	id := c.Param("beerId")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.Error{Info: "missing path param 'id'"})
	}

	beerId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Info: "missing a numeric value for the path param"})
	}

	log.Println(id)

	currency := c.QueryParam("currency")
	if currency == "" {
		currency = "USD"
	}

	quantity, err := strconv.Atoi(c.QueryParam("boxprice"))
	if err != nil {
		quantity = 6
	}

	beerBox, err := b.BeerStorage.GetBeerBox(beerId, quantity, currency)
	if _, ok := err.(model.NotFound); ok || errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, model.Error{Info: err.Error()})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Info: err.Error()})
	}

	return c.JSON(http.StatusOK, beerBox)
}
