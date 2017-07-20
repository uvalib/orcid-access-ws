package api

type Orcid struct {
	Id        string `json:"id,omitempty"`
	Cid       string `json:"cid,omitempty"`
	Orcid     string `json:"orcid,omitempty"`
	Uri       string `json:"uri,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
