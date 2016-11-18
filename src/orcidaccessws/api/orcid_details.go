package api

type OrcidDetails struct {
    Relevancy             string   `json:"relevancy,omitempty"`
    Id                    string   `json:"id,omitempty"`
    Uri                   string   `json:"uri,omitempty"`
    DisplayName           string   `json:"display_name,omitempty"`
    FirstName             string   `json:"first_name,omitempty"`
    LastName              string   `json:"last_name,omitempty"`
    Biography             string   `json:"biography,omitempty"`
    Keywords           [] string   `json:"keywords,omitempty"`
    ResearchUrls       [] string   `json:"urls,omitempty"`
}