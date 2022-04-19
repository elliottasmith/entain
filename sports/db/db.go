package db

import (
	"time"

	"syreclabs.com/go/faker"
)

func (r *eventsRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, sport_type TEXT, league TEXT, country TEXT, location_id INTEGER, name TEXT, round INTEGER, game INTEGER, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO events(id, sport_type, league, country, location_id, name, round, game, visible, advertised_start_time) VALUES (?,?,?,?,?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Creature(),
				faker.Company().Name(),
				faker.Team().State(),
				faker.Number().Between(1, 20),
				faker.Team().Name() + " vs " + faker.Team().Name(),
				faker.Number().Between(1, 24),
				faker.Number().Between(1, 8),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}
