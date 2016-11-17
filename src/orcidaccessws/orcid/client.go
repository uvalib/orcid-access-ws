package orcid

import (
    "fmt"
    "time"
    "net/http"
    "github.com/parnurzeal/gorequest"
    "orcidaccessws/api"
    "orcidaccessws/config"
    "orcidaccessws/logger"
    "encoding/json"
)

//
// v1.2 response structure
//

const PROTOCOL_VERSION = "1.2"

type orcidQueryResponse struct {
    Version             string                 `json:"message-version,omitempty"`
    Profile             orcidProfile           `json:"orcid-profile,omitempty"`
}

type orcidVersionResponse struct {
    Version             string                 `json:"message-version,omitempty"`
}

type orcidProfile struct {
    Bio                 orcidBio               `json:"orcid-bio,omitempty"`
}

type orcidBio struct {
    PersonalDetails     orcidPersonalDetails   `json:"orcid-bio,omitempty"`
    Biography           valueVisibilityPair    `json:"biography,omitempty"`
}

type orcidPersonalDetails struct {
    GivenName           valueVisibilityPair    `json:"given-names,omitempty"`
    FamilyName          valueVisibilityPair    `json:"family-name,omitempty"`
    CreditName          valueVisibilityPair    `json:"credit-name,omitempty"`
}

type valueVisibilityPair struct {
    Value               string                 `json:"value,omitempty"`
    Visibility          string                 `json:"visibility,omitempty"`
}

// status for the EZID objects
const STATUS_PUBLIC = "public"
const STATUS_RESERVED = "reserved"
const STATUS_UNAVAILABLE = "unavailable"

const USE_CROSS_REF_PROFILE = true

func GetOrcidDetails( orcid string ) ( [] * api.OrcidDetails, int ) {

    // construct target URL
    url := fmt.Sprintf( "%s/%s", config.Configuration.OrcidServiceUrl, orcid )
    fmt.Printf( "%s\n", url )

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
       SetDebug( config.Configuration.Debug ).
       Get( url  ).
       Timeout( time.Duration( config.Configuration.OrcidServiceTimeout ) * time.Second ).
       Set( "Accept", "application/json" ).
       End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return nil, http.StatusInternalServerError
    }

    defer resp.Body.Close( )

    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the protocol version to ensure we support it
    
    vr := orcidVersionResponse{ }
    err := json.Unmarshal( []byte( body ), &vr )
    if err != nil {
        return nil, http.StatusInternalServerError
    }

    if vr.Version != PROTOCOL_VERSION {
        logger.Log( fmt.Sprintf( "ORCID protocol version not supported. Require: %s, received: %s", PROTOCOL_VERSION, vr.Version ) )
        return nil, http.StatusHTTPVersionNotSupported
    }

    qr := orcidQueryResponse{ }
    err = json.Unmarshal( []byte( body ), &qr )
    if err != nil {
        return nil, http.StatusInternalServerError
    }

    //response := transformDetailsResponse( r.Profile )
    return nil, http.StatusNotFound
}

func SearchOrcid( search string ) ( [] * api.OrcidDetails, int ) {
    return nil, http.StatusNotFound
}

//
//
//

//
// get entity details when provided a DOI
//
func GetDoi( doi string ) ( api.Entity, int ) {

    // construct target URL
    url := fmt.Sprintf( "%s/id/%s", config.Configuration.OrcidServiceUrl, doi )

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
        SetDebug( config.Configuration.Debug ).
        Get( url  ).
        Timeout( time.Duration( config.Configuration.OrcidServiceTimeout ) * time.Second ).
        End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return blankEntity( ), http.StatusInternalServerError
    }

    defer resp.Body.Close( )

    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the body for errors
    if !statusIsOk( body ) {
        logger.Log( fmt.Sprintf( "Error response body: [%s]", body ) )
        return blankEntity( ), http.StatusBadRequest
    }

    // all good...
    return makeEntityFromBody( body ), http.StatusOK
}

//
// Create a new entity; we may or may not have complete entity details
//
func CreateDoi( shoulder string, entity api.Entity, status string ) ( api.Entity, int ) {

    // log if necessary
    logEntity( entity )

    // construct target URL
    url := fmt.Sprintf( "%s/shoulder/%s", config.Configuration.OrcidServiceUrl, shoulder )

    var body string
    var err error
    // construct the payload, set the status to reserved
    if USE_CROSS_REF_PROFILE == true {
        body, err = makeCrossRefBodyFromEntity(entity, status )
    } else {
        body, err = makeDataciteBodyFromEntity( entity, status )
    }

    // check for errors
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: creating service payload %s", err ) )
        return blankEntity( ), http.StatusInternalServerError
    }

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
        SetDebug( config.Configuration.Debug ).
        SetBasicAuth( config.Configuration.OrcidUser, config.Configuration.OrcidPassphrase ).
        Post( url  ).
        Send( body ).
        Timeout( time.Duration( config.Configuration.OrcidServiceTimeout ) * time.Second ).
        Set( "Content-Type", "text/plain" ).
        End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return blankEntity( ), http.StatusInternalServerError
    }

    defer resp.Body.Close( )

    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the body for errors
    if !statusIsOk( body ) {
        logger.Log( fmt.Sprintf( "Error response body: [%s]", body ) )
        return blankEntity( ), http.StatusBadRequest
    }

    // all good...
    return makeEntityFromBody( body ), http.StatusOK
}

//
// Update an existing DOI to match the provided entity
//
func UpdateDoi( entity api.Entity, status string ) int {

    // log if necessary
    logEntity( entity )

    // construct target URL
    url := fmt.Sprintf( "%s/id/%s", config.Configuration.OrcidServiceUrl, entity.Id )

    var body string
    var err error
    // construct the payload...
    if USE_CROSS_REF_PROFILE == true {
        body, err = makeCrossRefBodyFromEntity(entity, status )
    } else {
        body, err = makeDataciteBodyFromEntity( entity, status )
    }

    // check for errors
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: creating service payload %s", err ) )
        return http.StatusInternalServerError
    }

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
        SetDebug( config.Configuration.Debug ).
        SetBasicAuth( config.Configuration.OrcidUser, config.Configuration.OrcidPassphrase ).
        Post( url  ).
        Send( body ).
        Timeout( time.Duration( config.Configuration.OrcidServiceTimeout ) * time.Second ).
        Set( "Content-Type", "text/plain" ).
        End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return http.StatusInternalServerError
    }

    defer resp.Body.Close( )
    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the body for errors
    if !statusIsOk( body ) {
        logger.Log( fmt.Sprintf( "Error response body: [%s]", body ) )
        return http.StatusBadRequest
    }

    // all good...
    return http.StatusOK
}

//
// delete entity details when provided a DOI
//
func DeleteDoi( doi string ) int {

    // construct target URL
    url := fmt.Sprintf( "%s/id/%s", config.Configuration.OrcidServiceUrl, doi )

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
        SetDebug( config.Configuration.Debug ).
        SetBasicAuth( config.Configuration.OrcidUser, config.Configuration.OrcidPassphrase ).
        Delete( url  ).
        Timeout( time.Duration( config.Configuration.OrcidServiceTimeout ) * time.Second ).
        End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return http.StatusInternalServerError
    }

    defer resp.Body.Close( )
    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the body for errors
    if !statusIsOk( body ) {
        logger.Log( fmt.Sprintf( "Error response body: [%s]", body ) )
        return http.StatusBadRequest
    }

    // all good...
    return http.StatusOK
}

//
// get the status of the endpoint
//
func GetStatus( ) int {

    // construct target URL
    url := fmt.Sprintf( "%s/status", config.Configuration.OrcidServiceUrl )

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
        SetDebug( config.Configuration.Debug ).
        Get( url  ).
        Timeout( time.Duration( config.Configuration.OrcidServiceTimeout ) * time.Second ).
        End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return http.StatusInternalServerError
    }

    defer resp.Body.Close( )
    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the body for errors
    if !statusIsOk( body ) {
        logger.Log( fmt.Sprintf( "Error response body: [%s]", body ) )
        return http.StatusBadRequest
    }

    // all good...
    return http.StatusOK
}