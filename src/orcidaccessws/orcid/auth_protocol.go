package orcid

//
// oauth request/response structures
//

type oauthRequest struct {
   ClientID     string `json:"client_id,omitempty"`
   ClientSecret string `json:"client_secret,omitempty"`
   Scope        string `json:"scope,omitempty"`
   GrantType    string `json:"grant_type,omitempty"`
}

type oauthResponse struct {
   AccessToken      string `json:"access_token,omitempty"`
   RefreshToken     string `json:"refresh_token,omitempty"`
   TokenType        string `json:"token_type,omitempty"`
   Scope            string `json:"scope,omitempty"`
}

//
// end of file
//
