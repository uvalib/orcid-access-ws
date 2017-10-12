package handlers

import (
   "net/http"
   "orcidaccessws/dao"
   "orcidaccessws/orcid"
)

//
// HealthCheck -- do the healthcheck
//
func HealthCheck(w http.ResponseWriter, r *http.Request) {

   // update the statistics
   Statistics.RequestCount++
   Statistics.HeartbeatCount++

   status := http.StatusOK
   dbErr := dao.DB.CheckDB()
   orcidPublicErr := orcid.GetPublicEndpointStatus( )
   orcidSecureErr := orcid.GetSecureEndpointStatus( )

   var dbMessage, orcidPublicMessage, orcidSecureMessage string

   if dbErr != nil || orcidPublicErr != nil || orcidSecureErr != nil {

      status = http.StatusInternalServerError

      if dbErr != nil {
         dbMessage = dbErr.Error()
      }

      if orcidPublicErr != nil {
         orcidPublicMessage = orcidPublicErr.Error()
      }

      if orcidSecureErr != nil {
         orcidSecureMessage = orcidSecureErr.Error()
      }
   }

   encodeHealthCheckResponse(w, status, dbMessage, orcidPublicMessage, orcidSecureMessage)
}

//
// end of file
//
