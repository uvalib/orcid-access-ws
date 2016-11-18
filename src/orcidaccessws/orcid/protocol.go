package orcid

//import (
//)

//
// v1.2 response structure
//

const PROTOCOL_VERSION = "1.2"

// all responses have these attributes
type orcidCommonResponse struct {
    Version             string                 `json:"message-version,omitempty"`
    Error               valueField             `json:"error-desc,omitempty"`
}

// a profile query contains a profile response
type orcidProfileResponse struct {
    Profile             orcidProfile           `json:"orcid-profile,omitempty"`
}

// an ORCID user bio
type orcidProfile struct {
    Bio                 orcidBio               `json:"orcid-bio,omitempty"`
}

// a search contains ...
type orcidSearch struct {
}

//
// structures for ORCID user attributes
//

type orcidBio struct {
    PersonalDetails     orcidPersonalDetails   `json:"personal-details,omitempty"`
    Biography           valueVisibilityPair    `json:"biography,omitempty"`
}

type orcidPersonalDetails struct {
    GivenName           valueVisibilityPair    `json:"given-names,omitempty"`
    FamilyName          valueVisibilityPair    `json:"family-name,omitempty"`
    CreditName          valueVisibilityPair    `json:"credit-name,omitempty"`
}

type valueField struct {
    Value               string                 `json:"value,omitempty"`
}

type valueVisibilityPair struct {
    Value               string                 `json:"value,omitempty"`
    Visibility          string                 `json:"visibility,omitempty"`
}