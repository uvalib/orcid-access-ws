package client

import (
    "time"
    "fmt"
    "github.com/parnurzeal/gorequest"
    "net/http"
    "orcidaccessws/api"
    "encoding/json"
    "io"
    "io/ioutil"
)

const API_DEBUG = false

func HealthCheck( endpoint string ) int {

    url := fmt.Sprintf( "%s/healthcheck", endpoint )
    //fmt.Printf( "%s\n", url )

    resp, _, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Get( url ).
       Timeout( time.Duration( 5 ) * time.Second ).
       End( )

    if errs != nil {
        return http.StatusInternalServerError
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    return resp.StatusCode
}

func VersionCheck( endpoint string ) ( int, string ) {

    url := fmt.Sprintf( "%s/version", endpoint )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
    SetDebug( API_DEBUG ).
    Get( url ).
    Timeout( time.Duration( 5 ) * time.Second ).
    End( )

    if errs != nil {
        return http.StatusInternalServerError, ""
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    r := api.VersionResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, ""
    }

    return resp.StatusCode, r.Version
}

func Statistics( endpoint string ) ( int, * api.Statistics ) {

    url := fmt.Sprintf( "%s/statistics", endpoint )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Get( url ).
       Timeout( time.Duration( 5 ) * time.Second ).
       End( )

    if errs != nil {
        return http.StatusInternalServerError, nil
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    r := api.StatisticsResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, &r.Details
}

func GetOrcid( endpoint string, id string, token string ) ( int, [] * api.Orcid ) {

    url := fmt.Sprintf( "%s/cid/%s?auth=%s", endpoint, id, token )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Get( url ).
       Timeout( time.Duration( 5 ) * time.Second ).
       End( )

    if errs != nil {
       return http.StatusInternalServerError, nil
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    r := api.OrcidResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Orcids
}

func GetAllOrcid( endpoint string, token string ) ( int, [] * api.Orcid ) {

    url := fmt.Sprintf( "%s/cid?auth=%s", endpoint, token )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
            SetDebug( API_DEBUG ).
            Get( url ).
            Timeout( time.Duration( 5 ) * time.Second ).
            End( )

    if errs != nil {
        return http.StatusInternalServerError, nil
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    r := api.OrcidResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Orcids
}

func GetOrcidDetails( endpoint string, orcid string, token string ) ( int, * api.OrcidDetails ) {

    url := fmt.Sprintf( "%s/orcid/%s?auth=%s", endpoint, orcid, token )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
            SetDebug( API_DEBUG ).
            Get( url ).
            Timeout( time.Duration( 5 ) * time.Second ).
            End( )

    if errs != nil {
        return http.StatusInternalServerError, nil
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    r := api.OrcidDetailsResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Details
}

func SearchOrcid( endpoint string, search string, start string, max string, token string ) ( int, [] * api.OrcidDetails, int ) {

    url := fmt.Sprintf( "%s/orcid?q=%s&start=%s&max=%s&auth=%s", endpoint, search, start, max, token )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
            SetDebug( API_DEBUG ).
            Get( url ).
            Timeout( time.Duration( 5 ) * time.Second ).
            End( )

    if errs != nil {
        return http.StatusInternalServerError, nil, 0
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    r := api.OrcidSearchResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil, 0
    }

    return resp.StatusCode, r.Results, r.Total
}

func SetOrcid( endpoint string, cid string, orcid string, token string ) int {

    url := fmt.Sprintf( "%s/cid/%s/%s?auth=%s", endpoint, cid, orcid, token )
    //fmt.Printf( "%s\n", url )

    resp, _, errs := gorequest.New( ).
            SetDebug( API_DEBUG ).
            Put( url ).
            Timeout( time.Duration( 5 ) * time.Second ).
            End( )

    if errs != nil {
        return http.StatusInternalServerError
    }

    defer io.Copy( ioutil.Discard, resp.Body )
    defer resp.Body.Close( )

    return resp.StatusCode
}