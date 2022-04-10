package vup

import (
	"time"
)

type (
	ListResp struct {
		Page    int         `json:"page"`
		Size    int         `json:"size"`
		MaxPage int64       `json:"max_page"`
		List    []*UserResp `json:"list"`
	}

	UserResp struct {
		UserInfo
		Listening       bool      `json:"listening"`
		LastListenedAt  time.Time `json:"last_listened_at"`
		LastBehaviourAt time.Time `json:"last_behaviour_at"`
	}

	UserInfo struct {
		Uid           int64     `json:"uid"`
		Name          string    `json:"name"`
		Face          string    `json:"face"`
		FirstListenAt time.Time `json:"first_listen_at"`
		RoomId        int64     `json:"room_id"`
		Sign          string    `json:"sign"`
	}
)
