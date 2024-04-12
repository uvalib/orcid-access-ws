package orcid

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
)

var emptyUpdateCode = ""
var emptyAccessToken = ""
var currentAccessToken = ""

// we want to handle this differently from other requests because it is used as part of healthchecking
var authTimeout = 5 * time.Second

// UpdateOrcidActivity -- update the user activity
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

var employmentXML = `<?xml version="1.0" encoding="UTF-8"?>
<employment:employment
	xmlns:employment="http://www.orcid.org/ns/employment" xmlns:common="http://www.orcid.org/ns/common"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.orcid.org/ns/employment ../employment-3.0.xsd ">
	<common:organization>
		<common:name>University of Virginia</common:name>
		<common:address>
			<common:city>Charlottesville</common:city>
			<common:region>VA</common:region>
			<common:country>US</common:country>
		</common:address>
		<common:disambiguated-organization>
				<common:disambiguated-organization-identifier>https://ror.org/0153tk833</common:disambiguated-organization-identifier>
				<common:disambiguation-source>ROR</common:disambiguation-source>
			</common:disambiguated-organization>
	</common:organization>
</employment:employment>`

// SendEmployment updates the user's ORCID with UVA Employent info
func SendEmployment(attributes api.OrcidAttributes) {

	// First check for existing UVA Employment
	if hasEmp, err := hasExistingEmployment(attributes); err != nil || hasEmp {
		return
	}

	url := fmt.Sprintf("%s/%s/employment", config.Configuration.OrcidSecureURL, attributes.Orcid)
	// construct the auth field
	auth := fmt.Sprintf("Bearer %s", attributes.OauthAccessToken)
	start := time.Now()
	resp, _, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Post(url).
		Set("Accept", "application/xml").
		Set("Content-Type", "application/xml").
		Set("Authorization", auth).
		Send(employmentXML).
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	// check for errors
	if errs != nil || resp.StatusCode != http.StatusCreated {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s %s in %s \n %s", url, resp.Status, errs, duration, resp.Body))
		return
	}
	logger.Log(fmt.Sprintf("Employment created for %s", attributes.Orcid))

	return
}

// hasExistingEmployment checks for existing UVA Employment matching our OrcidClientID
func hasExistingEmployment(attributes api.OrcidAttributes) (bool, []error) {
	hasEmployment := false
	// Only need to know employment here
	var employmentStruct struct {
		AffiliationGroup []struct {
			Summaries []struct {
				Employment struct {
					Source struct {
						SourceClientID struct {
							Path string `json:"path"`
						} `json:"source-client-id"`
					} `json:"source"`
				} `json:"employment-summary"`
			} `json:"summaries"`
		} `json:"affiliation-group"`
	}
	url := fmt.Sprintf("%s/%s/employments", config.Configuration.OrcidSecureURL, attributes.Orcid)
	// construct the auth field
	auth := fmt.Sprintf("Bearer %s", attributes.OauthAccessToken)
	start := time.Now()
	resp, body, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", "application/json").
		Set("Content-Type", "application/json").
		Set("Authorization", auth).
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	// check for errors
	if errs != nil || resp.StatusCode != http.StatusOK {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s %s in %s \n %s", url, resp.Status, errs, duration, resp.Body))
		return false, errs
	}

	if err := json.Unmarshal([]byte(body), &employmentStruct); err != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		return false, errs
	}

	// Check the nested json
	for _, g := range employmentStruct.AffiliationGroup {
		for _, s := range g.Summaries {
			// Check if existing employments match our client ID
			if s.Employment.Source.SourceClientID.Path == config.Configuration.OrcidClientID {
				hasEmployment = true
			}
		}
	}
	logger.Log(fmt.Sprintf("INFO: %s has UVA Employment: %t", attributes.Orcid, hasEmployment))
	return hasEmployment, nil
}

var educationXML = `<?xml version="1.0" encoding="UTF-8"?>
<education:education
	xmlns:common="http://www.orcid.org/ns/common" xmlns:education="http://www.orcid.org/ns/education"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.orcid.org/ns/education ../education-3.0.xsd ">
	<common:organization>
		<common:name>University of Virginia</common:name>
		<common:address>
			<common:city>Charlottesville</common:city>
			<common:region>VA</common:region>
			<common:country>US</common:country>
		</common:address>
		<common:disambiguated-organization>
				<common:disambiguated-organization-identifier>https://ror.org/0153tk833</common:disambiguated-organization-identifier>
				<common:disambiguation-source>ROR</common:disambiguation-source>
			</common:disambiguated-organization>
	</common:organization>
</education:education>`

// SendEmployment updates the user's ORCID with UVA Employent info
func SendEducation(attributes api.OrcidAttributes) {

	// First check for existing UVA Employment
	if hasEmp, err := hasExistingEducation(attributes); err != nil || hasEmp {
		return
	}

	url := fmt.Sprintf("%s/%s/education", config.Configuration.OrcidSecureURL, attributes.Orcid)
	// construct the auth field
	auth := fmt.Sprintf("Bearer %s", attributes.OauthAccessToken)
	start := time.Now()
	resp, _, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Post(url).
		Set("Accept", "application/xml").
		Set("Content-Type", "application/xml").
		Set("Authorization", auth).
		Send(educationXML).
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	// check for errors
	if errs != nil || resp.StatusCode != http.StatusCreated {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s %s in %s \n %s", url, resp.Status, errs, duration, resp.Body))
		return
	}
	logger.Log(fmt.Sprintf("INFO: UVA Education created for %s", attributes.Orcid))

	return
}

// hasExistingEmployment checks for existing UVA Employment matching our OrcidClientID
func hasExistingEducation(attributes api.OrcidAttributes) (bool, error) {
	hasUVAEducation := false
	// Only need to know employment here
	var educationStruct struct {
		AffiliationGroup []struct {
			Summaries []struct {
				Education struct {
					Source struct {
						SourceClientID struct {
							Path string `json:"path"`
						} `json:"source-client-id"`
					} `json:"source"`
				} `json:"education-summary"`
			} `json:"summaries"`
		} `json:"affiliation-group"`
	}
	url := fmt.Sprintf("%s/%s/educations", config.Configuration.OrcidSecureURL, attributes.Orcid)
	// construct the auth field
	auth := fmt.Sprintf("Bearer %s", attributes.OauthAccessToken)
	start := time.Now()
	resp, body, errs := gorequest.New().
		SetDebug(config.Configuration.Debug).
		Get(url).
		Set("Accept", "application/json").
		Set("Content-Type", "application/json").
		Set("Authorization", auth).
		Timeout(time.Duration(config.Configuration.ServiceTimeout) * time.Second).
		End()
	duration := time.Since(start)

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	// check for errors
	if errs != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		return false, errs[0]
	}

	if err := json.Unmarshal([]byte(body), &educationStruct); err != nil {
		logger.Log(fmt.Sprintf("ERROR: service (%s) returns %s in %s", url, errs, duration))
		return false, errs[0]
	}

	// Check the nested json
	for _, g := range educationStruct.AffiliationGroup {
		for _, s := range g.Summaries {
			// Check if existing employments match our client ID
			if s.Education.Source.SourceClientID.Path == config.Configuration.OrcidClientID {
				hasUVAEducation = true
			}
		}
	}
	//logger.Log(fmt.Sprintf("hasEducation: %s ", resp.Body))
	logger.Log(fmt.Sprintf("INFO: %s has UVA education: %t", attributes.Orcid, hasUVAEducation))
	return hasUVAEducation, nil
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
		Timeout(authTimeout).
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

// RenewAccessToken -- renew the access token
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
		Set("grant_type", "refresh_token").
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

// GetOrcidDetails -- get details for the specified ORCID
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
		return nil, resp.StatusCode, errors.New(orcid)
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
//	logger.Log(fmt.Sprintf("INFO: Service (%s) returns http %d in %s", url, resp.StatusCode, duration))
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

// GetPublicEndpointStatus -- get the public endpoint status
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

// GetSecureEndpointStatus -- get the secure endpoint status
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

// issue a GET to the specified URL and use a bearer token for aurthorization
// ignore the response payload and return an error if we get a non-200 response
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
		Timeout(authTimeout).
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

// singleton type access to the access token
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
