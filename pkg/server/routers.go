package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/matetirpak/chessbot-playground-server/internal/api"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Swagger auto documentation
	router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	// Sessions
	Route{
		"DeleteSessions",
		strings.ToUpper("Delete"),
		"/chessserver/v1/sessions",
		api.DeleteSessions,
	},

	Route{
		"GetSessions",
		strings.ToUpper("Get"),
		"/chessserver/v1/sessions",
		api.GetSessions,
	},

	Route{
		"PostSessions",
		strings.ToUpper("Post"),
		"/chessserver/v1/sessions",
		api.PostSessions,
	},

	Route{
		"PutSessions",
		strings.ToUpper("Put"),
		"/chessserver/v1/sessions",
		api.PutSessions,
	},

	// Game
	Route{
		"GetGame",
		strings.ToUpper("Get"),
		"/chessserver/v1/game",
		api.GetGame,
	},

	Route{
		"PutGame",
		strings.ToUpper("Put"),
		"/chessserver/v1/game",
		api.PutGame,
	},
}
