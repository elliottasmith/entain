package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elliottasmith/entain/sports/proto/sports"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"syreclabs.com/go/faker"
)

func TestApplyFilterEmpty(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply filter with an empty events request filter.
	query := er.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{
	})

	// Assert query does not contain visible clause and args are empty.
	assert.Equal(t, "SELECT * FROM events", query)
}

func TestApplyFilterVisibleFalse(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply filter with a events request filter, filtering all results.
	query := er.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{
		VisibleOnly: false,
	})

	// Assert query does not contain visible clause and args are empty.
	assert.Equal(t, "SELECT * FROM events", query)
}

func TestApplyFilterVisibleTrue(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply filter with a events request filter, filtering visible only.
	query := er.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{
		VisibleOnly: true,
	})

	// Assert query contains visible clause and args are empty.
	assert.Equal(t, "SELECT * FROM events WHERE visible = 1", query)
}

func TestApplyOrderEmpty(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply order with an empty events request order.
	query := er.applyOrder("SELECT * FROM events", &sports.ListEventsRequestOrder{
	})

	// Assert query does not contain order clause.
	assert.Equal(t, "SELECT * FROM events", query)
}

func TestApplyOrderAscending(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply order with a events request order, order by advertised start time with a ascending direction.
	query := er.applyOrder("SELECT * FROM events", &sports.ListEventsRequestOrder{
		Field: "advertised_start_time",
		Direction: sports.Direction_ASC,
	})

	// Assert query contains order clause.
	assert.Equal(t, "SELECT * FROM events ORDER BY advertised_start_time ASC", query)
}

func TestApplyOrderDescending(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply order with a events request order, order by advertised start time with a descending direction.
	query := er.applyOrder("SELECT * FROM events", &sports.ListEventsRequestOrder{
		Field: "advertised_start_time",
		Direction: sports.Direction_DESC,
	})

	// Assert query contains order clause.
	assert.Equal(t, "SELECT * FROM events ORDER BY advertised_start_time DESC", query)
}

func TestApplyOrderInvalidField(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Call apply order with a events request order, order by invalid field with a descending direction.
	query := er.applyOrder("SELECT * FROM events", &sports.ListEventsRequestOrder{
		Field: "invalid_field",
		Direction: sports.Direction_DESC,
	})

	// Assert query does not contain order clause.
	assert.Equal(t, "SELECT * FROM events", query)
}

func TestScanRacesStatusClosed(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Create a mock event to be returned from the scan events function
	mockEvent := []*sports.Event{
		{
			Id: faker.RandomInt64(1, 10),
			SportType: faker.Team().Creature(),
			League: faker.Company().Name(),
			Country: faker.Team().State(),
			LocationId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Round: faker.RandomInt64(1, 24),
			Game: faker.RandomInt64(1, 8),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New((time.Now().AddDate(0, 0, -1))),
			Status: sports.Status_CLOSED,
		},
	}

	// To test scan events an sql row is generated using sqlmock (work around instead of creating complete mock database)
	mockRows := sqlmock.NewRows([]string{"id", "sport_type", "league", "country", "location_id", "name", "round", "game", "visible", "advertised_start_time"}).AddRow(mockEvent[0].GetId(), mockEvent[0].GetSportType(), mockEvent[0].GetLeague(), mockEvent[0].GetCountry(), mockEvent[0].GetLocationId(), mockEvent[0].GetName(), mockEvent[0].GetRound(), mockEvent[0].GetGame(), mockEvent[0].GetVisible(), mockEvent[0].GetAdvertisedStartTime().AsTime())
	sqlRows := mockRowsToSqlRows(mockRows)

	// Call scan events with mock sql rows
	events, err := er.scanEvents(sqlRows)

	// Assert mock event matches event result from scan events
	assert.NoError(t, err)
	assert.Equal(t, mockEvent, events)
}

func TestScanRacesStatusOpen(t *testing.T) {
	// Create EventsRepo without DB, not mocking SQL DB here.
	er := &eventsRepo{}

	// Create a mock event to be returned from the scan events function
	mockEvent := []*sports.Event{
		{
			Id: faker.RandomInt64(1, 10),
			SportType: faker.Team().Creature(),
			League: faker.Company().Name(),
			Country: faker.Team().State(),
			LocationId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Round: faker.RandomInt64(1, 24),
			Game: faker.RandomInt64(1, 8),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New((time.Now().AddDate(0, 0, 1))),
			Status: sports.Status_OPEN,
		},
	}

	// To test scan events an sql row is generated using sqlmock (work around instead of creating complete mock database)
	mockRows := sqlmock.NewRows([]string{"id", "sport_type", "league", "country", "location_id", "name", "round", "game", "visible", "advertised_start_time"}).AddRow(mockEvent[0].GetId(), mockEvent[0].GetSportType(), mockEvent[0].GetLeague(), mockEvent[0].GetCountry(), mockEvent[0].GetLocationId(), mockEvent[0].GetName(), mockEvent[0].GetRound(), mockEvent[0].GetGame(), mockEvent[0].GetVisible(), mockEvent[0].GetAdvertisedStartTime().AsTime())
	sqlRows := mockRowsToSqlRows(mockRows)

	// Call scan events with mock sql rows
	events, err := er.scanEvents(sqlRows)

	// Assert mock event matches event result from scan events
	assert.NoError(t, err)
	assert.Equal(t, mockEvent, events)
}

// Helpers
func mockRowsToSqlRows(mockRows *sqlmock.Rows) *sql.Rows {
    db, mock, _ := sqlmock.New()
    mock.ExpectQuery("select").WillReturnRows(mockRows)
    rows, _ := db.Query("select")
	db.Close()
    return rows
}