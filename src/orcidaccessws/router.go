package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"orcidaccessws/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{

	Route{
		"HealthCheck",
		"GET",
		"/healthcheck",
		handlers.HealthCheck,
	},

	Route{
		"VersionGet",
		"GET",
		"/version",
		handlers.VersionGet,
	},

	Route{
		"RuntimeInfo",
		"GET",
		"/runtime",
		handlers.RuntimeInfo,
	},

	Route{
		"Stats",
		"GET",
		"/statistics",
		handlers.StatsGet,
	},

	Route{
		"GetOneOrcid",
		"GET",
		"/cid/{id}",
		handlers.GetOneOrcid,
	},

	Route{
		"GetAllOrcid",
		"GET",
		"/cid",
		handlers.GetAllOrcid,
	},

	Route{
		"GetOrcidDetails",
		"GET",
		"/orcid/{id}",
		handlers.GetOrcidDetails,
	},

	Route{
		"SearchOrcid",
		"GET",
		"/orcid",
		handlers.SearchOrcid,
	},

	Route{
		"SetOrcid",
		"PUT",
		"/cid/{id}/{orcid}",
		handlers.SetOrcid,
	},

	Route{
		"DeleteOrcid",
		"DELETE",
		"/cid/{id}",
		handlers.DeleteOrcid,
	},
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler

		handler = route.HandlerFunc
		handler = HandlerLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
