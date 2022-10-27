package data

import (
	"context"
	"database/sql"
	"time"
)

type TeamUsers struct {
	UserID    int64  `json:"user_id"`
	TeamID    int64  `json:"team_id"`
	JoinDate  string `json:"join_date"`
	LeaveDate string `json:"leave_date"`
	Role      string `json:"role"`
}

type TeamUsersModel struct {
	DB *sql.DB
}

func (m TeamUsersModel) GetAllByTeam(teamID int64) ([]*TeamUsers, error) {
	query := `
		SELECT user_id, join_date, leave_date
		FROM teams_users
		WHERE teams_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	totalRecords := 0
	teamsUsers := []*TeamUsers{}

	for rows.Next() {
		var teamUser TeamUsers

		teamUser.TeamID = teamID

		err := rows.Scan(
			&totalRecords,
			&teamUser.UserID,
			&teamUser.JoinDate,
			&teamUser.LeaveDate,
		)
		if err != nil {
			return nil, err
		}

		teamsUsers = append(teamsUsers, &teamUser)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teamsUsers, nil
}

func (m TeamUsersModel) GetAllByUser(userID int64) ([]*TeamUsers, error) {
	query := `
		SELECT teams_id, join_date, leave_date
		FROM teams_users
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	totalRecords := 0
	teamsUsers := []*TeamUsers{}

	for rows.Next() {
		var teamUser TeamUsers

		teamUser.UserID = userID

		err := rows.Scan(
			&totalRecords,
			&teamUser.UserID,
			&teamUser.JoinDate,
			&teamUser.LeaveDate,
		)
		if err != nil {
			return nil, err
		}

		teamsUsers = append(teamsUsers, &teamUser)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teamsUsers, nil
}
