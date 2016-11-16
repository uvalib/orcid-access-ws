package handlers

import (
    "fmt"
    "net/http"
    "orcidaccessws/authtoken"
    "orcidaccessws/config"
    "orcidaccessws/dao"
    "orcidaccessws/logger"
)

func GetAllOrcid( w http.ResponseWriter, r *http.Request ) {

    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    Statistics.RequestCount++

    // parameters OK ?
    if NotEmpty( token ) == false {
        status := http.StatusBadRequest
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "getorcid", token ) == false {
        status := http.StatusForbidden
        encodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // get the ORCID details
    orcids, err := dao.Database.GetAllOrcid( )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: %s\n", err.Error( ) ) )
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