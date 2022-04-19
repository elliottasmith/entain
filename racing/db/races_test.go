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

func TestApplyOrderEmpty(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply order with an empty races request order.
	query := rr.applyOrder("SELECT * FROM races", &racing.ListRacesRequestOrder{
	})

	// Assert query does not contain order clause.
	assert.Equal(t, "SELECT * FROM races", query)
}

func TestApplyOrderAscending(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply order with a races request order, order by advertised start time with a ascending direction.
	query := rr.applyOrder("SELECT * FROM races", &racing.ListRacesRequestOrder{
		Field: "advertised_start_time",
		Direction: racing.Direction_ASC,
	})

	// Assert query contains order clause.
	assert.Equal(t, "SELECT * FROM races ORDER BY advertised_start_time ASC", query)
}

func TestApplyOrderDescending(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply order with a races request order, order by advertised start time with a descending direction.
	query := rr.applyOrder("SELECT * FROM races", &racing.ListRacesRequestOrder{
		Field: "advertised_start_time",
		Direction: racing.Direction_DESC,
	})

	// Assert query contains order clause.
	assert.Equal(t, "SELECT * FROM races ORDER BY advertised_start_time DESC", query)
}

func TestApplyOrderInvalidField(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Call apply order with a races request order, order by invalid field with a descending direction.
	query := rr.applyOrder("SELECT * FROM races", &racing.ListRacesRequestOrder{
		Field: "invalid_field",
		Direction: racing.Direction_DESC,
	})

	// Assert query does not contain order clause.
	assert.Equal(t, "SELECT * FROM races", query)
}