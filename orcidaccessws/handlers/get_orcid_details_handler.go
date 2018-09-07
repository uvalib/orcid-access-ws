package handlers

import (
	//"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/authtoken"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	//"github.com/uvalib/orcid-access-ws/orcidaccessws/orcid"
)

//
// GetOrcidDetails -- get the orcid details handler
//
func GetOrcidDetails(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	token := r.URL.Query().Get("auth")

	// parameters OK?
	if isEmpty(id) || isEmpty(token) {
		status := http.StatusBadRequest
		encodeOrcidDetailsResponse(w, status, http.StatusText(status), nil)
		return
	}

	// validate the token
	if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.ServiceTimeout) == false {
		status := http.StatusForbidden
		encodeOrcidDetailsResponse(w, status, http.StatusText(status), nil)
		return
	}

	//
	// not implemented as we have moved to the 2.0 API which supports different behavior
	//
	status := http.StatusNotImplemented
	encodeOrcidDetailsResponse(w, status, http.StatusText(status), nil)

	// get the ORCID details
	//orcid, status, err := orcid.GetOrcidDetails(id)

	// we did got an error, return it
	//if status != http.StatusOK {
	//   encodeOrcidDetailsResponse(w, status,
	//      fmt.Sprintf("%s (%s)", http.StatusText(status), err), nil)
	//   return
	//}

	//status = http.StatusOK
	//encodeOrcidDetailsResponse(w, status, http.StatusText(status), orcid)
}

//
// end of file
//
