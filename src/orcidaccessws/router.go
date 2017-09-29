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
		"GetOrcidAttributes",
		"GET",
		"/cid/{id}",
		handlers.GetOrcidAttributes,
	},

	Route{
		"GetAllOrcidAttributes",
		"GET",
		"/cid",
		handlers.GetAllOrcidAttributes,
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
		"SetOrcidAttributes",
		"PUT",
		"/cid/{id}",
		handlers.SetOrcidAttributes,
	},

	Route{
		"UpdateActivity",
		"PUT",
		"/cid/{id}/activity",
		handlers.UpdateActivity,
	},

	Route{
		"DeleteOrcidAttributes",
		"DELETE",
		"/cid/{id}",
		handlers.DeleteOrcidAttributes,
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
