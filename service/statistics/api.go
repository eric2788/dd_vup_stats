package statistics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"vup_dd_stats/service/blive"

	"github.com/corpix/uarand"
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

func GetVtbListLaplace() (map[int64]VupData, error) {
	// legacy: https://vup-json.bilibili.ooo/vup.json
	res, err := httpGet("https://vup-json.laplace.live/vup-slim.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp map[string]VupJsonData
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	var results = make(map[int64]VupData)
	for k, v := range resp {
		uid, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			logger.Errorf("parse uid %v failed for %s: %v, skipped", k, v.Name, err)
			continue
		}
		results[uid] = VupData{
			Name:   v.Name,
			RoomId: v.RoomId,
		}
	}
	return results, nil
}

func GetVtbListVtbMoe() (map[int64]VupData, error) {
	res, err := httpGet("https://api.vtbs.moe/v1/short")
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
	if err != nil {
		return nil, err
	}
	var results = make(map[int64]VupData)
	for _, v := range resp {
		results[v.Mid] = VupData{
			Name:   v.UName,
			RoomId: v.RoomId,
		}
	}
	return results, nil
}

func GetUserInfo(uid int64) (*UserResp, error) {
	return GetUserInfoRetry(uid, 0, 5)
}

func GetUserInfoRetry(uid int64, times, max int) (*UserResp, error) {
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
		if times > max {
			logger.Warnf("Retried %v times, returning error", max)
			return nil, errors.New(resp.Message)
		} else {
			logger.Warnf("User-Agent is blocked, retrying with another one")
		}
		times += 1
		return GetUserInfoRetry(uid, times, max)
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
