package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"orcidaccessws/api"
	"orcidaccessws/authtoken"
	"orcidaccessws/config"
	"orcidaccessws/dao"
	"orcidaccessws/logger"
)

//
// SetOrcidAttributes -- set the orcid attributes handler
//
func SetOrcidAttributes(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	token := r.URL.Query().Get("auth")

	// parameters OK?
	if isEmpty(id) || isEmpty(token) {
		status := http.StatusBadRequest
		encodeStandardResponse(w, status, http.StatusText(status))
		return
	}

	// validate the token
	if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "setorcid", token, config.Configuration.ServiceTimeout) == false {
		status := http.StatusForbidden
		encodeStandardResponse(w, status, http.StatusText(status))
		return
	}

	decoder := json.NewDecoder(r.Body)
	attributes := api.OrcidAttributes{}

	if err := decoder.Decode(&attributes); err != nil {
		logger.Log(fmt.Sprintf("ERROR: decoding set attributes request payload %s", err))
		status := http.StatusBadRequest
		encodeStandardResponse(w, status, http.StatusText(status))
		return
	}

	defer io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()

	// at minimum, the ORCID must be defined
	//if isEmpty(attributes.Orcid) {
	//   status := http.StatusBadRequest
	//   encodeStandardResponse(w, status, http.StatusText(status))
	//   return
	//}

	// set the ORCID attributes
	err := dao.DB.SetOrcidAttributesByCid(id, attributes)
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

//
// end of file
//
