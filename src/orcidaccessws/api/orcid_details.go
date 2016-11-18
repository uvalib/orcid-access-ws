package api

type OrcidDetails struct {
    Orcid                 string   `json:"id,omitempty"`
    DisplayName           string   `json:"display_name,omitempty"`
    FirstName             string   `json:"first_name,omitempty"`
    LastName              string   `json:"last_name,omitempty"`
    Keywords           [] string   `json:"keywords,omitempty"`
    //ResearchUrls       [] string   `json:"xxx,omitempty"`
}