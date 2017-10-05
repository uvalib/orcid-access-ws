package handlers

import (
	"net/http"
	"orcidaccessws/dao"
	"orcidaccessws/orcid"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	// update the statistics
	Statistics.RequestCount++
	Statistics.HeartbeatCount++

	status := http.StatusOK
	dbErr := dao.Database.Check()
	orcidPublicErr := orcid.GetPublicEndpointStatus( )
	orcidSecureErr := orcid.GetSecureEndpointStatus( )

	var db_msg, orcid_public_msg, orcid_secure_msg string

	if dbErr != nil || orcidPublicErr != nil || orcidSecureErr != nil {

		status = http.StatusInternalServerError

		if dbErr != nil {
			db_msg = dbErr.Error()
		}

		if orcidPublicErr != nil {
			orcid_public_msg = orcidPublicErr.Error()
		}

		if orcidSecureErr != nil {
			orcid_secure_msg = orcidSecureErr.Error()
		}
	}

	encodeHealthCheckResponse(w, status, db_msg, orcid_public_msg, orcid_secure_msg )
}
