package orcid

//
// v1.2 response structure
//

const publicProtocolVersion = "1.2"

// all responses have these attributes
type orcidCommonResponse struct {
	Version string           `json:"message-version,omitempty"`
	Error   stringValueField `json:"error-desc,omitempty"`
}

// a profile query contains a profile response
type orcidProfileResponse struct {
	Profile orcidProfile `json:"orcid-profile,omitempty"`
}

// a search contains a profile response
type orcidSearchResponse struct {
	SearchResults orcidResults `json:"orcid-search-results,omitempty"`
}

//
// structures within the ORCID protocol
//

type orcidProfile struct {
	ID  orcidID  `json:"orcid-identifier,omitempty"`
	Bio orcidBio `json:"orcid-bio,omitempty"`
}

type orcidID struct {
	URI string `json:"uri,omitempty"`
	ID  string `json:"path,omitempty"`
}

type orcidResults struct {
	Results    []orcidResult `json:"orcid-search-result,omitempty"`
	TotalFound int           `json:"num-found,omitempty"`
}

type orcidResult struct {
	Relevancy floatValueField `json:"relevancy-score,omitempty"`
	Profile   orcidProfile    `json:"orcid-profile,omitempty"`
}

type orcidBio struct {
	PersonalDetails orcidPersonalDetails `json:"personal-details,omitempty"`
	Biography       stringValueField     `json:"biography,omitempty"`
	Urls            orcidUrls            `json:"researcher-urls,omitempty"`
	Keywords        orcidKeywords        `json:"keywords,omitempty"`
}

type orcidPersonalDetails struct {
	GivenName  stringValueField `json:"given-names,omitempty"`
	FamilyName stringValueField `json:"family-name,omitempty"`
	CreditName stringValueField `json:"credit-name,omitempty"`
}

type orcidKeywords struct {
	Keywords []stringValueField `json:"keyword,omitempty"`
}

type orcidUrls struct {
	Urls []orcidURL `json:"researcher-url,omitempty"`
}

type orcidURL struct {
	Name stringValueField `json:"url-name,omitempty"`
	URL  stringValueField `json:"url,omitempty"`
}

type stringValueField struct {
	Value string `json:"value,omitempty"`
}

type floatValueField struct {
	Value float64 `json:"value,omitempty"`
}
