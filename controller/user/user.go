package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"vup_dd_stats/service/vup"
)

var logger = logrus.WithField("controller", "user")

func Register(gp *gin.RouterGroup) {
	gp.GET("", GetUsers)
	gp.GET("/:uid", GetUser)
}

func GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "page must be a number",
		})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "30"))

	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "size must be a number",
		})
		return
	}

	desc := c.DefaultQuery("desc", "true") == "true"

	resp, err := vup.GetVups(page, size, desc)

	if err != nil {
		logger.Errorf("索取 vup 列表時出現錯誤: %v", err)
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

func GetUser(c *gin.Context) {
	//userId := c.Param("uid")

}
