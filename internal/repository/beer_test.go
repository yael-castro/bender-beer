package repository

import (
	"errors"
	"github.com/yael-castro/bender-beer/internal/model"
	"reflect"
	"strconv"
	"testing"
)

var (
	defaultMemoryBeer = NewBS(Memory)

	// ttCreateBeer is the test table used to test the method CreateBeer of BeerStorage
	ttCreateBeer = []struct {
		BeerStorage
		model.Beer
		expectedError error
	}{
		// This case trys create an empty model.Beer
		{
			BeerStorage: defaultMemoryBeer,
		},
		// This case trys create a model.Beer with an exists id
		// In this case use the memoryBeerStorage as mock
		{
			BeerStorage: &memoryBeerStorage{
				data: map[int]model.Beer{
					1: model.Beer{},
				},
			},
			Beer: model.Beer{
				Id:      1,
				Name:    "myBeer",
				Brewery: "myBrewery",
			},
		},
	}

	// ttGetAllBeers is the test table used to test the method GetAllBeers of BeerStorage
	ttGetAllBeers = []struct {
		BeerStorage
		expectedError error
		expectedData  model.Beers
	}{
		// This case preload data in a memoryStorage to then get and
		//check if the expected data is the same that the preload data
		{
			BeerStorage: &memoryBeerStorage{
				data: map[int]model.Beer{
					1: {Id: 1, Name: "Beer 1"},
					2: {Id: 2, Name: "Beer 2"},
				},
			},
			expectedData: model.Beers{
				{Id: 1, Name: "Beer 1"},
				{Id: 2, Name: "Beer 2"},
			},
		},
	}

	// ttCreateBeer is the test table used to test the method GetBeerById of BeerStorage
	ttGetBeerById = []struct {
		BeerStorage
		expectedBeer  model.Beer
		expectedError error
	}{
		{
			BeerStorage: &memoryBeerStorage{
				data: map[int]model.Beer{
					1: {Id: 1, Name: "Beer 1"},
					2: {Id: 2, Name: "Beer 2"},
				},
			},
			expectedBeer: model.Beer{Id: 1, Name: "Beer 1"},
		},
		{
			// BeerStorage are empty
			BeerStorage: NewBS(Memory),
			// expectedBeer no exists
			expectedBeer:  model.Beer{Id: 1_000},
			expectedError: model.NotFound("not found beer with id '1000'"),
		},
	}
)

func TestBeerStorage_CreateBeer(t *testing.T) {
	for i, v := range ttCreateBeer {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			err := v.CreateBeer(&v.Beer)
			if !errors.Is(err, v.expectedError) {
				t.Fatalf(`expected error "%v" got "%v"`, v.expectedError, err)
			}

			gotBeer, err := v.GetBeerById(v.Id)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(v.Beer, gotBeer) {
				t.Logf(`expected beer "%+v" got "+%+v"`, v.Beer, gotBeer)
			}

			t.Logf("%+v", v.Beer)
		})
	}
}

func TestBeerStorage_GetAllBeers(t *testing.T) {
	for i, v := range ttGetAllBeers {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			beers, err := v.GetAllBeers()
			if !errors.Is(err, v.expectedError) {
				t.Fatalf(`expected error "%v" got "%v"`, v.expectedError, err)
			}

			if err != nil {
				t.Skip(err)
			}

			if !reflect.DeepEqual(beers, v.expectedData) {
				t.Logf(`expected beers "%+v" got "+%+v"`, v.expectedData, beers)
			}

			t.Logf("%+v", beers)
		})
	}
}

func TestBeerStorage_GetBeerById(t *testing.T) {
	for i, v := range ttGetBeerById {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			gotBeer, err := v.GetBeerById(v.expectedBeer.Id)
			if !errors.Is(err, v.expectedError) {
				t.Fatalf(`expected error "%v" got "%v"`, v.expectedError, err)
			}

			if err != nil {
				t.Skip(err)
			}

			if !reflect.DeepEqual(gotBeer, v.expectedBeer) {
				t.Logf(`expected beers "%+v" got "+%+v"`, v.expectedBeer, gotBeer)
			}

			t.Logf("%+v", gotBeer)
		})
	}
}
