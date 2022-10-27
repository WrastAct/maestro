package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/WrastAct/maestro/internal/validator"
)

type Game struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ValidateGame(v *validator.Validator, game *Game) {
	v.Check(game.Name != "", "name", "must be provided")
	v.Check(len(game.Name) <= 100, "name", "must not be more than 100 bytes long")
}

type GameModel struct {
	DB *sql.DB
}

func (m GameModel) Insert(game *Game) error {
	query := `
		INSERT INTO games (games_name)
		VALUES ($1)
		RETURNING games_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, game.Name).Scan(&game.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m GameModel) Get(id int64) (*Game, error) {
	query := `
		SELECT games_id, games_name
		FROM games
		WHERE games_id = $1`

	var game Game

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&game.ID,
		&game.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &game, nil
}

func (m GameModel) GetAll() ([]*Game, error) {
	query := `
		SELECT games_id, games_name
		FROM games`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	games := []*Game{}

	for rows.Next() {
		var game Game

		err := rows.Scan(
			&game.ID,
			&game.Name,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (m GameModel) Update(game *Game) error {
	query := `
		UPDATE games
		SET games_name = $1
		WHERE games_id = $2`

	args := []interface{}{
		game.Name,
		game.ID,
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

func (m GameModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM games
		WHERE games_id = $1`

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
