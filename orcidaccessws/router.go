package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/handlers"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routeSlice []route

var routes = routeSlice{

	route{
		"FaveIcon",
		"GET",
		"/favicon.ico",
		handlers.FavIconHandler,
	},

	route{
		"HealthCheck",
		"GET",
		"/healthcheck",
		handlers.HealthCheck,
	},

	route{
		"VersionInfo",
		"GET",
		"/version",
		handlers.VersionInfo,
	},

	route{
		"GetOrcidAttributes",
		"GET",
		"/cid/{id}",
		handlers.GetOrcidAttributes,
	},

	route{
		"GetAllOrcidAttributes",
		"GET",
		"/cid",
		handlers.GetAllOrcidAttributes,
	},

	route{
		"GetOrcidDetails",
		"GET",
		"/orcid/{id}",
		handlers.GetOrcidDetails,
	},

	route{
		"SearchOrcid",
		"GET",
		"/orcid",
		handlers.SearchOrcid,
	},

	route{
		"SetOrcidAttributes",
		"PUT",
		"/cid/{id}",
		handlers.SetOrcidAttributes,
	},

	route{
		"UpdateActivity",
		"PUT",
		"/cid/{id}/activity",
		handlers.UpdateActivity,
	},

	route{
		"DeleteOrcidAttributes",
		"DELETE",
		"/cid/{id}",
		handlers.DeleteOrcidAttributes,
	},
}

//
// NewRouter -- build and return the router
//
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// add the route for the prometheus metrics
	router.Handle("/metrics", HandlerLogger(promhttp.Handler(), "promhttp.Handler"))

	for _, route := range routes {

		var handler http.Handler = route.HandlerFunc
		handler = HandlerLogger(handler, route.Name)
		handler = prometheus.InstrumentHandler(route.Name, handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

//
// end of file
//
