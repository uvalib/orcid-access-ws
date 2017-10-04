package api

type HealthCheckResponse struct {
	DbCheck          HealthCheckResult `json:"mysql"`
	OrcidPublicCheck HealthCheckResult `json:"orcid_public"`
	OrcidSecureCheck HealthCheckResult `json:"orcid_member"`
}
