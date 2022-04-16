package records

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
)

type RecordGetter func(uid int64, limit int) ([]db.Behaviour, error)

var (
	logger        = logrus.WithField("controller", "records")
	recordTypeMap = map[string]RecordGetter{
		"self":  vup.GetTopSelfRecords,
		"guest": vup.GetTopGuestRecords,
		"dd":    vup.GetTopDDRecords,
	}
)

func Register(group *gin.RouterGroup) {
	group.GET("/:uid", GetRecordsByType)
	group.GET("", GetGlobalRecords)
}

func GetGlobalRecords(c *gin.Context) {

	query := c.DefaultQuery("q", "")

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

	showSelf := c.DefaultQuery("showSelf", "true") == "true"

	records, err := vup.GetGlobalRecords(query, page, pageSize, showSelf)

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

func GetRecordsByType(c *gin.Context) {

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

	getterType := c.DefaultQuery("type", "dd")

	// 最高拿到100
	if limit > 100 {
		limit = 100
	}

	if getter, ok := recordTypeMap[getterType]; ok {
		records, err := getter(userId, limit)
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
	} else {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "type not found",
		})
	}
}
