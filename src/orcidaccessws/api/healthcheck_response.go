package api

type HealthCheckResponse struct {
   DbCheck     HealthCheckResult `json:"mysql"`
   OrcidCheck  HealthCheckResult `json:"orcid"`
}

