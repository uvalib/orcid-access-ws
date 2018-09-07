package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/authtoken"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/dao"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
)

//
// DeleteOrcidAttributes -- delete orcid attributes handler
//
func DeleteOrcidAttributes(w http.ResponseWriter, r *http.Request) {

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

	// get the ORCID details
	err := dao.DB.DelOrcidAttributesByCid(id)
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
