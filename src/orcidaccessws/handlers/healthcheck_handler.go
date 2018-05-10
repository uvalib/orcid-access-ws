package handlers

import (
	"net/http"
	"orcidaccessws/dao"
	"orcidaccessws/orcid"
	"orcidaccessws/logger"
	"fmt"
)

//
// HealthCheck -- do the healthcheck
//
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	status := http.StatusOK
	dbErr := dao.DB.CheckDB()
	orcidPublicErr := orcid.GetPublicEndpointStatus()
	orcidSecureErr := orcid.GetSecureEndpointStatus()

	var dbMessage, orcidPublicMessage, orcidSecureMessage string

	if dbErr != nil || orcidPublicErr != nil || orcidSecureErr != nil {

		status = http.StatusInternalServerError

		if dbErr != nil {
			dbMessage = dbErr.Error()
			logger.Log(fmt.Sprintf( "ERROR: Database reports '%s'", dbMessage ) )
		}

		if orcidPublicErr != nil {
			orcidPublicMessage = orcidPublicErr.Error()
			logger.Log(fmt.Sprintf( "ERROR: ORCID public endpoint reports '%s'", orcidPublicMessage ) )
		}

		if orcidSecureErr != nil {
			orcidSecureMessage = orcidSecureErr.Error()
			logger.Log(fmt.Sprintf( "ERROR: ORCID secure endpoint reports '%s'", orcidSecureMessage ) )
		}
	}

	encodeHealthCheckResponse(w, status, dbMessage, orcidPublicMessage, orcidSecureMessage)
}

//
// end of file
//
