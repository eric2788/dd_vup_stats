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

var nonStatsCommands = set.FromArray([]string{Live, Preparing})

func GetRegisteredCommands() []string {
	commands := make([]string, 0, len(handlerMap))
	for command := range handlerMap {
		if nonStatsCommands.Has(command) {
			continue
		}
		commands = append(commands, command)
	}
	return commands
}

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
		logger.Debugf("未知的數據類型: %s", liveData.Command)
		exceptions.Add(liveData.Command)
	}
}

func RegisterHandler(command string, handler HandlerFunc) {
	handlerMap[command] = handler
}
