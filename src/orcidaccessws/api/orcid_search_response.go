package api

//
// OrcidSearchResponse -- orcid search response
//
type OrcidSearchResponse struct {
   Status  int             `json:"status"`
   Message string          `json:"message"`
   Start   int             `json:"start"`
   Count   int             `json:"count"`
   Total   int             `json:"total"`
   Results []*OrcidDetails `json:"results"`
}

//
// end of file
//
