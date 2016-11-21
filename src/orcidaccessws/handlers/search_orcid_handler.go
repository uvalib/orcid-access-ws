package handlers

import (
    "net/http"
  //  "github.com/gorilla/mux"
    "orcidaccessws/authtoken"
    "orcidaccessws/config"
    "orcidaccessws/orcid"
    "orcidaccessws/logger"
    "fmt"
)

const DEFAULT_SEARCH_START_IX = "0"
const DEFAULT_SEARCH_MAX_RESULTS = "50"

func SearchOrcid( w http.ResponseWriter, r *http.Request ) {

    //vars := mux.Vars( r )
    query := r.URL.Query( ).Get( "q" )
    token := r.URL.Query( ).Get( "auth" )
    start := r.URL.Query( ).Get( "start" )
    count := r.URL.Query( ).Get( "max" )

    // update the statistics
    Statistics.RequestCount++
    Statistics.SearchOrcidDetailsCount++

    // parameters OK ?
    if nonEmpty( query ) == false || nonEmpty( token ) == false {
        status := http.StatusBadRequest
        encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // check the supplied parameters and set defaults as necessary
    if nonEmpty( start ) == false {
        start = DEFAULT_SEARCH_START_IX
    }
    if nonEmpty( count ) == false {
        count = DEFAULT_SEARCH_MAX_RESULTS
    }

    // validate parameters as necessary
    if isNumeric( start ) == false || isNumeric( count ) == false {
        status := http.StatusBadRequest
        encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "getorcid", token ) == false {
        status := http.StatusForbidden
        encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // get the ORCID details
    orcids, status, err := orcid.SearchOrcid( query, start, count )

    // we did got an error, return it
    if status != http.StatusOK {
        encodeOrcidDetailsResponse( w, status,
            fmt.Sprintf( "%s (%s)", http.StatusText( status ), err ), nil )
        return
    }

    logger.Log( fmt.Sprintf( "ORCID search: %d result(s) located", len( orcids ) ) )

    status = http.StatusOK
    encodeOrcidDetailsResponse( w, status, http.StatusText( status ), orcids )
}