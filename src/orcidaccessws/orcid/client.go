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
    "io"
    "io/ioutil"
)

func GetOrcidDetails( orcid string ) ( * api.OrcidDetails, int, error ) {

    // construct target URL
    url := fmt.Sprintf( "%s/%s/orcid-bio", config.Configuration.OrcidServiceUrl, orcid )
    //fmt.Printf( "%s\n", url )

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
       SetDebug( config.Configuration.Debug ).
       Get( url  ).
       Timeout( time.Duration( config.Configuration.Timeout ) * time.Second ).
       Set( "Accept", "application/json" ).
       End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return nil, http.StatusInternalServerError, errs[ 0 ]
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the common response elements
    status, err := checkCommonResponse( body )
    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

    if status != http.StatusOK {
        return nil, status, nil
    }

    pr := orcidProfileResponse{ }
    err = json.Unmarshal( []byte( body ), &pr )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: json unmarshal: %s", err ) )
        return nil, http.StatusInternalServerError, err
    }

    return transformDetailsResponse( &pr.Profile ), http.StatusOK, nil
}

func SearchOrcid( search string, start_ix string, max_results string ) ( [] * api.OrcidDetails, int, int, error ) {

    // construct target URL
    url := fmt.Sprintf( "%s/search/orcid-bio?q=%s&start=%s&rows=%s", config.Configuration.OrcidServiceUrl,
        htmlEncode( search ), start_ix, max_results )
    fmt.Printf( "%s\n", url )

    // issue the request
    start := time.Now( )
    resp, body, errs := gorequest.New( ).
            SetDebug( config.Configuration.Debug ).
            Get( url  ).
            Timeout( time.Duration( config.Configuration.Timeout ) * time.Second ).
            Set( "Accept", "application/json" ).
            End( )
    duration := time.Since( start )

    // check for errors
    if errs != nil {
        logger.Log( fmt.Sprintf( "ERROR: service (%s) returns %s in %s", url, errs, duration ) )
        return nil, 0, http.StatusInternalServerError, errs[ 0 ]
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    logger.Log( fmt.Sprintf( "Service (%s) returns http %d in %s", url, resp.StatusCode, duration ) )

    // check the common response elements
    status, err := checkCommonResponse( body )
    if err != nil {
        return nil, 0, http.StatusInternalServerError, err
    }

    if status != http.StatusOK {
        return nil, 0, status, err
    }

    sr := orcidSearchResponse{ }
    err = json.Unmarshal( []byte( body ), &sr )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: json unmarshal: %s", err ) )
        return nil, 0, http.StatusInternalServerError, err
    }

    //if sr.SearchResults.Results == nil || len( sr.SearchResults.Results ) == 0 {
    //    return nil, sr.SearchResults.TotalFound, http.StatusNotFound, nil
    //}

    return transformSearchResponse( sr.SearchResults ), sr.SearchResults.TotalFound, http.StatusOK, nil
}
