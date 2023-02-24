package watcher

import (
	"fmt"
	"vup_dd_stats/service/db"
)

func GetStatsByType(top int, t string) (interface{}, error) {
	switch t {
	case "count":
		return GetTotalCount()
	case "dd":
		return GetMostDDWatchers(top)
	case "behaviours":
		return GetMostBehaviourWatchers(top)
	case "spent":
		return GetMostSpentWatchers(top)
	case "famous":
		return GetMostFamousVups(top)
	case "interacted":
		return GetMostInteractedVups(top)
	case "earned":
		return GetMostEarnedVups(top)
	default:
		return nil, fmt.Errorf("unknown stats type: %q", t)
	}
}

// GetMostEarnedVups 获取从普通B站用户打赏营收最多的vups
func GetMostEarnedVups(limit int) ([]AnalysisVupInfo, error) {
	var mostEarnedVups = make([]AnalysisVupInfo, 0)

	limitStr := fmt.Sprintf("%d", limit)
	if limit == -1 {
		limitStr = "all"
	}

	err := db.Database.
		Raw(fmt.Sprintf(`
			select 
				vups.uid, 
				vups.name, 
				vups.face, 
				s.earned as count 
			from vups_with_watcher_behaviours s 
			inner join vups on vups.uid = s.uid 
			order by count desc 
			limit %s
		`, limitStr)).
		Scan(&mostEarnedVups).
		Error

	return mostEarnedVups, err
}

// GetMostDDVups 獲取最受普通B站用户欢迎的vups (最多人访问的vup)
func GetMostFamousVups(limit int) ([]AnalysisVupInfo, error) {
	var mostFamousVups = make([]AnalysisVupInfo, 0)
	limitStr := fmt.Sprintf("%d", limit)
	if limit == -1 {
		limitStr = "all"
	}

	err := db.Database.
		Raw(fmt.Sprintf(`
			select 
				vups.uid, 
				vups.name, 
				vups.face, 
				s.famous as count 
			from vups_with_watcher_behaviours s 
			inner join vups on vups.uid = s.uid 
			order by count desc 
			limit %s
		`, limitStr)).
		Scan(&mostFamousVups).
		Error
	return mostFamousVups, err
}

// GetMostInteractedVups 获取经常被普通B站用户互动的vups (被普通B站用户互动次数最多)
func GetMostInteractedVups(limit int) ([]AnalysisVupInfo, error) {
	var mostInteractedVups = make([]AnalysisVupInfo, 0)

	limitStr := fmt.Sprintf("%d", limit)
	if limit == -1 {
		limitStr = "all"
	}

	err := db.Database.
		Raw(fmt.Sprintf(`
			select 
				vups.uid, 
				vups.name, 
				vups.face, 
				s.interacted as count 
			from vups_with_watcher_behaviours s 
			inner join vups on vups.uid = s.uid 
			order by count desc 
			limit %s
		`, limitStr)).
		Scan(&mostInteractedVups).
		Error

	return mostInteractedVups, err
}

// GetMostDDWatchers 獲取進入最多不同直播間的 dd
func GetMostDDWatchers(limit int) ([]AnalysisWatcherInfo, error) {

	var mostDDWatchers = make([]AnalysisWatcherInfo, 0)

	// maximum limit is 50000
	// having 7.7M+ records in watcher_stats since 2023/2/23
	// will cause to web lag
	if limit == -1 {
		limit = 50000
	}

	err := db.Database.
		Raw(fmt.Sprintf("select uid, u_name, dd as count from watchers_stats order by count desc limit %d", limit)).
		Scan(&mostDDWatchers).
		Error

	return mostDDWatchers, err
}

func GetTotalCount() (int64, error) {
	var count int64
	err := db.Database.
		Raw(db.CountStatement, "watcher_behaviours").
		Count(&count).
		Error
	return count, err
}

// GetMostBehaviourWatchers 獲取最多行為的 dd
func GetMostBehaviourWatchers(limit int) ([]AnalysisWatcherInfo, error) {

	var mostBehaviourWatchers = make([]AnalysisWatcherInfo, 0)

	// maximum limit is 50000
	// having 7.7M+ records in watcher_stats since 2023/2/23
	// will cause to web lag
	if limit == -1 {
		limit = 50000
	}

	err := db.Database.
		Raw(fmt.Sprintf("select uid, u_name, count from watchers_stats order by count desc limit %d", limit)).
		Scan(&mostBehaviourWatchers).
		Error

	return mostBehaviourWatchers, err
}

// GetMostSpentWatchers 獲取花費最多的 dd
func GetMostSpentWatchers(limit int) ([]AnalysisWatcherInfo, error) {

	var mostSpentWatchers = make([]AnalysisWatcherInfo, 0)

	// maximum limit is 50000
	// having 7.7M+ records in watcher_stats since 2023/2/23
	// will cause to web lag
	if limit == -1 {
		limit = 50000
	}

	err := db.Database.
		Raw(fmt.Sprintf("select uid, u_name, spent as price from watchers_stats order by price desc limit %d", limit)).
		Scan(&mostSpentWatchers).
		Error

	return mostSpentWatchers, err
}

// GetMostBehaviourWatchersByCommand 獲取最多行為的 dd
// Still have performance issue
func GetMostBehaviourWatchersByCommand(limit int, command string, price bool) ([]AnalysisWatcherInfo, error) {
	var mostDDBehaviourVups = make([]AnalysisWatcherInfo, 0)

	orderBy := "count"
	if price {
		orderBy = "price"
	}

	// maximum limit is 50000
	// having 7.7M+ records in watcher_stats since 2023/2/23
	if limit == -1 {
		limit = 50000
	}

	err := db.Database.
		Raw(fmt.Sprintf(`
			select
				uid,
				u_name,
				SUM(count) as count,
				SUM(price) as price
			from watchers_fans
			where command = ?
			group by target_uid, uid, u_name
			order by %s desc
			limit %d;
		`, orderBy, limit), command).
		Scan(&mostDDBehaviourVups).
		Error

	return mostDDBehaviourVups, err
}
