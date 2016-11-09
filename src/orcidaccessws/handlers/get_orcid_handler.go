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

    // parameters OK ?
    if NotEmpty( id ) == false || NotEmpty( token ) == false {
        status := http.StatusBadRequest
        EncodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "getorcid", token ) == false {
        status := http.StatusForbidden
        EncodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    // get the authorization details
    orcids, err := dao.Database.GetOrcidByCid( id )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: %s\n", err.Error( ) ) )
        status := http.StatusInternalServerError
        EncodeOrcidResponse( w, status,
            fmt.Sprintf( "%s (%s)", http.StatusText( status ), err ),
            nil )
        return
    }

    // we did not find the item, return 404
    if orcids == nil || len( orcids ) == 0 {
        status := http.StatusNotFound
        EncodeOrcidResponse( w, status, http.StatusText( status ), nil )
        return
    }

    status := http.StatusOK
    EncodeOrcidResponse( w, status, http.StatusText( status ), orcids )
}