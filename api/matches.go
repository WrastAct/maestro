package main

import (
	"errors"
	"net/http"

	"github.com/WrastAct/maestro/internal/data"
	"github.com/WrastAct/maestro/internal/validator"
)

func (app *application) createMatchHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TournamentID int64  `json:"tournament_id"`
		MatchData    string `json:"data"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	match := &data.Match{
		TournamentID: input.TournamentID,
		Matchdata:    input.MatchData,
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()

	if data.ValidateMatch(v, match); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Match.Insert(match)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"match": match}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMatchHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Match.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "match successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMatchHandler(w http.ResponseWriter, r *http.Request) {
	match, err := app.models.Match.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"match": match}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
