package api

type OrcidDetailsResponse struct {
   Status        int           `json:"status"`
   Message       string        `json:"message"`
   Details  [] * OrcidDetails  `json:"orcid"`
}

