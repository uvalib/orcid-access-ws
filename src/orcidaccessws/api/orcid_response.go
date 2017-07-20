package api

type OrcidResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Orcids  []*Orcid `json:"orcids"`
}
