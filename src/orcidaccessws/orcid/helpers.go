package orcid

import (
    "fmt"
    "orcidaccessws/api"
    "encoding/json"
    "net/http"
    "orcidaccessws/logger"
    "strings"
    "errors"
    "html"
)

func checkCommonResponse( body string ) ( int, error ) {

    cr := orcidCommonResponse{ }
    err := json.Unmarshal( []byte( body ), &cr )
    if err != nil {
        logger.Log( fmt.Sprintf( "ERROR: json unmarshal: %s", err ) )
        return http.StatusInternalServerError, err
    }

    // check protocol version to ensure we know what to do with this
    if cr.Version != PROTOCOL_VERSION {
        logger.Log( fmt.Sprintf( "ORCID protocol version not supported. Require: %s, received: %s", PROTOCOL_VERSION, cr.Version ) )
        return http.StatusHTTPVersionNotSupported, nil
    }

    // is there an error string
    if cr.Error.Value != "" {
       if strings.HasPrefix( cr.Error.Value, "Not found" ) == true {
           return http.StatusNotFound, nil
       }

       // not sure, just return a general error
       return http.StatusInternalServerError, errors.New( cr.Error.Value )
    }

    return http.StatusOK, nil
}

func transformDetailsResponse( profile orcidProfile ) [] * api.OrcidDetails {
    results := make([ ] * api.OrcidDetails, 0 )
    return append( results, constructDetails( profile ) )
}

func transformSearchResponse( search orcidResults ) [] * api.OrcidDetails {
    results := make([ ] * api.OrcidDetails, 0 )
    for _, e := range search.Results {
        od := constructDetails( e.Profile )
        od.Relevancy = fmt.Sprintf( "%.6f", e.Relevancy.Value )
        results = append( results, od )
    }
    return( results )
}

func constructDetails( profile orcidProfile ) * api.OrcidDetails {

    od := new( api.OrcidDetails )

    od.Id = profile.Id.Id
    od.Uri = profile.Id.Uri
    od.DisplayName = profile.Bio.PersonalDetails.CreditName.Value
    od.FirstName = profile.Bio.PersonalDetails.GivenName.Value
    od.LastName = profile.Bio.PersonalDetails.FamilyName.Value
    od.Biography = profile.Bio.Biography.Value

    od.Keywords = make([ ] string, 0 )
    for _, e := range profile.Bio.Keywords.Keywords {
        od.Keywords = append( od.Keywords, e.Value )
    }

    od.ResearchUrls = make([ ] string, 0 )
    for _, e := range profile.Bio.Urls.Urls {
        od.ResearchUrls = append( od.ResearchUrls, e.Url.Value )
    }

    return( od )
}

func htmlEncode( value string ) string {
    // HTML encoding
    return html.EscapeString( value )
}