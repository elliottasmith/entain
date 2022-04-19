package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elliottasmith/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"syreclabs.com/go/faker"
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

func TestScanRacesStatusClosed(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Create a mock race to be returned from the scan races function
	mockRace := []*racing.Race{
		{
			Id: faker.RandomInt64(1, 10),
			MeetingId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Number: faker.RandomInt64(1, 10),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New((time.Now().AddDate(0, 0, -1))),
			Status: racing.Status_CLOSED,
		},
	}

	// To test scan races an sql row is generated using sqlmock (work around instead of creating complete mock database)
	mockRows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).AddRow(mockRace[0].GetId(), mockRace[0].GetMeetingId(), mockRace[0].GetName(), mockRace[0].GetNumber(), mockRace[0].GetVisible(), mockRace[0].GetAdvertisedStartTime().AsTime())
	sqlRows := mockRowsToSqlRows(mockRows)

	// Call scan races with mock sql rows
	races, err := rr.scanRaces(sqlRows)

	// Assert mock race matches race result from scan races
	assert.NoError(t, err)
	assert.Equal(t, mockRace, races)
}

func TestScanRacesStatusOpen(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Create a mock race to be returned from the scan races function
	mockRace := []*racing.Race{
		{
			Id: faker.RandomInt64(1, 10),
			MeetingId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Number: faker.RandomInt64(1, 10),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New((time.Now().AddDate(0, 0, 1))),
			Status: racing.Status_OPEN,
		},
	}

	// To test scan races an sql row is generated using sqlmock (work around instead of creating complete mock database)
	mockRows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).AddRow(mockRace[0].GetId(), mockRace[0].GetMeetingId(), mockRace[0].GetName(), mockRace[0].GetNumber(), mockRace[0].GetVisible(), mockRace[0].GetAdvertisedStartTime().AsTime())
	sqlRows := mockRowsToSqlRows(mockRows)

	// Call scan races with mock sql rows
	races, err := rr.scanRaces(sqlRows)

	// Assert mock race matches race result from scan races
	assert.NoError(t, err)
	assert.Equal(t, mockRace, races)
}

func TestScanRace(t *testing.T) {
	// Create RacesRepo without DB, not mocking SQL DB here.
	rr := &racesRepo{}

	// Create a mock race to be returned from the scan race function
	mockRace := []*racing.Race{
		{
			Id: faker.RandomInt64(1, 10),
			MeetingId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Number: faker.RandomInt64(1, 10),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New((time.Now().AddDate(0, 0, -1))),
			Status: racing.Status_CLOSED,
		},
	}

	// To test scan race an sql row is generated using sqlmock (work around instead of creating complete mock database)
	mockRows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).AddRow(mockRace[0].GetId(), mockRace[0].GetMeetingId(), mockRace[0].GetName(), mockRace[0].GetNumber(), mockRace[0].GetVisible(), mockRace[0].GetAdvertisedStartTime().AsTime())
	sqlRow := mockRowsToSqlRow(mockRows)

	// Call scan races with mock sql rows
	race, err := rr.scanRace(sqlRow)

	// Assert mock race matches race result from scan races
	assert.NoError(t, err)
	assert.Equal(t, mockRace[0], race)
}

// Helpers
func mockRowsToSqlRows(mockRows *sqlmock.Rows) *sql.Rows {
    db, mock, _ := sqlmock.New()
    mock.ExpectQuery("select").WillReturnRows(mockRows)
    rows, _ := db.Query("select")
	db.Close()
    return rows
}

func mockRowsToSqlRow(mockRows *sqlmock.Rows) *sql.Row {
    db, mock, _ := sqlmock.New()
    mock.ExpectQuery("select").WillReturnRows(mockRows)
    row := db.QueryRow("select")
	db.Close()
    return row
}