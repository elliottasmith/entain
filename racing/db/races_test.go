package db

import (
	"testing"

	"github.com/elliottasmith/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
)

func TestApplyFilterEmpty(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply filter with an empty races request filter.
	query, args := rr.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
	})

	// Assert query does not contain visible clause and args are empty.
	assert.Equal(t, "SELECT * FROM races", query)
	assert.Nil(t, args)
}

func TestApplyFilterVisibleFalse(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply filter with a races request filter, filtering all results.
	query, args := rr.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		VisibleOnly: false,
	})

	// Assert query does not contain visible clause and args are empty.
	assert.Equal(t, "SELECT * FROM races", query)
	assert.Nil(t, args)
}

func TestApplyFilterVisibleTrue(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply filter with a races request filter, filtering visible only.
	query, args := rr.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		VisibleOnly: true,
	})

	// Assert query contains visible clause and args are empty.
	assert.Equal(t, "SELECT * FROM races WHERE visible = 1", query)
	assert.Nil(t, args)
}
