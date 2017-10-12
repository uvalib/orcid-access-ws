package api

//
// OrcidDetailsResponse -- response to the orcid details request
//
type OrcidDetailsResponse struct {
   Status  int           `json:"status"`
   Message string        `json:"message"`
   Details *OrcidDetails `json:"orcid"`
}

//
// end of file
//
