package watcher

import (
	"strconv"
	"vup_dd_stats/service/watcher"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("controller", "watcher")

func Register(group *gin.RouterGroup) {

	// disabled global statistics for now
	// due to slow and high impact on performance
	//group.GET("/stats", GetGlobalStats)
	//group.GET("/stats/command/:command", GetCommandStatus)

	group.GET("/stats/:uid", GetWatcherStats)
	group.GET("/stats/:uid/:command", GetWatcherStatsCommand)
	group.GET("/record/:uid", GetWatcherRecord)
	group.GET("/:uid", GetWatcher)
}

func GetCommandStatus(c *gin.Context) {
	top, err := strconv.Atoi(c.DefaultQuery("top", "3"))

	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "top must be an integer",
		})
		return
	}

	// 最高拿到 50
	if top > 50 {
		top = 50
	} else if top <= 0 {
		top = 3
	}

	price := c.DefaultQuery("price", "false") != "false"

	resp, err := watcher.GetMostBehaviourWatchersByCommand(top, c.Param("command"), price)

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

func GetWatcherStatsCommand(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	command := c.Param("command")

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

	orderPrice := c.DefaultQuery("price", "false") != "false"

	resp, err := watcher.GetStatsCommand(userId, limit, command, orderPrice)

	if err != nil {
		logger.Warn(err)
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

func GetWatcherRecord(c *gin.Context) {

	userId, err := strconv.ParseInt(c.Param("uid"), 10, 64)

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "uid must be a number",
		})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "30"))

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "pageSize must be a number",
		})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "page must be a number",
		})
		return
	}

	// 每頁最高拿到50筆
	if pageSize > 50 {
		pageSize = 50
	} else if pageSize <= 0 {
		pageSize = 0
	}

	command := c.DefaultQuery("command", "")

	records, err := watcher.GetRecords(userId, page, pageSize, command)

	if err != nil {
		logger.Warn(err)
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    records,
	})
}

func GetWatcher(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("uid"), 10, 64)

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "UID 必须为数字",
		})
		return
	}

	resp, err := watcher.GetWatcher(userId)

	if err != nil {
		logger.Errorf("索取 watcher 時出現錯誤: %v", err)
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	if resp == nil {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    resp,
		})
	}
}

func GetGlobalStats(c *gin.Context) {
	top, err := strconv.Atoi(c.DefaultQuery("top", "3"))

	if err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "top must be an integer",
		})
		return
	}

	// 最高拿到 50
	if top > 50 {
		top = 50
	} else if top <= 0 {
		top = 3
	}

	statsType := c.DefaultQuery("type", "")
	resp, err := watcher.GetStatsByType(top, statsType)

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

	resp, err := watcher.GetStats(userId, limit)

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
