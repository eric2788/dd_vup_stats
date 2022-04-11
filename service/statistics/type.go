package statistics

type (
	ListeningStats struct {
		TotalListeningCount int64   `json:"total_listening_count"`
		Rooms               []int64 `json:"rooms"`
	}

	VtbsMoeResp struct {
		Mid  int64  `json:"mid"`
		Uuid string `json:"uuid"`
	}
)
