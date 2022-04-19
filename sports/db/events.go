package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/elliottasmith/entain/sports/proto/sports"
)

// EventsRepo provides repository access to events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter, order *sports.ListEventsRequestOrder) ([]*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the event repository dummy data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = r.seed()
	})

	return err
}

func (r *eventsRepo) List(filter *sports.ListEventsRequestFilter, order *sports.ListEventsRequestOrder) ([]*sports.Event, error) {
	var (
		err   error
		query string
	)

	query = getEventQueries()[eventsList]

	query = r.applyFilter(query, filter)

	query = r.applyOrder(query, order)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (r *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string) {
	var clauses []string

	if filter == nil {
		return query
	}

	// Check the VisibleOnly request value and only return visible results if true.
	if filter.GetVisibleOnly() {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query
}

func (r *eventsRepo) applyOrder(query string, order *sports.ListEventsRequestOrder) string {
	// Check that the input field validates and return the order by clause with the input field and direction
	if validateField(order.GetField()) {
		return fmt.Sprintf("%s ORDER BY %s %s", query, order.GetField(), order.GetDirection())
	}
	return query
}

func (m *eventsRepo) scanEvents(rows *sql.Rows) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.SportType, &event.League, &event.Country, &event.LocationId, &event.Name, &event.Round, &event.Game, &event.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		if err := setStartTime(&event, advertisedStart); err != nil {
			return nil, err
		}

		setStatus(&event, advertisedStart)

		events = append(events, &event)
	}

	return events, nil
}

func setStartTime(race *sports.Event, startTime time.Time) (error) {
	// Sets event start time to proto timestamp format
	ts, err := ptypes.TimestampProto(startTime)
	if err != nil {
		return err
	}

	race.AdvertisedStartTime = ts

	return nil
}

func setStatus(race *sports.Event, startTime time.Time) {
	// Set event status to open if the start time is in the future
	if (time.Now().Before(startTime)) {
		race.Status = sports.Status_OPEN
	}
}

func validateField(inputField string) bool {
	// Slice of valid fields
	validFields := []string{"id", "sport_type", "league", "country", "location_id", "name", "round", "game", "visible", "advertised_start_time" }

	// Validate that the input field is a valid field to prevent sql injection
	for _, validField := range validFields {
		if validField == inputField {
			return true
		}
	}
	return false
}
