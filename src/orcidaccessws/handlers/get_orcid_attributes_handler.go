package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"orcidaccessws/authtoken"
	"orcidaccessws/config"
	"orcidaccessws/dao"
	"orcidaccessws/logger"
)

//
// GetOrcidAttributes -- get orcid attributes handler
//
func GetOrcidAttributes(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	token := r.URL.Query().Get("auth")

	// parameters OK?
	if isEmpty(id) || isEmpty(token) {
		status := http.StatusBadRequest
		encodeOrcidAttributesResponse(w, status, http.StatusText(status), nil)
		return
	}

	// validate the token
	if authtoken.Validate(config.Configuration.AuthTokenEndpoint, "getorcid", token, config.Configuration.ServiceTimeout) == false {
		status := http.StatusForbidden
		encodeOrcidAttributesResponse(w, status, http.StatusText(status), nil)
		return
	}

	// get the ORCID details
	attributes, err := dao.DB.GetOrcidAttributesByCid(id)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: %s", err.Error()))
		status := http.StatusInternalServerError
		encodeOrcidAttributesResponse(w, status,
			fmt.Sprintf("%s (%s)", http.StatusText(status), err),
			nil)
		return
	}

	// we did not find the item, return 404
	if attributes == nil || len(attributes) == 0 {
		status := http.StatusNotFound
		encodeOrcidAttributesResponse(w, status, http.StatusText(status), nil)
		return
	}

	status := http.StatusOK
	encodeOrcidAttributesResponse(w, status, http.StatusText(status), attributes)
}

//
// end of file
//
