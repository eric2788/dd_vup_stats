package db

import (
	"database/sql"
	"time"
)

type (
	Vup struct {
		Uid              int64 `gorm:"primaryKey;autoIncrement:false"`
		Name             string
		Face             string
		FirstListenAt    time.Time
		RoomId           int64
		Sign             string
		Behaviours       []*Behaviour `gorm:"foreignKey:Uid;references:Uid;OnDelete:CASCADE"`
		TargetBehaviours []*Behaviour `gorm:"foreignKey:TargetUid;references:Uid;OnDelete:CASCADE"`

		LastListen *LastListen `gorm:"foreignKey:Uid;references:Uid;OnDelete:CASCADE"`
	}

	Behaviour struct {
		ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
		Uid       int64          `json:"uid"`
		CreatedAt time.Time      `json:"created_at"`
		TargetUid int64          `json:"target_uid"`
		Command   string         `json:"command"`
		Display   string         `json:"display"`
		Image     sql.NullString `json:"image"`
		Price     float64        `json:"price" gorm:"default:0"`
	}

	LastListen struct {
		Uid          int64 `gorm:"primaryKey;autoIncrement:false"`
		LastListenAt time.Time
	}

	Analysis struct {
		AccessCount int64 `gorm:"default:0"`
		LastAccessDate string
	}

	UserAnalysis struct {
		Analysis
		Uid int64 `gorm:"primaryKey;autoIncrement:false"`
	}

	SearchAnalysis struct {
		Analysis
		SearchHash string `gorm:"primaryKey"`
		SearchText string
	}
)
