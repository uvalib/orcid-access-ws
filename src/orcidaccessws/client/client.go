package client

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io"
	"io/ioutil"
	"net/http"
	"orcidaccessws/api"
	"time"
)

var debugHttp = false
var serviceTimeout = 5

func HealthCheck(endpoint string) int {

	url := fmt.Sprintf("%s/healthcheck", endpoint)
	//fmt.Printf( "%s\n", url )

	resp, _, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return resp.StatusCode
}

func VersionCheck(endpoint string) (int, string) {

	url := fmt.Sprintf("%s/version", endpoint)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, ""
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.VersionResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, ""
	}

	return resp.StatusCode, r.Version
}

func RuntimeCheck(endpoint string) (int, *api.RuntimeResponse) {

	url := fmt.Sprintf("%s/runtime", endpoint)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(false).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, nil
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.RuntimeResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return resp.StatusCode, &r
}

func Statistics(endpoint string) (int, *api.Statistics) {

	url := fmt.Sprintf("%s/statistics", endpoint)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, nil
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.StatisticsResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return resp.StatusCode, &r.Details
}

func GetOrcid(endpoint string, id string, token string) (int, []*api.Orcid) {

	url := fmt.Sprintf("%s/cid/%s?auth=%s", endpoint, id, token)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, nil
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.OrcidResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return resp.StatusCode, r.Orcids
}

func DelOrcid(endpoint string, id string, token string) int {

	url := fmt.Sprintf("%s/cid/%s?auth=%s", endpoint, id, token)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Delete(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.StandardResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError
	}

	return resp.StatusCode
}

func GetAllOrcid(endpoint string, token string) (int, []*api.Orcid) {

	url := fmt.Sprintf("%s/cid?auth=%s", endpoint, token)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, nil
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.OrcidResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return resp.StatusCode, r.Orcids
}

func GetOrcidDetails(endpoint string, orcid string, token string) (int, *api.OrcidDetails) {

	url := fmt.Sprintf("%s/orcid/%s?auth=%s", endpoint, orcid, token)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, nil
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.OrcidDetailsResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return resp.StatusCode, r.Details
}

func SearchOrcid(endpoint string, search string, start string, max string, token string) (int, []*api.OrcidDetails, int) {

	url := fmt.Sprintf("%s/orcid?q=%s&start=%s&max=%s&auth=%s", endpoint, search, start, max, token)
	//fmt.Printf( "%s\n", url )

	resp, body, errs := gorequest.New().
		SetDebug(debugHttp).
		Get(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError, nil, 0
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	r := api.OrcidSearchResponse{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return http.StatusInternalServerError, nil, 0
	}

	return resp.StatusCode, r.Results, r.Total
}

func SetOrcid(endpoint string, cid string, orcid string, token string) int {

	url := fmt.Sprintf("%s/cid/%s/%s?auth=%s", endpoint, cid, orcid, token)
	//fmt.Printf( "%s\n", url )

	resp, _, errs := gorequest.New().
		SetDebug(debugHttp).
		Put(url).
		Timeout(time.Duration(serviceTimeout) * time.Second).
		End()

	if errs != nil {
		return http.StatusInternalServerError
	}

	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return resp.StatusCode
}
