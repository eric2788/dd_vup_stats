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
		Behaviours       []*Behaviour `gorm:"foreignKey:Uid;references:Uid"`
		TargetBehaviours []*Behaviour `gorm:"foreignKey:TargetUid;references:Uid"`

		LastListen *LastListen `gorm:"foreignKey:Uid;references:Uid"`
	}

	Behaviour struct {
		ID        uint `gorm:"primaryKey;autoIncrement"`
		Uid       int64
		CreatedAt time.Time
		TargetUid int64
		Command   string
		Display   string
		Image     sql.NullString
	}

	LastListen struct {
		Uid          int64     `gorm:"primaryKey"`
		LastListenAt time.Time `gorm:"default:null"`
	}
)
