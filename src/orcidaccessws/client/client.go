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

var debugHTTP = false
var serviceTimeout = 30

//
// HealthCheck -- calls the service health check method
//
func HealthCheck(endpoint string) int {

   url := fmt.Sprintf("%s/healthcheck", endpoint)
   //fmt.Printf( "%s\n", url )

   resp, _, errs := gorequest.New().
      SetDebug(debugHTTP).
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

//
// VersionCheck -- calls the service version check method
//
func VersionCheck(endpoint string) (int, string) {

   url := fmt.Sprintf("%s/version", endpoint)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
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

//
// RuntimeCheck -- calls the service runtime method
//
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

//
// Statistics -- call the statistics method on the service
//
func Statistics(endpoint string) (int, *api.Statistics) {

   url := fmt.Sprintf("%s/statistics", endpoint)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
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

//
// GetOrcidAttributes -- call the get orcid attributes method on the service
//
func GetOrcidAttributes(endpoint string, cid string, token string) (int, []*api.OrcidAttributes) {

   url := fmt.Sprintf("%s/cid/%s?auth=%s", endpoint, cid, token)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
      Get(url).
      Timeout(time.Duration(serviceTimeout) * time.Second).
      End()

   if errs != nil {
      return http.StatusInternalServerError, nil
   }

   defer io.Copy(ioutil.Discard, resp.Body)
   defer resp.Body.Close()

   r := api.OrcidAttributesResponse{}
   err := json.Unmarshal([]byte(body), &r)
   if err != nil {
      return http.StatusInternalServerError, nil
   }

   return resp.StatusCode, r.Attributes
}

//
// GetAllOrcidAttributes -- call get all orcid attributes on the service
//
func GetAllOrcidAttributes(endpoint string, token string) (int, []*api.OrcidAttributes) {

   url := fmt.Sprintf("%s/cid?auth=%s", endpoint, token)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
      Get(url).
      Timeout(time.Duration(serviceTimeout) * time.Second).
      End()

   if errs != nil {
      return http.StatusInternalServerError, nil
   }

   defer io.Copy(ioutil.Discard, resp.Body)
   defer resp.Body.Close()

   r := api.OrcidAttributesResponse{}
   err := json.Unmarshal([]byte(body), &r)
   if err != nil {
      return http.StatusInternalServerError, nil
   }

   return resp.StatusCode, r.Attributes
}

//
// SetOrcidAttributes -- call set orcid attributes on the service
//
func SetOrcidAttributes(endpoint string, cid string, token string, attributes api.OrcidAttributes) int {

   url := fmt.Sprintf("%s/cid/%s?auth=%s", endpoint, cid, token)
   //fmt.Printf( "%s\n", url )

   resp, _, errs := gorequest.New().
      SetDebug(debugHTTP).
      Put(url).
      Send(attributes).
      Timeout(time.Duration(serviceTimeout)*time.Second).
      Set("Content-Type", "application/json").
      End()

   if errs != nil {
      return http.StatusInternalServerError
   }

   defer io.Copy(ioutil.Discard, resp.Body)
   defer resp.Body.Close()

   return resp.StatusCode
}

//
// DelOrcidAttributes -- call delete orcid attributes on the service
//
func DelOrcidAttributes(endpoint string, cid string, token string) int {

   url := fmt.Sprintf("%s/cid/%s?auth=%s", endpoint, cid, token)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
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

//
// UpdateActivity -- call update activity on the service
//
func UpdateActivity(endpoint string, cid string, token string, activity api.ActivityUpdate) (int, string) {

   url := fmt.Sprintf("%s/cid/%s/activity?auth=%s", endpoint, cid, token)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
      Put(url).
      Send(activity).
      Timeout(time.Duration(serviceTimeout)*time.Second).
      Set("Content-Type", "application/json").
      End()

   if errs != nil {
      return http.StatusInternalServerError, ""
   }

   defer io.Copy(ioutil.Discard, resp.Body)
   defer resp.Body.Close()

   r := api.UpdateActivityResponse{}
   err := json.Unmarshal([]byte(body), &r)
   if err != nil {
      return http.StatusInternalServerError, ""
   }

   return resp.StatusCode, r.UpdateCode
}

//
// GetOrcidDetails -- call get orcid details on the service
//
func GetOrcidDetails(endpoint string, orcid string, token string) (int, *api.OrcidDetails) {

   url := fmt.Sprintf("%s/orcid/%s?auth=%s", endpoint, orcid, token)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
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

//
// SearchOrcid -- call search orcid on the service
//
func SearchOrcid(endpoint string, search string, start string, max string, token string) (int, []*api.OrcidDetails, int) {

   url := fmt.Sprintf("%s/orcid?q=%s&start=%s&max=%s&auth=%s", endpoint, search, start, max, token)
   //fmt.Printf( "%s\n", url )

   resp, body, errs := gorequest.New().
      SetDebug(debugHTTP).
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

//
// end of file
//
