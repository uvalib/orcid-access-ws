package orcid

//
// v2.0 response structure
//

type activityUpdateResponse struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ResponseCode     int    `json:"response-code,omitempty"`
	UserMessage      string `json:"user-message,omitempty"`
	DeveloperMessage string `json:"developer-message,omitempty"`
}

//
// end of file
//
