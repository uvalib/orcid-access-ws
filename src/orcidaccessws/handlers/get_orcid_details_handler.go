package handlers

import (
    "net/http"
    "github.com/gorilla/mux"
    "orcidaccessws/authtoken"
    "orcidaccessws/config"
    "orcidaccessws/orcid"
    "fmt"
)

func GetOrcidDetails( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    id := vars[ "id" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    Statistics.RequestCount++
    Statistics.GetOrcidDetailsCount++

    // parameters OK ?
    if nonEmpty( id ) == false || nonEmpty( token ) == false {
        status := http.StatusBadRequest
        encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.Timeout ) == false {
        status := http.StatusForbidden
        encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // get the ORCID details
    orcid, status, err := orcid.GetOrcidDetails( id )

    // we did got an error, return it
    if status != http.StatusOK {
        encodeOrcidDetailsResponse( w, status,
            fmt.Sprintf( "%s (%s)", http.StatusText( status ), err ), nil )
        return
    }

    status = http.StatusOK
    encodeOrcidDetailsResponse( w, status, http.StatusText( status ), orcid )
}