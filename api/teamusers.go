package main

import (
	"net/http"

	"github.com/WrastAct/maestro/internal/data"
)

func (app *application) createTeamUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID    int64  `json:"user_id"`
		TeamID    int64  `json:"team_id"`
		JoinDate  string `json:"join_date"`
		LeaveDate string `json:"leave_date"`
		Role      string `json:"role"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	teamUsers := &data.TeamUsers{
		UserID:    input.UserID,
		TeamID:    input.TeamID,
		JoinDate:  input.JoinDate,
		LeaveDate: input.LeaveDate,
		Role:      input.Role,
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.TeamUsers.Insert(teamUsers)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"team_players": teamUsers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTeamUsersHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	teamUsers, err := app.models.TeamUsers.GetAllByTeam(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"team_players": teamUsers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
