package api

type UpdateActivityResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	UpdateCode string `json:"update_code"`
}
