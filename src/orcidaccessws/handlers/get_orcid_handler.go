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

func GetOneOrcid( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    id := vars[ "id" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    Statistics.RequestCount++
    Statistics.GetOrcidCount++

    // parameters OK ?
    if nonEmpty( id ) == false || nonEmpty( token ) == false {
        status := http.StatusBadRequest
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.Timeout ) == false {
        status := http.StatusForbidden
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // get the ORCID details
    orcids, err := dao.Database.GetOrcidByCid( id )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: %s", err.Error( ) ) )
        status := http.StatusInternalServerError
        encodeOrcidResponse( w, status,
            fmt.Sprintf( "%s (%s)", http.StatusText( status ), err ),
            nil )
        return
    }

    // we did not find the item, return 404
    if orcids == nil || len( orcids ) == 0 {
        status := http.StatusNotFound
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    status := http.StatusOK
    encodeOrcidResponse( w, status, http.StatusText( status ), orcids )
}