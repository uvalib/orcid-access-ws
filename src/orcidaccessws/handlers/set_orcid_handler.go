package handlers

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "orcidaccessws/authtoken"
    "orcidaccessws/config"
    "orcidaccessws/dao"
    "orcidaccessws/logger"
)

func SetOrcid( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    id := vars[ "id" ]
    orcid := vars[ "orcid" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    Statistics.RequestCount++
    Statistics.SetOrcidCount++

    // parameters OK ?
    if nonEmpty( id ) == false || nonEmpty( orcid ) == false || nonEmpty( token ) == false {
        status := http.StatusBadRequest
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "setorcid", token ) == false {
        status := http.StatusForbidden
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // get the ORCID details
    err := dao.Database.SetOrcidByCid( id, orcid )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: %s", err.Error( ) ) )
        status := http.StatusInternalServerError
        encodeOrcidResponse( w, status,
            fmt.Sprintf( "%s (%s)", http.StatusText( status ), err ),
            nil )
        return
    }

    status := http.StatusOK
    encodeStandardResponse( w, status, http.StatusText( status ) )
}