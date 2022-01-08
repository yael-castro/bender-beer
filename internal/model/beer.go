package model

type (
	// Beers alias for []Beer
	Beers = []Beer

	Beer struct {
		Id       int
		Name     string
		Brewery  string
		Country  string
		Price    float64
		Currency string
	}

	BeerBox struct {
		PriceTotal float64 `json:"Price Total"`
	}

	CurrencyInfo struct {
		Success   bool               `json:"success"`
		Timestamp uint               `json:"timestamp"`
		Source    string             `json:"source"`
		Quotes    map[string]float64 `json:"quotes"`
		Error     `json:"error,omitempty"`
	}
)
