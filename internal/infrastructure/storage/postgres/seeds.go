package postgres

import (
	"database/sql"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedTeams(db *sql.DB, count int) error {
	gofakeit.Seed(0)

	for i := 0; i < count; i++ {
		_, err := db.Exec(`
            INSERT INTO teams (team_name)
            VALUES ($1)
        `, gofakeit.AppName())

		if err != nil {
			return err
		}
	}

	return nil
}

func SeedUsers(db *sql.DB, count int) error {
	gofakeit.Seed(0)

	for i := 0; i < count; i++ {
		_, err := db.Exec(`
            INSERT INTO users (username, user_id)
            VALUES ($1, $2)
        `, gofakeit.Username(), gofakeit.ID())

		if err != nil {
			return err
		}
	}

	return nil
}

func GetTeams(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT team_name FROM teams`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		teams = append(teams, name)
	}

	return teams, rows.Err()
}

func GetUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT user_id FROM users ORDER BY user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		users = append(users, id)
	}

	return users, rows.Err()
}


func AssignUsersToTeams(db *sql.DB) error {
	teams, err := GetTeams(db)
	if err != nil {
		return err
	}

	users, err := GetUsers(db)
	if err != nil {
		return err
	}

	if len(teams) == 0 {
		return nil
	}

	const perTeam = 10
	teamCount := len(teams)

	for i, userID := range users {
		teamIndex := (i / perTeam) % teamCount

		_, err := db.Exec(`
			UPDATE users
			SET team_name = $1
			WHERE user_id = $2
		`, teams[teamIndex], userID)

		if err != nil {
			return err
		}
	}

	return nil
}
