package blive

import (
	"encoding/json"
	"vup_dd_stats/utils/set"
)

type HandlerFunc func(data *LiveData) error

var (
	handlerMap = make(map[string]HandlerFunc)
	exceptions = set.New[string]()
)

func handleMessage(b []byte) {

	liveData := &LiveData{}

	if err := json.Unmarshal(b, liveData); err != nil {
		logger.Warnf("解析 JSON 數據時出現錯誤: %v", err)
		return
	}

	if handler, ok := handlerMap[liveData.Command]; ok {
		if err := handler(liveData); err != nil {
			logger.Warnf("處理數據 %v 時出現錯誤: %v", liveData.Command, err)
		}
	} else if !exceptions.Has(liveData.Command) {
		logger.Warnf("未知的數據類型: %s", liveData.Command)
		exceptions.Add(liveData.Command)
	}
}

func RegisterHandler(command string, handler HandlerFunc) {
	handlerMap[command] = handler
}
