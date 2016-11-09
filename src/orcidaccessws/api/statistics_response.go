package api

type StatisticsResponse struct {
    StandardResponse
    Details       Statistics `json:"statistics"`
}


