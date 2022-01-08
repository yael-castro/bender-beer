// Package handler handles all http requests made to the app (is the presentation layer)
package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// defaultMaxAge default age to cache
const defaultMaxAge = 100

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

// New constructs an empty *Handler
func New() *Handler {
	return &Handler{}
}

// Handler main handler used in the ...
type Handler struct {
	*echo.Echo
	Beer
}

func (h Handler) Init() {
	h.Pre(middleware.AddTrailingSlash())
	h.Use(middleware.Recover(), middleware.Logger(), middleware.CORS())

	// Endpoint settings
	h.GET("/bender-beer/v1/health-check/", HealthCheck)

	h.POST("/bender-beer/v1/beers/", h.CreateBeer)

	h.GET("/bender-beer/v1/beers/:beerId/", h.GetBeerById)
	h.GET("/bender-beer/v1/beers/:beerId/boxprice/", h.GetBeerBox)
	h.GET("/bender-beer/v1/beers/", h.GetAllBeers)
}
