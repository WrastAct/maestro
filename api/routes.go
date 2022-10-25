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

	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requirePermission("admin", app.healthcheckHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/rooms", app.requireAuthenticatedUser(app.listRoomHandler))
	// router.HandlerFunc(http.MethodPost, "/v1/rooms", app.requirePermission("user", app.createRoomHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/rooms/:id", app.requirePermission("user", app.showRoomHandler))
	// router.HandlerFunc(http.MethodPatch, "/v1/rooms/:id", app.requirePermission("user", app.updateRoomHandler))
	// router.HandlerFunc(http.MethodDelete, "/v1/rooms/:id", app.requirePermission("user", app.deleteRoomHandler))

	// router.HandlerFunc(http.MethodGet, "/v1/furniture", app.listFurnitureHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/furniture", app.createFurnitureHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/furniture/:id", app.showFurnitureHandler)
	// router.HandlerFunc(http.MethodPatch, "/v1/furniture/:id", app.updateFurnitureHandler)
	// router.HandlerFunc(http.MethodDelete, "/v1/furniture/:id", app.deleteFurnitureHandler)

	// router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	//todo add authentication e.g. app.authenticate
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
