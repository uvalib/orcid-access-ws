package api

//
// HealthCheckResponse -- response to the health check query
//
type HealthCheckResponse struct {
   DbCheck          HealthCheckResult `json:"mysql"`
   OrcidPublicCheck HealthCheckResult `json:"orcid_public"`
   OrcidSecureCheck HealthCheckResult `json:"orcid_member"`
}

//
// end of file
//
