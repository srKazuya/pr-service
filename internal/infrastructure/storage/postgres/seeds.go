package postgres

import (
	"database/sql"


	"github.com/brianvoe/gofakeit/v7"
)

func SeedUsersWithTeams(db *sql.DB) error {
	gofakeit.Seed(12345)

	teams := []string{"Alpha", "Beta", "Gamma"}
	for _, t := range teams {
		db.Exec(`INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT DO NOTHING`, t)
	}

	for i := 0; i < 10; i++ {
		db.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, true)
		`, gofakeit.UUID(), gofakeit.Username(), teams[i%3])
	}

	for i := 0; i < 5; i++ {
		db.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, NULL, false)
		`, gofakeit.UUID(), gofakeit.Username())
	}

	for i := 0; i < 3; i++ {
		db.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, NULL, true)
		`, gofakeit.UUID(), gofakeit.Username())
	}


	return nil
}