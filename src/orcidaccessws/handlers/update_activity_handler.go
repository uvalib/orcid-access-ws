package handlers

import (
   "fmt"
   "github.com/gorilla/mux"
   "net/http"
   "orcidaccessws/authtoken"
   "orcidaccessws/config"
   "orcidaccessws/dao"
   "orcidaccessws/logger"
   "encoding/json"
   "orcidaccessws/api"
   "io"
   "io/ioutil"
   "errors"
   "orcidaccessws/orcid"
)

var emptyUpdateCode = ""

func UpdateActivity(w http.ResponseWriter, r *http.Request) {

   vars := mux.Vars(r)
   id := vars["id"]
   token := r.URL.Query().Get("auth")

   // update the statistics
   Statistics.RequestCount++
   Statistics.UpdateActivityCount++

   // parameters OK?
   if isEmpty(id) || isEmpty(token) {
      status := http.StatusBadRequest
      encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode )
      return
   }

   // validate the token
   if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.Timeout) == false {
      status := http.StatusForbidden
      encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode )
      return
   }

   decoder := json.NewDecoder(r.Body)
   activity := api.ActivityUpdate{}

   if err := decoder.Decode(&activity); err != nil {
      logger.Log(fmt.Sprintf("ERROR: decoding request payload: %s", err))
      status := http.StatusBadRequest
      encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode )
      return
   }

   defer io.Copy(ioutil.Discard, r.Body)
   defer r.Body.Close()

   if err := validateRequestPayload( activity ); err != nil {
      logger.Log(fmt.Sprintf("ERROR: invalid request payload: %s", err))
      status := http.StatusBadRequest
      encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode )
      return
   }

   // get the users ORCID attributes
   attributes, err := dao.Database.GetOrcidAttributesByCid(id)
   if err != nil {
      logger.Log(fmt.Sprintf("ERROR: %s", err.Error()))
      status := http.StatusInternalServerError
      encodeUpdateActivityResponse(w, status,
         fmt.Sprintf("%s (%s)", http.StatusText(status), err),
         "")
      return
   }

   // we did not find the item, return 404
   if attributes == nil || len(attributes) == 0 {
      status := http.StatusNotFound
      encodeUpdateActivityResponse(w, status, http.StatusText(status), emptyUpdateCode )
      return
   }

   // update the activity
   updateCode, status, err := orcid.UpdateOrcidActivity( attributes[0].Orcid, attributes[0].OauthAccessToken, activity )

   // TODO handle the access token update and retry
   //if stuff {
   //   renew the access token...
   //   retry the update activity if we can
   //}

   // we did got an error, return it
   if status != http.StatusOK {
      encodeUpdateActivityResponse(w, status,
         fmt.Sprintf("%s (%s)", http.StatusText(status), err), emptyUpdateCode )
      return
   }

   encodeUpdateActivityResponse(w, status, http.StatusText(status), updateCode )
}

func validateRequestPayload( activity api.ActivityUpdate ) error {

   //
   // basic validation that the required fields exist
   //

   if len( activity.Work.Title ) == 0 {
      return errors.New( "Empty work title" )
   }

   if len( activity.Work.ResourceType ) == 0 {
      return errors.New( "Empty work resource type" )
   }

   if len( activity.Work.Url ) == 0 {
      return errors.New( "Empty work url" )
   }

   return nil
}