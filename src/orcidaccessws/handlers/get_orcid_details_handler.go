package handlers

import (
    "net/http"
    "github.com/gorilla/mux"
    "orcidaccessws/authtoken"
    "orcidaccessws/config"
    "orcidaccessws/orcid"
)

func GetOrcidDetails( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    id := vars[ "id" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    Statistics.RequestCount++
    Statistics.GetOrcidDetailsCount++

    // parameters OK ?
    if NotEmpty( id ) == false || NotEmpty( token ) == false {
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
    orcids, status := orcid.GetOrcidDetails( id )

    // we did got an error, return it
    if status != http.StatusOK {
        encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // we did not find the item, return 404
    //if orcid == nil {
    //    status := http.StatusNotFound
    //    encodeOrcidDetailsResponse( w, status, http.StatusText( status ), nil )
    //    return
    //}

    status = http.StatusOK
    encodeOrcidDetailsResponse( w, status, http.StatusText( status ), orcids )
}