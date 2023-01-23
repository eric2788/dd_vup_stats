package watcher

// SaveWatcher fetch user info from api.bilibili.com and save it to database
// before that, will check whether the user is in database or not
// if not, then save it
func SaveWatcher(uid int64) {
	// TODO
}

// MigrateToVup will migrate specific watcher to vup
// it will first move the records of the watcher from watcher_behaviour to behaviour table
// then delete the records of that watcher in watcher_behaviour table
// finally, delete the watcher in watcher table
// but one thing important is that, first needs to check whether the uid is in watcher table or not
func MigrateToVup(uid int64) {
	// TODO
}
