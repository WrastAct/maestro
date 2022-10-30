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

func (m TeamUsersModel) Insert(teamUsers *TeamUsers) error {
	query := `
		INSERT INTO teams_users (user_id, teams_id, join_date, leave_date, role)
		VALUES ($1, $2, $3, $4, $5)`

	args := []interface{}{teamUsers.UserID, teamUsers.TeamID, teamUsers.JoinDate,
		teamUsers.LeaveDate, teamUsers.Role}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m TeamUsersModel) Update(teamUsers *TeamUsers) error {
	query := `
		UPDATE teams_users
		SET join_date = $1, leave_date = $2, role = $3
		WHERE teams_id = $4 AND user_id = $5`

	args := []interface{}{
		teamUsers.JoinDate,
		teamUsers.LeaveDate,
		teamUsers.Role,
		teamUsers.TeamID,
		teamUsers.UserID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m TeamUsersModel) GetAllByTeam(teamID int64) ([]*TeamUsers, error) {
	query := `
		SELECT user_id, join_date, leave_date, role
		FROM teams_users
		WHERE teams_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teamsUsers := []*TeamUsers{}

	for rows.Next() {
		var teamUser TeamUsers

		teamUser.TeamID = teamID

		err := rows.Scan(
			&teamUser.UserID,
			&teamUser.JoinDate,
			&teamUser.LeaveDate,
			&teamUser.Role,
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

	teamsUsers := []*TeamUsers{}

	for rows.Next() {
		var teamUser TeamUsers

		teamUser.UserID = userID

		err := rows.Scan(
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
