package watcher

import (
	"strconv"
	"vup_dd_stats/service/watcher"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//TODO create controller of watchers

var logger = logrus.WithField("controller", "watcher")


// test only, so only stats/:uid for now
func Register(group *gin.RouterGroup) {
	group.GET("/stats/:uid", GetWatcherStats)
}

func GetWatcherStats(c *gin.Context) {

	userId, err := strconv.ParseInt(c.Param("uid"), 10, 64)

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "uid must be a number",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "limit must be a number",
		})
		return
	}

	// 最高拿到 50
	if limit > 50 {
		limit = 50
	} else if limit <= 0 {
		limit = 1
	}

	concurrent := c.DefaultQuery("concurrent", "false") != "false"

	var getStats = watcher.GetStats

	if concurrent {
		getStats = watcher.GetStatsConcurrent
	}

	resp, err := getStats(userId, limit)

	if err != nil {
		logger.Error(err)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": resp,
	})
}
