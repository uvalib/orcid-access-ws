package api

type OrcidResponse struct {
	Status      int                `json:"status"`
	Message     string             `json:"message"`
	Attributes  []*OrcidAttributes `json:"results"`
}
