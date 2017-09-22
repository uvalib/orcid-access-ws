package api

type OrcidAttributes struct {
	Id                string `json:"id,omitempty"`
	Cid               string `json:"cid,omitempty"`
	Orcid             string `json:"orcid,omitempty"`
	Uri               string `json:"uri,omitempty"`
	OauthAccessToken  string `json:"oauth_access_token,omitempty"`
	OauthRefreshToken string `json:"oauth_refresh_token,omitempty"`
	OauthScope        string `json:"scope,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}