package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"vup_dd_stats/service/blive"
)

func GetListening() (*ListeningStats, error) {
	res, err := http.Get(fmt.Sprintf("%v/listening", os.Getenv("BILIGO_HOST")))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	stats := &ListeningStats{}
	err = json.Unmarshal(b, stats)
	return stats, err
}

func GetLiveInfo(roomId int64) (*blive.LiveInfo, error) {
	res, err := http.Get(fmt.Sprintf("%v/listening/%v", os.Getenv("BILIGO_HOST"), roomId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	info := &blive.LiveInfo{}
	err = json.Unmarshal(b, info)
	return info, err
}
