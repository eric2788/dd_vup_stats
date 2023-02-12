package vup

import (
	"time"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"
)

type (
	RecordResp struct {
		db.Behaviour
		VupFace string `json:"vup_face"`
	}

	UserResp struct {
		UserInfo
		Listening bool `json:"listening"`
	}

	UserDetailResp struct {
		UserResp
		BehavioursCount map[string]stats.TotalStats `json:"behaviours_count"`
	}

	UserInfo struct {
		SimpleUserInfo
		FirstListenAt   time.Time  `json:"first_listen_at"`
		LastBehaviourAt *time.Time `json:"last_behaviour_at"`
		DDCount         int64      `json:"dd_count"`
		TotalSpent      float64    `json:"total_spent"`
		LastListenedAt  time.Time  `json:"last_listened_at"`
	}

	SimpleUserInfo struct {
		Uid    int64  `json:"uid"`
		Name   string `json:"name"`
		Face   string `json:"face"`
		RoomId int64  `json:"room_id"`
		Sign   string `json:"sign"`
	}

	AnalysisUserInfo struct {
		SimpleUserInfo
		Count int64   `json:"count"`
		Price float64 `json:"price"`
	}

	PricedUserInfo struct {
		SimpleUserInfo
		Spent float64 `json:"spent"`
	}

	Analysis struct {
		TopDDVups    []AnalysisUserInfo `json:"top_dd_vups"`
		TopGuestVups []AnalysisUserInfo `json:"top_guest_vups"`
		TopSpentVups []PricedUserInfo   `json:"top_spent_vups,omitempty"`
	}

	CountStats struct {
		TotalVupRecorded      int64 `json:"total_vup_recorded"`
		CurrentListeningCount int64 `json:"current_listening_count"`
		TotalDDBehaviours     int64 `json:"total_dd_behaviours"`
	}

	Globalstats struct {
		TotalVupRecorded      int64              `json:"total_vup_recorded"`
		CurrentListeningCount int64              `json:"current_listening_count"`
		TotalDDBehaviours     int64              `json:"total_dd_behaviours"`
		MostDDBehaviourVups   []AnalysisUserInfo `json:"most_dd_behaviour_vups"`
		MostDDVups            []AnalysisUserInfo `json:"most_dd_vups"` // D 最多直播間的人
		MostSpentVups         []PricedUserInfo   `json:"most_spent_vups"`
	}
)
