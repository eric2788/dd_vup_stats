package db

import (
	"database/sql"
	"time"
)

type (
	Vup struct {
		Uid              int64  `gorm:"primaryKey;autoIncrement:false"`
		Name             string `gorm:"uniqueIndex"`
		Face             string `gorm:"uniqueIndex"`
		FirstListenAt    time.Time
		RoomId           int64        `gorm:"uniqueIndex"`
		Sign             string       `gorm:"uniqueIndex"`
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
)
