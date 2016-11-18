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

func GetOrcidDetails( orcid string ) ( [] * api.OrcidDetails, int ) {

    // construct target URL
    url := fmt.Sprintf( "%s/%s/orcid-bio", config.Configuration.OrcidServiceUrl, orcid )
    //fmt.Printf( "%s\n", url )

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

    // check the common response elements
    status, err := checkCommonResponse( body )
    if err != nil {
        return nil, http.StatusInternalServerError
    }

    if status != http.StatusOK {
        return nil, status
    }

    pr := orcidProfileResponse{ }
    err = json.Unmarshal( []byte( body ), &pr )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: json unmarshal: %s", err ) )
        return nil, http.StatusInternalServerError
    }

    return transformDetailsResponse( pr.Profile ), http.StatusOK
}

func SearchOrcid( search string ) ( [] * api.OrcidDetails, int ) {

    // construct target URL
    url := fmt.Sprintf( "%s/search/orcid-bio?q=%s&start=0&rows=1", config.Configuration.OrcidServiceUrl, htmlEncode( search ) )
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

    // check the common response elements
    status, err := checkCommonResponse( body )
    if err != nil {
        return nil, http.StatusInternalServerError
    }

    if status != http.StatusOK {
        return nil, status
    }

    sr := orcidSearchResponse{ }
    err = json.Unmarshal( []byte( body ), &sr )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: json unmarshal: %s", err ) )
        return nil, http.StatusInternalServerError
    }

    if sr.SearchResults.TotalFound == 0 {
        return nil, http.StatusNotFound
    }

    return transformSearchResponse( sr.SearchResults ), http.StatusOK
}
