package Playlist

import (
	"github.com/gospider007/requests"
	"hwyy/Ablum"
	"hwyy/Song"
	"hwyy/config"
	"hwyy/log"
)

var (
	logger = log.Logger
	Conf   = config.Conf
	req    = config.Req
)

func GetPlayListInfo(id string) (detail Ablum.Detail, ret int) {
	url := "https://api-drcn.music.dbankcloud.cn/music-operation-service/v2/service/musiclist/detail/bymusiclistid"
	start := 0
	var Songs []Song.Info
	for {
		data := map[string]any{
			"copyrightSwitch": "1",
			"limit":           "100",
			"musicListID":     id,
			"sortType":        "0",
			"start":           start,
		}
		headers := map[string]string{
			"User-Agent": "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
		}
		resp, err := req.Post(nil, url, requests.RequestOption{
			Json:    data,
			Headers: headers,
		})
		if err != nil || resp.StatusCode() != 200 {
			logger.Error("请求歌单信息失败", "status", resp.StatusCode(), "err", err)
			ret = 1
			return
		}
		_, err = resp.Json(&detail)
		if err != nil {
			logger.Error("解析歌单信息失败", "err", err)
			ret = 1
			return
		}
		if detail.Result.ResultCode != "000000" {
			logger.Error("获取歌单信息失败", "code", detail.Result.ResultCode, "msg", detail.Result.ResultMessage)
			ret = 1
			return
		}
		Songs = append(Songs, detail.MusicListInfoEx.SongSimpleInfos...)
		if len(detail.MusicListInfoEx.SongSimpleInfos) < 100 || len(Songs) >= Conf.MaxCount {
			break
		}
		start += 100
	}
	detail.MusicListInfoEx.SongSimpleInfos = Songs
	return
}

func GetFiles(id string) (files []Song.FileInfo) {
	detail, ret := GetPlayListInfo(id)
	if ret != 0 {
		return
	}
	musicListDetail := detail.MusicListInfoEx.MusicListDetail
	musicListName := musicListDetail.ContentName
	musicListDesc := musicListDetail.Description
	logger.Info("歌单信息", "歌单名", musicListName, "歌单描述", musicListDesc)
	infos := detail.MusicListInfoEx.SongSimpleInfos
	infos = Song.GetSong(infos)
	for _, info := range infos {
		info.ListName = musicListName
		Infos := Song.GetInfos(info)
		files = append(files, Infos...)
	}
	return
}
