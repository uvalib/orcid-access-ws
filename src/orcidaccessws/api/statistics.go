package api

type Statistics struct {
	RequestCount            int `json:"request_count"`
	SetOrcidCount           int `json:"set_orcid_count"`
	GetOrcidCount           int `json:"get_orcid_count"`
	DelOrcidCount           int `json:"del_orcid_count"`
	GetOrcidDetailsCount    int `json:"get_details_count"`
	SearchOrcidDetailsCount int `json:"search_count"`
	HeartbeatCount          int `json:"heartbeat_count"`
}
