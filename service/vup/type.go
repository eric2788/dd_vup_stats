package vup

import (
	"time"
	"vup_dd_stats/service/db"
)

type (
	ListResp[K any] struct {
		Page    int   `json:"page"`
		Size    int   `json:"size"`
		MaxPage int64 `json:"max_page"`
		Total   int64 `json:"total"`
		List    []K   `json:"list"`
	}

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
		BehavioursCount map[string]int64 `json:"behaviours_count"`
	}

	UserInfo struct {
		SimpleUserInfo
		FirstListenAt   time.Time  `json:"first_listen_at"`
		LastBehaviourAt *time.Time `json:"last_behaviour_at"`
		DDCount         int64      `json:"dd_count"`
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
		Count int64 `json:"count"`
	}

	Analysis struct {
		TopDDVups    []AnalysisUserInfo `json:"top_dd_vups"`
		TopGuestVups []AnalysisUserInfo `json:"top_guest_vups"`
	}

	GlobalStatistics struct {
		TotalVupRecorded           int64                         `json:"total_vup_recorded"`
		CurrentListeningCount      int64                         `json:"current_listening_count"`
		MostDDBehaviourVupCommands map[string][]AnalysisUserInfo `json:"most_dd_behaviour_vup_commands"`
		MostDDBehaviourVups        []AnalysisUserInfo            `json:"most_dd_behaviour_vups"`
		MostDDVups                 []AnalysisUserInfo            `json:"most_dd_vups"` // D 最多直播間的人
		TotalDDBehaviours          int64                         `json:"total_dd_behaviours"`
	}
)
