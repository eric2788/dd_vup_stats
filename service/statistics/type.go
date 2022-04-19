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

	UserResp struct {
		Code int `json:"code"`
		Data struct {
			Mid      int64  `json:"mid"`
			Name     string `json:"name"`
			Official struct {
				Role  int    `json:"role"`
				Title string `json:"title"`
				Type  int    `json:"type"`
			}
		} `json:"data"`
	}
)
