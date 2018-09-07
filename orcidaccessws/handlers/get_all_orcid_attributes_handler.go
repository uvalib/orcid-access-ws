package handlers

import (
	"fmt"
	"net/http"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/authtoken"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/config"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/dao"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
)

//
// GetAllOrcidAttributes - get all orcid attributes handler
//
func GetAllOrcidAttributes(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("auth")

	// parameters OK?
	if isEmpty(token) {
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
	attributes, err := dao.DB.GetAllOrcidAttributes()
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