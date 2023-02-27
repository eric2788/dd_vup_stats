package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"vup_dd_stats/service/vup"
)

var logger = logrus.WithField("controller", "user")

var orderAllows = []string{
	"last_listened_at",
	"dd_count",
	"last_behaviour_at",
	"first_listen_at",
	"total_spent",
}

func Register(gp *gin.RouterGroup) {
	gp.GET("", GetUsers)
	gp.GET("/:uid", GetUser)
}

func GetUsers(c *gin.Context) {

	searchStr := c.DefaultQuery("q", "")

	// 如果有 uid: 前綴，就只搜尋 uid (laplace.live caused)
	searchStr = strings.TrimPrefix(strings.ToLower(searchStr), "uid:")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "page must be a number",
		})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "30"))

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "size must be a number",
		})
		return
	}

	// 每頁最高拿到50筆
	if size > 50 {
		size = 50
	} else if size <= 0 {
		size = 0
	}

	desc := c.DefaultQuery("desc", "true") == "true"
	orderBy := c.DefaultQuery("orderBy", "last_listened_at")

	if !slices.Contains(orderAllows, orderBy) {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "orderBy must be one of " + strings.Join(orderAllows, ", "),
		})
		return
	}

	resp, err := vup.SearchVups(searchStr, page, size, orderBy, desc)

	if err != nil {
		logger.Errorf("搜索 vup 列表時出現錯誤: %v", err)
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
	userId, err := strconv.ParseInt(c.Param("uid"), 10, 64)

	if err != nil {
		logger.Warn(err)
		c.JSON(400, gin.H{
			"code":    400,
			"message": "UID 必须为数字",
		})
		return
	}

	resp, err := vup.GetVup(userId)

	if err != nil {
		logger.Errorf("索取 vup 時出現錯誤: %v", err)
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
