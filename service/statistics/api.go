package statistics

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
	res, err := http.Get("https://api.tokyo.vtbs.moe/v1/vtbs")
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
	res, err := http.Get(fmt.Sprintf("https://api.bilibili.com/x/space/acc/info?mid=%v&jsonp=jsonp", uid))
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
	return &resp, err
}
