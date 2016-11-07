package main

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "orcidaccessws/api"
    "orcidaccessws/orcid"
    "log"
    "strings"
    "orcidaccessws/authtoken"
    "orcidaccessws/config"
    "fmt"
)

func IdLookup( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    doi := vars[ "doi" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    statistics.RequestCount++
    statistics.LookupCount++

    // validate inbound parameters
    if parameterOK( doi ) == false || parameterOK( token ) == false {
        respond( w, http.StatusBadRequest )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "lookup", token ) == false {
        respond( w, http.StatusForbidden )
        return
    }

    entity, status := orcid.GetDoi( doi )
    respondWithDetails( w, status, entity )
}

func IdCreate( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    shoulder := vars[ "shoulder" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    statistics.RequestCount++
    statistics.CreateCount++

    // validate inbound parameters
    if parameterOK( shoulder ) == false || parameterOK( token ) == false {
        respond( w, http.StatusBadRequest )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "create", token ) == false {
        respond( w, http.StatusForbidden )
        return
    }

    decoder := json.NewDecoder( r.Body )
    entity := api.Entity{ }

    if err := decoder.Decode( &entity ); err != nil {
        respond( w, http.StatusBadRequest )
        return
    }

    entity, status := orcid.CreateDoi( shoulder, entity, orcid.STATUS_RESERVED )
    respondWithDetails( w, status, entity )
}

func IdUpdate( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    doi := vars[ "doi" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    statistics.RequestCount++
    statistics.UpdateCount++

    // validate inbound parameters
    if parameterOK( doi ) == false || parameterOK( token ) == false {
        respond( w, http.StatusBadRequest )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "update", token ) == false {
        respond( w, http.StatusForbidden )
        return
    }

    decoder := json.NewDecoder( r.Body )
    entity := api.Entity{ }

    if err := decoder.Decode( &entity ); err != nil {
        respond( w, http.StatusBadRequest )
        return
    }

    entity.Id = doi
    status := orcid.UpdateDoi( entity, orcid.STATUS_PUBLIC )
    respond( w, status )
}

func IdDelete( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    doi := vars[ "doi" ]
    token := r.URL.Query( ).Get( "auth" )

    // update the statistics
    statistics.RequestCount++
    statistics.DeleteCount++

    // validate inbound parameters
    if parameterOK( doi ) == false || parameterOK( token ) == false {
        respond( w, http.StatusBadRequest )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "delete", token ) == false {
        respond( w, http.StatusForbidden )
        return
    }

    status := orcid.DeleteDoi( doi )
    respond( w, status )
}

func IdRevoke( w http.ResponseWriter, r *http.Request ) {

    vars := mux.Vars( r )
    doi := vars[ "doi" ]
    token := r.URL.Query( ).Get( "auth" )

    fmt.Printf( "NEW REVOKE: %s\n", doi )

    // update the statistics
    statistics.RequestCount++
    statistics.RevokeCount++

    // validate inbound parameters
    if parameterOK( doi ) == false || parameterOK( token ) == false {
        respond( w, http.StatusBadRequest )
        return
    }

    // validate the token
    if authtoken.Validate( config.Configuration.AuthTokenEndpoint, "delete", token ) == false {
        respond( w, http.StatusForbidden )
        return
    }

    // get the existing metadata
    entity, status := orcid.GetDoi( doi )
    if status == http.StatusOK {

        // update the status
        entity.Id = doi
        status = orcid.UpdateDoi( entity, orcid.STATUS_UNAVAILABLE )
    }

    respond( w, status )
}

func Stats( w http.ResponseWriter, r *http.Request ) {

    status := http.StatusOK

    jsonResponse( w )
    w.WriteHeader( status )

    if err := json.NewEncoder(w).Encode( api.StatisticsResponse { Status: status, Message: http.StatusText( status ), Details: statistics } ); err != nil {
        log.Fatal( err )
    }
}

func HealthCheck( w http.ResponseWriter, r *http.Request ) {

    // update the statistics
    statistics.RequestCount++
    statistics.HeartbeatCount++

    status := orcid.GetStatus( )
    healthy := status == http.StatusOK
    message := ""

    jsonResponse( w )
    w.WriteHeader( status )

    if err := json.NewEncoder(w).Encode( api.HealthCheckResponse { CheckType: api.HealthCheckResult{ Healthy: healthy, Message: message } } ); err != nil {
        log.Fatal( err )
    }
}

func GetVersion( w http.ResponseWriter, r *http.Request ) {

    jsonResponse( w )
    w.WriteHeader( http.StatusOK )

    if err := json.NewEncoder(w).Encode( api.VersionResponse{ Version: Version( ) } ); err != nil {
        log.Fatal( err )
    }
}

func respond( w http.ResponseWriter, status int ) {

    jsonResponse( w )
    w.WriteHeader( status )
    if err := json.NewEncoder(w).Encode( api.StandardResponse{ Status: status, Message: http.StatusText( status ) } ); err != nil {
        log.Fatal( err )
    }
}

func respondWithDetails( w http.ResponseWriter, status int, entity api.Entity ) {

    jsonResponse( w )
    w.WriteHeader( status )
    if status == http.StatusOK {
        if err := json.NewEncoder(w).Encode( api.StandardResponse{ Status: status, Message: http.StatusText( status ), Details: &entity } ); err != nil {
            log.Fatal( err )
        }
    } else {
        if err := json.NewEncoder(w).Encode( api.StandardResponse{ Status: status, Message: http.StatusText( status ) } ); err != nil {
            log.Fatal( err )
        }
    }
}

func jsonResponse( w http.ResponseWriter ) {
    w.Header( ).Set( "Content-Type", "application/json; charset=UTF-8" )
}

func parameterOK( param string ) bool {
    return len( strings.TrimSpace( param ) ) != 0
}
