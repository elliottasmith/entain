package db

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
	"time"
	"fmt"

	"github.com/elliottasmith/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, order *racing.ListRacesRequestOrder) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, order *racing.ListRacesRequestOrder) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.applyOrder(query, order)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// Check the VisibleOnly request value and only return visible results if true.
	if filter.GetVisibleOnly() {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applyOrder(query string, order *racing.ListRacesRequestOrder) string {
	// Check that the input field validates and return the order by clause with the input field and direction
	if validateField(order.GetField()) {
		return fmt.Sprintf("%s ORDER BY %s %s", query, order.GetField(), order.GetDirection())
	}
	return query
}

func (m *racesRepo) scanRaces(rows *sql.Rows,) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		setStatus(&race, advertisedStart)

		races = append(races, &race)
	}

	return races, nil
}

func setStatus(race *racing.Race, startTime time.Time) {
	// Set race status to open if the start time is in the future
	if (time.Now().Before(startTime)) {
		race.Status = racing.Status_OPEN
	}
}

func validateField(inputField string) bool {
	// Slice of valid fields
	validFields := []string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}

	// Validate that the input field is a valid field to prevent sql injection
	for _, validField := range validFields {
		if validField == inputField {
			return true
		}
	}
	return false
}