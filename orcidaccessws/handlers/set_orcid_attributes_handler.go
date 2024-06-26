package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/api"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/authtoken"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/dao"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/orcid"
)

// SetOrcidAttributes -- set the orcid attributes handler
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
	if authtoken.Validate(config.Configuration.SharedSecret, token) == false {
		status := http.StatusForbidden
		encodeStandardResponse(w, status, http.StatusText(status))
		return
	}

	decoder := json.NewDecoder(r.Body)
	attributes := api.OrcidAttributes{}

	if err := decoder.Decode(&attributes); err != nil {
		logger.Log(fmt.Sprintf("ERROR: decoding set attributes request payload %s", err))
		status := http.StatusBadRequest
		encodeStandardResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err.Error()))
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
	err := dao.Store.SetOrcidAttributesByCid(id, attributes)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: %s", err.Error()))
		status := http.StatusInternalServerError
		encodeStandardResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err.Error()))
		return
	}

	// Send Employment or Education to ORCID
	if !isEmpty(attributes.UserTypes) {
		sendEmployment := false
		sendEducation := false
		types := strings.Split(attributes.UserTypes, ";")
		employmentCheck := regexp.MustCompile(`Staff|Employee|Faculty`)

		for _, userType := range types {
			matched := employmentCheck.MatchString(userType)
			if matched {
				sendEmployment = true
			} else {
				sendEducation = true
			}
		}
		if sendEducation {
			orcid.SendEducation(attributes)
		}
		if sendEmployment {
			orcid.SendEmployment(attributes)
		}
	}

	status := http.StatusOK
	encodeStandardResponse(w, status, http.StatusText(status))
}

//
// end of file
//
