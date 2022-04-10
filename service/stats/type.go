package stats

type (
	ListeningStats struct {
		TotalListeningCount int64   `json:"total_listening_count"`
		Rooms               []int64 `json:"rooms"`
	}
)
