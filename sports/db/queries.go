package db

const (
	eventsList = "list"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT 
				id,
				sport_type,
				league,
				country,
				location_id,
				name,
				round,
				game,
				visible, 
				advertised_start_time 
			FROM events
		`,
	}
}
