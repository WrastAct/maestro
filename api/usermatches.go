package main

import (
	"net/http"

	"github.com/WrastAct/maestro/internal/data"
)

func (app *application) createUserMatchHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID        int64   `json:"user_id"`
		MatchID       int64   `json:"match_id"`
		TournamentID  int64   `json:"tournament_id"`
		Result        string  `json:"result"`
		AverageStress float64 `json:"avg_stress"`
		Humidity      float64 `json:"humidity"`
		Temperature   float64 `json:"temperature"`
		Pressure      float64 `json:"pressure"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userMatch := &data.UserMatch{
		UserID:        input.UserID,
		MatchID:       input.MatchID,
		TournamentID:  input.TournamentID,
		Result:        input.Result,
		AverageStress: input.AverageStress,
		Humidity:      input.Humidity,
		Temperature:   input.Temperature,
		Pressure:      input.Pressure,
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.UserMatch.Insert(userMatch)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"player_match": userMatch}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listUserMatchHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	userMatch, err := app.models.UserMatch.GetMatchesByUser(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"player_match": userMatch}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
