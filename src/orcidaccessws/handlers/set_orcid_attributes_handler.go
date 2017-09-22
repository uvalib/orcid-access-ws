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
)

func SetOrcidAttributes(w http.ResponseWriter, r *http.Request) {

   vars := mux.Vars(r)
   id := vars["id"]
   token := r.URL.Query().Get("auth")

   // update the statistics
   Statistics.RequestCount++
   Statistics.SetOrcidCount++

   // parameters OK?
   if isEmpty(id) || isEmpty(token) {
      status := http.StatusBadRequest
      encodeStandardResponse(w, status, http.StatusText(status))
      return
   }

   // validate the token
   if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "setorcid", token, config.Configuration.Timeout) == false {
      status := http.StatusForbidden
      encodeStandardResponse(w, status, http.StatusText(status))
      return
   }

   decoder := json.NewDecoder(r.Body)
   attributes := api.OrcidAttributes{}

   if err := decoder.Decode(&attributes); err != nil {
      logger.Log(fmt.Sprintf("ERROR: decoding request payload %s", err))
      status := http.StatusBadRequest
      encodeStandardResponse(w, status, http.StatusText(status))
      return
   }

   defer io.Copy(ioutil.Discard, r.Body)
   defer r.Body.Close()

   // at minimum, the ORCID must be defined
   if isEmpty(attributes.Orcid) {
      status := http.StatusBadRequest
      encodeStandardResponse(w, status, http.StatusText(status))
      return
   }

   // set the ORCID attributes
   err := dao.Database.SetOrcidAttributesByCid(id, attributes)
   if err != nil {
      logger.Log(fmt.Sprintf("ERROR: %s", err.Error()))
      status := http.StatusInternalServerError
      encodeStandardResponse(w, status,
         fmt.Sprintf("%s (%s)", http.StatusText(status), err))
      return
   }

   status := http.StatusOK
   encodeStandardResponse(w, status, http.StatusText(status))
}