package api

type OrcidAttributesResponse struct {
	Status     int                `json:"status"`
	Message    string             `json:"message"`
	Attributes []*OrcidAttributes `json:"results"`
}
