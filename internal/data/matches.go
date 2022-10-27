package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Match struct {
	ID           int64  `json:"id"`
	TournamentID int64  `json:"tournament_id"`
	Matchdata    string `json:"data"`
}

type MatchModel struct {
	DB *sql.DB
}

func (m MatchModel) Insert(match *Match) error {
	query := `
		INSERT INTO matches (matches_id, tournaments_id, match_data)
		VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, match.ID, match.TournamentID, match.Matchdata)
	if err != nil {
		return err
	}
	return nil
}

func (m MatchModel) GetByTournamentID(tournamentID int64) ([]*Match, error) {
	query := `
		SELECT tournaments_id, matches_id, match_data
		FROM matches
		WHERE tournaments_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, tournamentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matches := []*Match{}

	for rows.Next() {
		var match Match

		err := rows.Scan(
			&match.TournamentID,
			&match.ID,
			&match.Matchdata,
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

func (m MatchModel) Update(match *Match) error {
	query := `
		UPDATE matches
		SET match_data = $1
		WHERE matches_id = $2 AND tournaments_id = $3`

	args := []interface{}{
		match.Matchdata,
		match.ID,
		match.TournamentID,
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
