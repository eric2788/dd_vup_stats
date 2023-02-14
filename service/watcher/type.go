package watcher

import (
	"time"
	"vup_dd_stats/service/stats"
)

type (
	WatcherInfo struct {
		Uid   int64  `json:"uid"`
		UName string `json:"u_name"`
	}

	VupInfo struct {
		Uid  int64  `json:"uid"`
		Name string `json:"name"`
		Face string `json:"face"`
	}

	AnalysisWatcherInfo struct {
		WatcherInfo
		Count int64   `json:"count"`
		Price float64 `json:"price"`
	}

	AnalysisVupInfo struct {
		VupInfo
		Count int64   `json:"count"`
		Price float64 `json:"price"`
	}

	PricedWatcherInfo struct {
		WatcherInfo
		Spent float64 `json:"spent"`
	}

	PricedVupInfo struct {
		VupInfo
		Spent float64 `json:"spent"`
	}

	Analysis struct {
		TopDDVups    []AnalysisVupInfo `json:"top_dd_vups"`
		TopSpentVups []PricedVupInfo   `json:"top_spent_vups"`
	}

	WatcherResp struct {
		WatcherInfo
		UNames          string             `json:"u_names"`
		FirstListenAt   time.Time          `json:"first_listen_at"`
		LastBehaviourAt time.Time          `json:"last_behaviour_at"`
		DDCount         int64              `json:"dd_count"`
		TotalSpent      float64            `json:"total_spent"`
		Behaviours      []stats.TotalStats `json:"behaviours" gorm:"-"`
	}
)
