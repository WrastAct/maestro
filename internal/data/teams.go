package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/WrastAct/maestro/internal/validator"
)

type Team struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Region      string `json:"region"`
}

func ValidateTeam(v *validator.Validator, team *Team) {
	v.Check(team.Name != "", "name", "must be provided")
	v.Check(len(team.Name) <= 100, "name", "must not be more than 100 bytes long")
	v.Check(len(team.Region) <= 32, "region", "must not be more than 32 bytes long")
}

type TeamModel struct {
	DB *sql.DB
}

func (m TeamModel) Get(id int64) (*Team, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT teams_id, teams_name, teams_description, teams_region
		FROM teams
		WHERE teams_id = $1`

	var team Team

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.Region,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &team, nil
}

func (m TeamModel) Insert(team *Team) error {
	query := `
		INSERT INTO teams (teams_name, teams_description, teams_region)
		VALUES ($1, $2, $3)
		RETURNING teams_id`

	args := []interface{}{team.Name, team.Description, team.Region}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&team.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m TeamModel) Update(team *Team) error {
	query := `
		UPDATE teams
		SET teams_name = $1, teams_description = $2, teams_region = $3
		WHERE teams_id = $4`

	args := []interface{}{
		team.Name,
		team.Description,
		team.Region,
		team.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m TeamModel) GetAllByRegion(region string) ([]*Team, error) {
	query := `
		SELECT teams_name, teams_description
		FROM teams
		WHERE teams_region = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, region)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	totalRecords := 0
	teams := []*Team{}

	for rows.Next() {
		var team Team

		err := rows.Scan(
			&totalRecords,
			&team.Name,
			&team.Description,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, &team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
