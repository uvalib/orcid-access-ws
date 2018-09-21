package orcid

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
	"strings"
	"time"
)

var emptyUpdateCode = ""
var emptyAccessToken = ""
var currentAccessToken = ""

//
// UpdateOrcidActivity -- update the user activity
//
func UpdateOrcidActivity(orcid string, oauthToken string, activity api.ActivityUpdate) (string, int, error) {

	logActivityUpdateRequest(activity)

	// determine if we are creating a new activity or updating an existing one
	existingActivity := len(activity.UpdateCode) != 0

	// construct target URL
	url := fmt.Sprintf("%s/%s/work", config.Configuration.OrcidSecureURL, orcid)
	if existingActivity == true {
		url = fmt.Sprintf("%s/%s", url, activity.UpdateCode)
	}
	//fmt.Printf( "%s\n", url )

	// build the request body
	requestBody, err := makeUpdateActivityBody(activity)

	// check for errors
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: creating service payload %s", err))
		return emptyUpdateCode, http.StatusBadRequest, err
	}

	// construct the auth field
	auth := fmt.Sprintf("Bearer %s", oauthToken)

	// issue the request
	start := time.Now()
	var resp gorequest.Response
	var body string
	var errs []error
	if existingActivity == true {
		resp, body, errs = gorequest.New().
			SetDebug(config.Configuration.Debug).
			Put(url).
			Set("Accept", "application/json").
			Set("Content-Type", "application/xml").
			Set("Authorization", auth).
			Send(requestBody).
			Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
			End()
	} else {
		resp, body, errs = gorequest.New().
			SetDebug(config.Configuration.Debug).
			Post(url).
			Set("Accept", "application/json").
			Set("Content-Type", "application/xml").
			Set("Authorization", auth).
			Send(requestBody).
			Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
			End()
	}
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus(errs[0])
		return emptyUpdateCode, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	// the happy update path; just return the original update code
	if existingActivity == true && resp.StatusCode == http.StatusOK {
		return activity.UpdateCode, http.StatusOK, nil
	}

	// the happy create path, pull the push code from the location header
	if existingActivity == false && resp.StatusCode == http.StatusCreated {
		tokens := strings.Split(resp.Header.Get("Location"), "/")
		if len(tokens) != 0 {
			return tokens[len(tokens)-1], http.StatusOK, nil
		}

		// unexpected, return an error
		return emptyUpdateCode, http.StatusInternalServerError, errors.New("Unexpected/missing location header in response")
	}

	//
	// something unexpected happened and we did not get an error report (handled above)
	//

	// check for a 500 error and return if we find it
	if resp.StatusCode == http.StatusInternalServerError {
		return emptyUpdateCode, http.StatusInternalServerError, errors.New("Server reports Internal Server Error")
	}

	// otherwise, attempt to decode the response
	aur := activityUpdateResponse{}
	err = json.Unmarshal([]byte(body), &aur)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
		httpStatus := mapErrorResponseToStatus(err)
		return emptyUpdateCode, httpStatus, err
	}

	//
	// ORCID reported an error
	//
	if len(aur.Error) != 0 {
		logger.Log(fmt.Sprintf("ERROR: service reports %s (%s)", aur.Error, aur.ErrorDescription))
		return emptyUpdateCode, resp.StatusCode, errors.New(aur.ErrorDescription)
	}

	if aur.ResponseCode != http.StatusOK {
		logger.Log(fmt.Sprintf("ERROR: service reports: %d (%s)", aur.ResponseCode, aur.DeveloperMessage))
		return emptyUpdateCode, aur.ResponseCode, errors.New(aur.DeveloperMessage)
	}

	// unclear why we are here but an error occurred
	return emptyUpdateCode, http.StatusInternalServerError, errors.New("Unhandled error case")
}

func getOauthToken() (string, int, error) {

	// construct target URL
	url := fmt.Sprintf("%s/oauth/token", config.Configuration.OrcidOauthURL)
	//fmt.Printf("%s\n", url)

	// create the request payload
	pl := oauthRequest{
		ClientID:     config.Configuration.OrcidClientID,
		ClientSecret: config.Configuration.OrcidClientSecret,
		Scope:        "/read-public",
		GrantType:    "client_credentials"}

	// issue the request
	start := time.Now()
	resp, body, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Post(url).
		Send(pl).
		Set("Accept", "application/json").
		Set("Content-Type", "application/x-www-form-urlencoded").
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus(errs[0])
		return emptyAccessToken, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	// check for a non-success status and return if we find it
	if resp.StatusCode != http.StatusOK {
		return emptyAccessToken, resp.StatusCode, fmt.Errorf("Service (%s) returns http %d", url, resp.StatusCode)
	}

	// otherwise, attempt to decode the response
	oar := oauthResponse{}
	err := json.Unmarshal([]byte(body), &oar)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
		httpStatus := mapErrorResponseToStatus(err)
		return emptyAccessToken, httpStatus, err
	}

	return oar.AccessToken, http.StatusOK, nil
}

//
// RenewAccessToken -- renew the access token
//
func RenewAccessToken(staleToken string) (string, string, int, error) {

	// construct target URL
	url := fmt.Sprintf("%s/oauth/token", config.Configuration.OrcidOauthURL)
	//fmt.Printf("%s\n", url)

	// issue the request
	start := time.Now()
	resp, _, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", "application/json").
		Set("refresh_token", staleToken).
		Set("grant_type", "refresh").
		Set("client_id", config.Configuration.OrcidClientID).
		Set("client_secret", config.Configuration.OrcidClientSecret).
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus(errs[0])
		return emptyUpdateCode, emptyUpdateCode, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	//logger.Log(fmt.Sprintf("BODY [%s]", body ) )

	return emptyUpdateCode, emptyUpdateCode, http.StatusInternalServerError, errors.New("Not implemented")
}

//
// GetOrcidDetails -- get details for the specified ORCID
//
func GetOrcidDetails(orcid string) (*api.OrcidDetails, int, error) {

	// construct target URL
	url := fmt.Sprintf("%s/%s/person", config.Configuration.OrcidPublicURL, orcid)
	//fmt.Printf( "%s\n", url )

	// get an access token
	token, err := getAccessToken()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// issue the request
	start := time.Now()
	resp, body, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Authorization", token).
		Set("Accept", "application/json").
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus(errs[0])
		return nil, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	// check for an http status
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New( orcid )
	}

	pr := orcidPersonResponse{}
	err = json.Unmarshal([]byte(body), &pr)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
		return nil, http.StatusInternalServerError, err
	}

	return transformDetailsResponse(&pr), http.StatusOK, nil
}

//
// SearchOrcid -- search ORCID given the supplied parameters and return the set of ORCID details that match
//
//func SearchOrcid(search string, startIx string, maxResults string) ([]*api.OrcidDetails, int, int, error) {
//
//	// construct target URL
//	url := fmt.Sprintf("%s/search?q=%s&start=%s&rows=%s", config.Configuration.OrcidPublicURL,
//		htmlEncodeString(search), startIx, maxResults)
//	fmt.Printf("%s\n", url)
//
//	// issue the request
//	start := time.Now()
//	resp, body, errs := gorequest.New().
//		SetDebug(config.Configuration.Debug).
//		Get(url).
//		Set("Accept", "application/json").
//		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
//		End()
//	duration := time.Since(start)
//
//	// check for errors
//	if errs != nil {
//		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
//		httpStatus := mapErrorResponseToStatus(errs[0])
//		return nil, 0, httpStatus, errs[0]
//	}
//
//	defer io.Copy(ioutil.Discard, resp.Body)
//	defer resp.Body.Close()
//
//	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))
//
//	// check the common response elements
//	status, err := checkCommonResponse(body)
//	if err != nil {
//		httpStatus := mapErrorResponseToStatus(err)
//		return nil, 0, httpStatus, err
//	}
//
//	if status != http.StatusOK {
//		return nil, 0, status, err
//	}
//
//	sr := orcidSearchResponse{}
//	err = json.Unmarshal([]byte(body), &sr)
//	if err != nil {
//		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
//		httpStatus := mapErrorResponseToStatus(err)
//		return nil, 0, httpStatus, err
//	}
//
//	//if sr.SearchResults.Results == nil || len( sr.SearchResults.Results ) == 0 {
//	//    return nil, sr.SearchResults.TotalFound, http.StatusNotFound, nil
//	//}
//
//	return transformSearchResponse(sr.SearchResults), sr.SearchResults.TotalFound, http.StatusOK, nil
//}

//
// GetPublicEndpointStatus -- get the public endpoint status
//
func GetPublicEndpointStatus() error {

	// construct target URL
	url := fmt.Sprintf("%s/status", config.Configuration.OrcidPublicURL)
	//fmt.Printf( "%s\n", url )

	// get an access token
	token, err := getAccessToken()
	if err != nil {
		return err
	}

	return issueAuthorizedGet(url, "text/plain", token)
}

//
// GetSecureEndpointStatus -- get the secure endpoint status
//
func GetSecureEndpointStatus() error {

	// construct target URL
	url := fmt.Sprintf("%s/status", config.Configuration.OrcidSecureURL)
	//fmt.Printf( "%s\n", url )

	// get the access token
	token, err := getAccessToken()
	if err != nil {
		return err
	}
	return issueAuthorizedGet(url, "text/plain", token)
}

//
// issue a GET to the specified URL and use a bearer token for aurthorization
// ignore the response payload and return an error if we get a non-200 response
//
func issueAuthorizedGet(url string, accept string, authToken string) error {

	// construct the auth field
	auth := fmt.Sprintf("Bearer %s", authToken)

	// issue the request
	start := time.Now()
	resp, _, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", accept).
		Set("Authorization", auth).
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		return errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	// check for a non-success status and return if we find it
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Service (%s) returns http %d", url, resp.StatusCode)
	}

	return nil
}

//
// singleton type access to the access token
//
func getAccessToken() (string, error) {

	// do we need an access token
	if len(currentAccessToken) == 0 {
		token, status, err := getOauthToken()
		if status != http.StatusOK {
			return emptyAccessToken, err
		}
		currentAccessToken = token
	}
	return currentAccessToken, nil
}

//
// end of file
//
