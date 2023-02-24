package watcher

import (
	"fmt"
	"vup_dd_stats/service/db"
)

func GetFanStatsForVup(uid int64, limit int, t string) ([]AnalysisWatcherInfo, error) {
	switch t {
	case "behaviours":
		return GetMostBehavioursByVup(uid, limit)
	case "spent":
		return GetMostSpentByVup(uid, limit)
	default:
		return nil, fmt.Errorf("不支持的类型: %q", t)
	}
}

// GetMostBehavioursByVup 返回该 vup 中最高互动的 watcher
func GetMostBehavioursByVup(uid int64, limit int) ([]AnalysisWatcherInfo, error) {

	var mostDDWatchers = make([]AnalysisWatcherInfo, 0)

	if limit == -1 {
		limit = 50000
	}

	err := db.Database.
		Raw(fmt.Sprintf(`
			select
				uid,
				u_name,
				sum(price) as price,
				sum(count) as count
			from watchers_fans
			where target_uid = ?
			group by uid, u_name
			order by count desc
			limit %d;
		`, limit), uid).
		Scan(&mostDDWatchers).
		Error

	if err != nil {
		return nil, err
	}

	return mostDDWatchers, nil
}

// GetMostSpentByVup 返回该 vup 中花费最多的 watcher
func GetMostSpentByVup(uid int64, limit int) ([]AnalysisWatcherInfo, error) {

	var mostSpentWatchers = make([]AnalysisWatcherInfo, 0)

	if limit == -1 {
		limit = 50000
	}

	err := db.Database.
		Raw(fmt.Sprintf(`
			select
				uid,
				u_name,
				sum(price) as price,
				sum(count) as count
			from watchers_fans
			where target_uid = ?
			group by uid, u_name
			order by price desc
			limit %d;
		`, limit), uid).
		Scan(&mostSpentWatchers).
		Error

	if err != nil {
		return nil, err
	}

	return mostSpentWatchers, nil
}
