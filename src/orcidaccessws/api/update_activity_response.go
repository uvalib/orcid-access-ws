package api

//
// UpdateActivityResponse -- response for the update activity request
//
type UpdateActivityResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	UpdateCode string `json:"update_code"`
}

//
// end of file
//
