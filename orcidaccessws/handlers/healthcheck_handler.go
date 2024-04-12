package handlers

import (
	"fmt"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/dao"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/logger"
	"github.com/uvalib/orcid-access-ws/orcidaccessws/orcid"
	"net/http"
)

// HealthCheck -- do the healthcheck
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	status := http.StatusOK
	dbErr := dao.Store.Check()
	orcidPublicErr := orcid.GetPublicEndpointStatus()
	orcidSecureErr := orcid.GetSecureEndpointStatus()

	var dbMessage, orcidPublicMessage, orcidSecureMessage string

	if dbErr != nil || orcidPublicErr != nil || orcidSecureErr != nil {

		if dbErr != nil {
			// only a database connection problem is considered an error (cos we can actually do something
			// about it)...
			status = http.StatusInternalServerError

			dbMessage = dbErr.Error()
			logger.Log(fmt.Sprintf("ERROR: Datastore reports '%s'", dbMessage))
		}

		if orcidPublicErr != nil {
			orcidPublicMessage = orcidPublicErr.Error()
			logger.Log(fmt.Sprintf("ERROR: ORCID public endpoint reports '%s'", orcidPublicMessage))
		}

		if orcidSecureErr != nil {
			orcidSecureMessage = orcidSecureErr.Error()
			logger.Log(fmt.Sprintf("ERROR: ORCID secure endpoint reports '%s'", orcidSecureMessage))
		}
	}

	encodeHealthCheckResponse(w, status, dbMessage, orcidPublicMessage, orcidSecureMessage)
}

//
// end of file
//
