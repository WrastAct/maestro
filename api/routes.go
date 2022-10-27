package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//TODO: add handlers

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requireAuthenticatedUser(app.healthcheckHandler))

	router.HandlerFunc(http.MethodGet, "/v1/tournaments/", app.requireAuthenticatedUser(app.listTournamentHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tournaments/", app.requirePermission("admin", app.createTournamentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tournaments/:id", app.showTournamentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tournaments/:id", app.requirePermission("admin", app.updateTournamentHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tournaments/:id", app.requirePermission("admin", app.deleteTournamentHandler))

	router.HandlerFunc(http.MethodGet, "/v1/teams/", app.requireAuthenticatedUser(app.listTeamHandler))
	router.HandlerFunc(http.MethodPost, "/v1/teams/", app.requirePermission("admin", app.createTeamHandler))
	router.HandlerFunc(http.MethodGet, "/v1/teams/:id", app.showTeamHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/teams/:id", app.requirePermission("admin", app.updateTeamHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/teams/:id", app.requirePermission("admin", app.deleteTeamHandler))

	router.HandlerFunc(http.MethodGet, "/v1/games/", app.requireAuthenticatedUser(app.listGameHandler))
	router.HandlerFunc(http.MethodPost, "/v1/games/", app.requirePermission("admin", app.createGameHandler))
	router.HandlerFunc(http.MethodGet, "/v1/games/:id", app.showGameHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/games/:id", app.requirePermission("admin", app.updateGameHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/games/:id", app.requirePermission("admin", app.deleteGameHandler))

	router.HandlerFunc(http.MethodPost, "/v1/teams_players/", app.requirePermission("admin", app.createTeamUsersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/teams_players/", app.requireActivatedUser(app.listTeamUsersHandler))

	router.HandlerFunc(http.MethodPost, "/v1/player_matches", app.requirePermission("admin", app.createUserMatchHandler))
	router.HandlerFunc(http.MethodGet, "/v1/player_matches", app.requireAuthenticatedUser(app.listUserMatchHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
