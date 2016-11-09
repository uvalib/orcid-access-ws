package handlers

import (
    "log"
    //"fmt"
    "strings"
    "encoding/json"
    "net/http"
    "orcidaccessws/api"
    //"orcidaccessws/mapper"
    //"orcidaccessws/logger"
)

func EncodeStandardResponse( w http.ResponseWriter, status int, message string ) {

    //logger.Log( fmt.Sprintf( "Status: %d (%s)\n", status, message ) )
    jsonAttributes( w )
//    coorsAttributes( w )
    w.WriteHeader( status )
    if err := json.NewEncoder(w).Encode( api.StandardResponse{ Status: status, Message: message } ); err != nil {
        log.Fatal( err )
    }
}

func EncodeOrcidResponse( w http.ResponseWriter, status int, message string, orcids [] * api.Orcid ) {

    //logger.Log( fmt.Sprintf( "Status: %d (%s)\n", status, message ) )
    jsonAttributes( w )
    //    coorsAttributes( w )
    w.WriteHeader( status )
    if err := json.NewEncoder(w).Encode( api.OrcidResponse{ Status: status, Message: message, Orcids: orcids } ); err != nil {
        log.Fatal( err )
    }
}

func EncodeHealthCheckResponse( w http.ResponseWriter, status int, dbmsg string, orcidmsg string ) {

    db_healthy, orcid_healthy := true, true
    if len( dbmsg ) != 0 {
        db_healthy = false
    }
    if len( orcidmsg ) != 0 {
        orcid_healthy = false
    }
    jsonAttributes( w )
    w.WriteHeader( status )
    if err := json.NewEncoder(w).Encode( api.HealthCheckResponse {
        DbCheck: api.HealthCheckResult{ Healthy: db_healthy, Message: dbmsg },
        OrcidCheck: api.HealthCheckResult{ Healthy: orcid_healthy, Message: orcidmsg } } ); err != nil {
        log.Fatal( err )
    }
}

func EncodeStatsResponse( w http.ResponseWriter ) {

    status := http.StatusOK

    jsonAttributes( w )
    w.WriteHeader( status )

    //if err := json.NewEncoder(w).Encode( api.StatisticsResponse { Status: status, Message: http.StatusText( status ), Details: statistics } ); err != nil {
    //    log.Fatal( err )
    //}

}

func encodeVersionResponse( w http.ResponseWriter, status int, version string ) {
    jsonAttributes( w )
    w.WriteHeader( status )
    if err := json.NewEncoder(w).Encode( api.VersionResponse { Version: version } ); err != nil {
        log.Fatal( err )
    }
}

func jsonAttributes( w http.ResponseWriter ) {
    w.Header( ).Set( "Content-Type", "application/json; charset=UTF-8" )
}

func coorsAttributes( w http.ResponseWriter ) {
    w.Header( ).Set( "Access-Control-Allow-Origin", "*" )
    w.Header( ).Set( "Access-Control-Allow-Headers", "Content-Type" )
}

func NotEmpty( param string ) bool {
    return len( strings.TrimSpace( param ) ) != 0
}