package api

// OrcidAttributesResponse -- response to the attributes request
type OrcidAttributesResponse struct {
	Status     int                `json:"status"`
	Message    string             `json:"message"`
	Attributes []*OrcidAttributes `json:"results"`
}

//
// end of file
//
