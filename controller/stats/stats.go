package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"vup_dd_stats/service/vup"
)

var logger = logrus.WithField("controller", "stats")

func Register(group *gin.RouterGroup) {
	group.GET("", GetGlobalStats)
	group.GET("/:uid", GetUserStats)
	group.GET("/:uid/:command", GetUserStatsCommand)
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

	resp, err := vup.GetGlobalStats(top)

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

	resp, err := vup.GetStatsCommand(userId, limit, command)

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
