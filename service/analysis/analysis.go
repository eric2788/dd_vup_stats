package analysis

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"
	"vup_dd_stats/service/db"

	"github.com/sirupsen/logrus"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

// all function here must be invoke suspended

var logger = logrus.WithField("service", "analysis")

// RecordSearchText with annoymous
func RecordSearchText(txt string) {

	if strings.TrimSpace(txt) == "" {
		return
	}

	ana := &db.SearchAnalysis{}

	hash := md5.Sum([]byte(txt))
	hashStr := hex.EncodeToString(hash[:])

	err := db.Database.FirstOrCreate(ana, db.SearchAnalysis{
		SearchText: txt,
		SearchHash: hashStr,
	}).Error

	if err != nil {
		logger.Errorf("尋找或創建匿名搜索數據時出現錯誤: %v", err)
		return
	}

	ana.AccessCount += 1
	ana.LastAccessDate = time.Now().Format(TimeFormat)

	err = db.Database.Save(ana).Error

	if err != nil {
		logger.Errorf("儲存匿名搜索數據時出現錯誤: %v", err)
		return
	}

	logger.Infof("搜索數據儲存成功: %q 已被搜索 %d 次", ana.SearchText, ana.AccessCount)
}

func RecordSearchUser(uid int64) {
	ana := &db.UserAnalysis{}

	err := db.Database.FirstOrCreate(ana, db.UserAnalysis{Uid: uid}).Error

	if err != nil {
		logger.Errorf("尋找或創建匿名vup訪問數據時出現錯誤: %v", err)
		return
	}

	ana.AccessCount += 1
	ana.LastAccessDate = time.Now().Format(TimeFormat)

	err = db.Database.Save(ana).Error

	if err != nil {
		logger.Errorf("儲存vup訪問數據時出現錯誤: %v", err)
		return
	}

	logger.Infof("vup統計數據儲存成功: %d 已被訪問 %d 次", ana.Uid, ana.AccessCount)
}
