package stats

import (
	"strconv"
	"strings"
	"vup_dd_stats/service/vup"
	"vup_dd_stats/service/watcher"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("controller", "stats")

func Register(group *gin.RouterGroup) {
	group.GET("", GetGlobalStats)
	group.GET("/command/:command", GetCommandStatus)
	group.GET("/:uid", GetUserStats)
	group.GET("/:uid/fans", GetUserFanStats)
	group.GET("/:uid/:command", GetUserStatsCommand)
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

	resp, err := vup.GetMostBehaviourVupsByCommand(top, c.Param("command"), price)

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
	resp, err := vup.GetStatsByType(top, statsType)

	if err != nil {
		logger.Error(err)
		code := 500
		if strings.Contains(err.Error(), "不支持") {
			code = 400
		}
		c.JSON(code, gin.H{
			"code": code,
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

func GetUserStatsCommand(c *gin.Context) {

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

	resp, err := vup.GetStatsCommand(userId, limit, command, orderPrice)

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

func GetUserFanStats(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("uid"), 10, 64)

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "uid must be a number",
		})
		return
	}

	statType := c.DefaultQuery("type", "count")

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

	resp, err := watcher.GetFanStatsForVup(userId, limit, statType)

	if err != nil {
		logger.Warn(err)
		code := 500
		if strings.Contains(err.Error(), "不支持") {
			code = 400
		}
		c.JSON(code, gin.H{
			"code": code,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

func GetUserStats(c *gin.Context) {

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

	resp, err := vup.GetStats(userId, limit)

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
