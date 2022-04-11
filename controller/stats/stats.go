package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

var logger = logrus.WithField("controller", "stats")

func Register(group *gin.RouterGroup) {
	group.GET("", GetGlobalStats)
	group.GET("/:uid", GetUserStats)
}

func GetGlobalStats(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
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

	//TODO
	c.JSON(200, gin.H{
		"message": userId,
	})

}
