package analysis

import (
	"crypto/md5"
	"encoding/hex"
	"vup_dd_stats/service/db"

	"github.com/sirupsen/logrus"
)


// all function here must be invoke suspended

var logger = logrus.WithField("service", "analysis")

// RecordSearchText with annoymous
func RecordSearchText(txt string) {

	ana := &db.SearchAnalysis{}

	hash := md5.Sum([]byte(txt))
	hashStr := hex.EncodeToString(hash[:])

	err := db.Database.FirstOrCreate(&ana, db.SearchAnalysis{
		SearchText: txt,
		SearchHash: hashStr,
	}).Error

	if err != nil {
		logger.Errorf("尋找或創建匿名搜索數據時出現錯誤: %v", err)
		return
	}

	ana.AccessCount += 1

	err = db.Database.Save(ana).Error

	if err != nil {
		logger.Errorf("儲存匿名搜索數據時出現錯誤: %v", err)
		return
	}

	logger.Infof("搜索數據儲存成功: %q 已被搜索 %d 次", ana.SearchText, ana.AccessCount)
}


func RecordSearchUser(uid int64){
	ana := &db.UserAnalysis{}

	err := db.Database.FirstOrCreate(&ana, db.UserAnalysis{ Uid: uid })

	if err != nil {
		logger.Errorf("尋找或創建匿名vup訪問數據時出現錯誤: %v", err)
		return
	}

	ana.AccessCount += 1

	err = db.Database.Save(ana)

	if err != nil {
		logger.Errorf("儲存vup訪問數據時出現錯誤: %v", err)
		return
	}

	logger.Infof("vup統計數據儲存成功: %d 已被訪問 %d 次", ana.Uid, ana.AccessCount)
}


