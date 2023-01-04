package statistics

type (
	ListeningStats struct {
		TotalListeningCount int64   `json:"total_listening_count"`
		Rooms               []int64 `json:"rooms"`
	}

	VtbsMoeResp struct {
		Mid    int64  `json:"mid"`
		UName  string `json:"uname"`
		RoomId int64  `json:"roomid"`
	}

	// Common
	VupData struct {
		Name   string
		RoomId int64
	}

	VupFetcher func() (map[int64]VupData, error)

	VupJsonData struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		RoomId int64  `json:"room_id"`
		Face   string `json:"face"`
	}

	UserResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Mid      int64  `json:"mid"`
			Name     string `json:"name"`
			Official struct {
				// 0 普通人
				// 1 知名up
				// 2 高能主播
				// 3 B站机构账户
				// 4 政府相关账户 (大概)
				// 5 企业账户 (大概)
				Role  int    `json:"role"`
				Title string `json:"title"`
				Type  int    `json:"type"`
			}
		} `json:"data"`
	}
)
