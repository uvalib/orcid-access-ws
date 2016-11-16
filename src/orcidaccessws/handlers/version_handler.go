package handlers

import (
    "net/http"
)

func VersionGet( w http.ResponseWriter, r *http.Request ) {

    // update the statistics
    Statistics.RequestCount++

    // do the response
    encodeVersionResponse( w, http.StatusOK, Version( ) )
}