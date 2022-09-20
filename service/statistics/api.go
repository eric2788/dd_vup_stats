package statistics

import (
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
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

func GetListeningInfo(roomId int64) (*blive.ListeningInfo, error) {
	res, err := http.Get(fmt.Sprintf("%v/listening/%v", os.Getenv("BILIGO_HOST"), roomId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	info := &blive.ListeningInfo{}
	err = json.Unmarshal(b, info)
	return info, err
}

func GetVtbListVtbMoe() ([]VtbsMoeResp, error) {
	res, err := httpGet("https://api.tokyo.vtbs.moe/v1/vtbs")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp []VtbsMoeResp
	err = json.Unmarshal(b, &resp)
	return resp, err
}

func GetUserInfo(uid int64) (*UserResp, error) {
	res, err := httpGet(fmt.Sprintf("https://api.bilibili.com/x/space/acc/info?mid=%v&jsonp=jsonp", uid))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp UserResp
	err = json.Unmarshal(b, &resp)
	// change other agent
	if resp.Code == -401 {
		logger.Warnf("User-Agent is blocked, retrying with another one")
		return GetUserInfo(uid)
	}
	return &resp, err
}

// httpGet with user-agent
func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	logger.Debugf("Using User-Agent: %v\n", req.Header.Get("User-Agent"))
	return http.DefaultClient.Do(req)
}
