package main

import (
   "net/http"
   "github.com/gorilla/mux"
   "orcidaccessws/handlers"
)

type Route struct {
   Name        string
   Method      string
   Pattern     string
   HandlerFunc http.HandlerFunc
}

type Routes [] Route

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
        "GetOneOrcid",
        "GET",
        "/cid",
        handlers.GetOneOrcid,
    },

        /*
    Route{
        "IdCreate",
        "POST",
        "/{shoulder:.*}",
        IdCreate,
    },

    Route{
        "IdRevoke",
        "PUT",
        "/revoke/{doi:.*}",
        IdRevoke,
    },

    Route{
        "IdUpdate",
        "PUT",
        "/{doi:.*}",
        IdUpdate,
    },

    Route{
        "IdDelete",
        "DELETE",
        "/{doi:.*}",
        IdDelete,
    },
    */
}

func NewRouter( ) *mux.Router {

   router := mux.NewRouter().StrictSlash( true )
   for _, route := range routes {

      var handler http.Handler

      handler = route.HandlerFunc
      handler = HandlerLogger( handler, route.Name )

      router.
         Methods( route.Method ).
         Path( route.Pattern ).
         Name( route.Name ).
         Handler( handler )
   }

   return router
}
