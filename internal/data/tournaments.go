package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/WrastAct/maestro/internal/validator"
)

type Tournament struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	GameID    int64  `json:"game_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func ValidateTournament(v *validator.Validator, tournament *Tournament) {
	v.Check(tournament.Name != "", "name", "must be provided")
	v.Check(len(tournament.Name) <= 100, "name", "must not be more than 100 bytes long")
	ValidateDate(v, tournament.StartDate)
	ValidateDate(v, tournament.EndDate)
}

type TournamentModel struct {
	DB *sql.DB
}

func (m TournamentModel) Get(id int64) (*Tournament, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT tournaments_id, tournaments_name, games_id, start_date, end_date
		FROM tournaments
		WHERE tournaments_id = $1`

	var tournament Tournament

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&tournament.ID,
		&tournament.Name,
		&tournament.GameID,
		&tournament.StartDate,
		&tournament.EndDate,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tournament, nil
}

func (m TournamentModel) Insert(tournament *Tournament) error {
	query := `
		INSERT INTO tournaments (tournaments_name, games_id, start_date, end_date)
		VALUES ($1, $2, $3, $4)
		RETURNING tournaments_id`

	args := []interface{}{tournament.Name, tournament.GameID, tournament.StartDate, tournament.EndDate}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&tournament.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m TournamentModel) Update(tournament *Tournament) error {
	query := `
		UPDATE tournaments
		SET tournaments_name = $1, games_id = $2, start_date = $3, end_date = $4
		WHERE tournaments_id = $5`

	args := []interface{}{
		tournament.Name,
		tournament.GameID,
		tournament.StartDate,
		tournament.EndDate,
		tournament.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m TournamentModel) GetAll() ([]*Tournament, error) {
	query := `
		SELECT tournaments_id, tournaments_name, games_id, start_date, end_date
		FROM tournaments`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tournaments := []*Tournament{}

	for rows.Next() {
		var tournament Tournament

		err := rows.Scan(
			&tournament.ID,
			&tournament.Name,
			&tournament.GameID,
			&tournament.StartDate,
			&tournament.EndDate,
		)
		if err != nil {
			return nil, err
		}

		tournaments = append(tournaments, &tournament)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tournaments, nil
}

func (m TournamentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM tournaments
		WHERE tournaments_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
