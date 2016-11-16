package client

import (
    "time"
    "fmt"
    "github.com/parnurzeal/gorequest"
    "net/http"
    "orcidaccessws/api"
    "encoding/json"
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

    defer resp.Body.Close( )

    r := api.StatisticsResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, &r.Details
}

func GetOneOrcid( endpoint string, id string, token string ) ( int, [] * api.Orcid ) {

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

    defer resp.Body.Close( )

    r := api.OrcidResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Orcids
}

func GetOrcidDetails( endpoint string, orcid string, token string ) ( int, [] * api.OrcidDetails ) {

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

    defer resp.Body.Close( )

    r := api.OrcidDetailsResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Details
}

func SearchOrcid( endpoint string, search string, token string ) ( int, [] * api.OrcidDetails ) {

    url := fmt.Sprintf( "%s/orcid?q=%s&auth=%s", endpoint, search, token )
    //fmt.Printf( "%s\n", url )

    resp, body, errs := gorequest.New( ).
            SetDebug( API_DEBUG ).
            Get( url ).
            Timeout( time.Duration( 5 ) * time.Second ).
            End( )

    if errs != nil {
        return http.StatusInternalServerError, nil
    }

    defer resp.Body.Close( )

    r := api.OrcidDetailsResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Details
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

    defer resp.Body.Close( )

    return resp.StatusCode
}

/*
func Create( endpoint string, shoulder string, token string ) ( int, * api.Entity ) {

    url := fmt.Sprintf("%s/%s?auth=%s", endpoint, shoulder, token)
    //fmt.Printf( "%s\n", url )

    entity := api.Entity{ Title : "my title", Url: "http://google.com" }

    resp, body, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Post( url ).
       Send( entity ).
       Timeout( time.Duration( 5 ) * time.Second ).
       Set( "Content-Type", "application/json" ).
       End( )

    if errs != nil {
        return http.StatusInternalServerError, nil
    }

    defer resp.Body.Close( )

    r := api.StandardResponse{ }
    err := json.Unmarshal( []byte( body ), &r )
    if err != nil {
        return http.StatusInternalServerError, nil
    }

    return resp.StatusCode, r.Details
}

func Update( endpoint string, entity api.Entity, token string ) int {

    url := fmt.Sprintf("%s/%s?auth=%s", endpoint, entity.Id, token )
    //fmt.Printf( "%s\n", url )

    resp, _, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Put( url ).
       Send( entity ).
       Timeout( time.Duration( 5 ) * time.Second ).
       Set( "Content-Type", "application/json" ).
       End( )

    if errs != nil {
        return http.StatusInternalServerError
    }

    defer resp.Body.Close( )

    return resp.StatusCode
}

func Delete( endpoint string, doi string, token string ) int {

    url := fmt.Sprintf("%s/%s?auth=%s", endpoint, doi, token )
    //fmt.Printf( "%s\n", url )

    resp, _, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Delete( url ).
       Timeout( time.Duration( 5 ) * time.Second ).
       End( )

    if errs != nil {
        return http.StatusInternalServerError
    }

    defer resp.Body.Close( )

    return resp.StatusCode
}

//
// revoke an entity when provided a DOI
//
func Revoke( endpoint string, doi string, token string ) int {

    // construct target URL
    url := fmt.Sprintf("%s/revoke/%s?auth=%s", endpoint, doi, token )
    //fmt.Printf( "%s\n", url )

    resp, _, errs := gorequest.New( ).
       SetDebug( API_DEBUG ).
       Put( url ).
       Timeout( time.Duration( 5 ) * time.Second ).
       End( )

    if errs != nil {
        return http.StatusInternalServerError
    }

    defer resp.Body.Close( )

    return resp.StatusCode
}

*/
