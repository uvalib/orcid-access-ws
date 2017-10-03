package orcid

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"orcidaccessws/api"
	"orcidaccessws/config"
	"orcidaccessws/logger"
	"strings"
	"time"
)

var emptyUpdateCode = ""

//
// update the user activity
//
func UpdateOrcidActivity(orcid string, oauth_token string, activity api.ActivityUpdate) (string, int, error) {

	logActivityUpdateRequest(activity)

	// determine if we are creating a new activity or updating an existing one
	existingActivity := len(activity.UpdateCode) != 0

	// construct target URL
	url := fmt.Sprintf("%s/%s/work", config.Configuration.OrcidSecureUrl, orcid)
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
	auth := fmt.Sprintf("Bearer %s", oauth_token)

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
			Timeout(time.Duration(config.Configuration.Timeout) * time.Second).
			End()
	} else {
		resp, body, errs = gorequest.New().
			SetDebug(config.Configuration.Debug).
			Post(url).
			Set("Accept", "application/json").
			Set("Content-Type", "application/xml").
			Set("Authorization", auth).
			Send(requestBody).
			Timeout(time.Duration(config.Configuration.Timeout) * time.Second).
			End()
	}
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus( errs[0] )
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
		httpStatus := mapErrorResponseToStatus( err )
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

//
// renew the access token
//
func RenewAccessToken(staleToken string) (string, string, int, error) {

	// construct target URL
	url := fmt.Sprintf("%s/oauth/token", config.Configuration.OrcidOauthUrl)
	fmt.Printf("%s\n", url)

	// issue the request
	start := time.Now()
	resp, _, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", "application/json").
		Set("refresh_token", staleToken).
		Set("grant_type", "refresh").
		Set("client_id", config.Configuration.OrcidClientId).
		Set("client_secret", config.Configuration.OrcidClientSecret).
		Timeout(time.Duration(config.Configuration.Timeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus( errs[0] )
		return emptyUpdateCode, emptyUpdateCode, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	//logger.Log(fmt.Sprintf("BODY [%s]", body ) )

	return emptyUpdateCode, emptyUpdateCode, http.StatusInternalServerError, errors.New("Not implemented")
}

//
// get details for the specified ORCID
//
func GetOrcidDetails(orcid string) (*api.OrcidDetails, int, error) {

	// construct target URL
	url := fmt.Sprintf("%s/%s/orcid-bio", config.Configuration.OrcidPublicUrl, orcid)
	//fmt.Printf( "%s\n", url )

	// issue the request
	start := time.Now()
	resp, body, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", "application/json").
		Timeout(time.Duration(config.Configuration.Timeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus( errs[0] )
		return nil, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	// check the common response elements
	status, err := checkCommonResponse(body)
	if err != nil {
		httpStatus := mapErrorResponseToStatus( err )
		return nil, httpStatus, err
	}

	if status != http.StatusOK {
		return nil, status, nil
	}

	pr := orcidProfileResponse{}
	err = json.Unmarshal([]byte(body), &pr)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
		return nil, http.StatusInternalServerError, err
	}

	return transformDetailsResponse(&pr.Profile), http.StatusOK, nil
}

//
// search ORCID given the supplied parameters and return the set of ORCID details that match
//
func SearchOrcid(search string, start_ix string, max_results string) ([]*api.OrcidDetails, int, int, error) {

	// construct target URL
	url := fmt.Sprintf("%s/search/orcid-bio?q=%s&start=%s&rows=%s", config.Configuration.OrcidPublicUrl,
		htmlEncodeString(search), start_ix, max_results)
	fmt.Printf("%s\n", url)

	// issue the request
	start := time.Now()
	resp, body, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", "application/json").
		Timeout(time.Duration(config.Configuration.Timeout) * time.Second).
		End()
	duration := time.Since(start)

	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		httpStatus := mapErrorResponseToStatus( errs[0] )
		return nil, 0, httpStatus, errs[0]
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	logger.Log(fmt.Sprintf("Service (%s) returns http %d in %s", url, resp.StatusCode, duration))

	// check the common response elements
	status, err := checkCommonResponse(body)
	if err != nil {
		httpStatus := mapErrorResponseToStatus( err )
		return nil, 0, httpStatus, err
	}

	if status != http.StatusOK {
		return nil, 0, status, err
	}

	sr := orcidSearchResponse{}
	err = json.Unmarshal([]byte(body), &sr)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: json unmarshal: %s", err))
		httpStatus := mapErrorResponseToStatus( err )
		return nil, 0, httpStatus, err
	}

	//if sr.SearchResults.Results == nil || len( sr.SearchResults.Results ) == 0 {
	//    return nil, sr.SearchResults.TotalFound, http.StatusNotFound, nil
	//}

	return transformSearchResponse(sr.SearchResults), sr.SearchResults.TotalFound, http.StatusOK, nil
}

//
// end of file
//
