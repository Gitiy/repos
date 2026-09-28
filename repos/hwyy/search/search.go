package search

import (
	"github.com/gospider007/requests"
	"hwyy/Song"
	"hwyy/config"
	"hwyy/log"
)

var (
	logger = log.Logger

	req = config.Req
)

type Resp struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	Experiment         string   `json:"experiment"`
	SortedContentTypes []string `json:"sortedContentTypes"`
	HasNextPage        string   `json:"hasNextPage"`
	DynamicResultNums  struct {
	} `json:"dynamicResultNums"`
	QuResult struct {
		ContentType   string `json:"contentType"`
		QueryWord     string `json:"queryWord"`
		MediasearchQu []struct {
			Pattern string  `json:"pattern"`
			Weight  float64 `json:"weight"`
			Prob    float64 `json:"prob"`
		} `json:"mediasearchQu"`
		MediasearchTop1QuPattern string `json:"mediasearchTop1QuPattern"`
	} `json:"quResult"`
	TotalNum        string      `json:"totalNum"`
	SongSimpleInfos []Song.Info `json:"songSimpleInfos"`
}

func Search(word string) (files []Song.FileInfo) {
	url := "https://api-drcn.music.dbankcloud.cn/music-search-service/v10/service/fuzzysearch"
	headers := map[string]string{
		"User-Agent": "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
	}
	data := map[string]any{
		"contentType": "1",
		"correct":     "1",
		"limit":       "20",
		"queryWord":   word,
		"start":       "0",
	}
	resp, err := req.Post(nil, url, requests.RequestOption{
		Json:    data,
		Headers: headers,
	})
	if err != nil || resp.StatusCode() != 200 {
		logger.Error("请求专辑信息失败", "status", resp.StatusCode(), "err", err)
		return
	}
	var res Resp
	_, err = resp.Json(&res)
	if err != nil {
		logger.Error("解析专辑信息失败", "err", err)
		return
	}
	infos := res.SongSimpleInfos
	infos = Song.GetSong(infos)
	for _, info := range infos {
		Infos := Song.GetInfos(info)
		files = append(files, Infos...)
	}
	return

}
