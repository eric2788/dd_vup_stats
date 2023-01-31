package watcher

import (
	"math"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"

	"gorm.io/gorm"
)

func GetRecords(uid int64, page, pageSize int, cmd string) (*stats.ListResp[db.WatcherBehaviour], error) {

	// ensure page is valid
	page = int(math.Max(1, float64(page)))

	//ensure pageSize is valid
	pageSize = int(math.Max(1, float64(pageSize)))

	var behaviours []db.WatcherBehaviour

	r := db.Database.Model(&db.WatcherBehaviour{}).Where("uid = ?", uid)

	if cmd != "" {
		r = r.Where("command = ?", cmd)
	}

	r = r.Session(&gorm.Session{})

	err := r.
		Order("created_at desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&behaviours).
		Error

	if err != nil {
		return nil, err
	}

	var totalSearchCount int64

	err = r.
		Count(&totalSearchCount).
		Error

	if err != nil {
		return nil, err
	}

	return &stats.ListResp[db.WatcherBehaviour]{
		Total:   totalSearchCount,
		List:    behaviours,
		Page:    page,
		Size:    pageSize,
		MaxPage: int64(math.Ceil(float64(totalSearchCount) / float64(pageSize))),
	}, nil
}
