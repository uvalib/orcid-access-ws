package orcid

//
// v2.0 response structure
//

const SECURE_PROTOCOL_VERSION = "2.0"

type activityUpdateResponse struct {
   Error            string  `json:"error,omitempty"`
   ErrorDescription string  `json:"error_description,omitempty"`
}

//
// end of file
//