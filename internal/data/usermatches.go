package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserMatch struct {
	UserID        int64   `json:"user_id"`
	MatchID       int64   `json:"match_id"`
	TournamentID  int64   `json:"tournament_id"`
	Result        string  `json:"result"`
	AverageStress float64 `json:"avg_stress"`
	Humidity      float64 `json:"humidity"`
	Temperature   float64 `json:"temperature"`
	Pressure      float64 `json:"pressure"`
}

type UserMatchModel struct {
	DB *sql.DB
}

func (m UserMatchModel) Insert(userMatch *UserMatch) error {
	query := `
		INSERT INTO users_matches (user_id, matches_id, tournaments_id, result, average_stress,
			humidity, temperature, pressure)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []interface{}{
		userMatch.UserID,
		userMatch.MatchID,
		userMatch.TournamentID,
		userMatch.Result,
		userMatch.AverageStress,
		userMatch.Humidity,
		userMatch.Temperature,
		userMatch.Pressure,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m UserMatchModel) GetMatchesByUser(userID int64) ([]*UserMatch, error) {
	query := `
		SELECT users_id, matches_id, tournaments_id, result, average_stress, 
			humidity, temperature, pressure
		FROM users_matches
		WHERE users_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matches := []*UserMatch{}

	for rows.Next() {
		var match UserMatch

		err := rows.Scan(
			&match.UserID,
			&match.MatchID,
			&match.TournamentID,
			&match.Result,
			&match.AverageStress,
			&match.Humidity,
			&match.Temperature,
			&match.Pressure,
		)
		if err != nil {
			return nil, err
		}

		matches = append(matches, &match)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (m UserMatchModel) Update(userMatch *UserMatch) error {
	query := `
		UPDATE users_matches
		SET result = $1, average_stress = $2, humidity = $3, temperature = $4, pressure = $5
		WHERE users_id = $6 
		 AND matches_id = $7 
		 AND tournaments_id = $8`

	args := []interface{}{
		userMatch.Result,
		userMatch.AverageStress,
		userMatch.Humidity,
		userMatch.Temperature,
		userMatch.Pressure,
		userMatch.UserID,
		userMatch.MatchID,
		userMatch.TournamentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
