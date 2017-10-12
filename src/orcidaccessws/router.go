package main

import (
   "github.com/gorilla/mux"
   "net/http"
   "orcidaccessws/handlers"
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
      "RuntimeInfo",
      "GET",
      "/runtime",
      handlers.RuntimeInfo,
   },

   route{
      "Stats",
      "GET",
      "/statistics",
      handlers.StatsGet,
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

//
// end of file
//