package api

type Statistics struct {
	RequestCount            int `json:"request_count"`
	SetOrcidAttribsCount    int `json:"set_attribs_count"`
	GetOrcidAttribsCount    int `json:"get_attribs_count"`
	DelOrcidAttribsCount    int `json:"del_attribs_count"`
	UpdateActivityCount     int `json:"update_activity_count"`
	GetOrcidDetailsCount    int `json:"get_details_count"`
	SearchOrcidDetailsCount int `json:"search_count"`
	HeartbeatCount          int `json:"heartbeat_count"`
}
