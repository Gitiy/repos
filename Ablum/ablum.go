package Ablum

import (
	"github.com/gospider007/requests"
	"hwyy/Song"
	"hwyy/config"
	"hwyy/log"
)

var (
	logger = log.Logger
	Conf   = config.Conf
	req    = config.Req
)

type Detail struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	AlbumInfoEx struct {
		AlbumInfo       AlbumInfo   `json:"albumInfo"`
		SongSimpleInfos []Song.Info `json:"songSimpleInfos"`
	} `json:"albumInfoEx"`
	MusicListInfoEx struct {
		MusicListDetail MusicListDetail `json:"musicListDetail"`
		SongSimpleInfos []Song.Info     `json:"songSimpleInfos"`
	} `json:"musicListInfoEx"`
}
type AlbumInfo struct {
	ContentID   string `json:"contentID"`
	ContentName string `json:"contentName"`
	ContentType string `json:"contentType"`
	Picture     struct {
		BigImgURL    string `json:"bigImgURL"`
		MiddleImgURL string `json:"middleImgURL"`
		SmallImgURL  string `json:"smallImgURL"`
	} `json:"picture"`
	SubTitle      string        `json:"subTitle"`
	Status        string        `json:"status"`
	Times         string        `json:"times"`
	TotalCount    string        `json:"totalCount"`
	ArtistID      []interface{} `json:"artistID"`
	CollectTimes  string        `json:"collectTimes"`
	ContentExInfo string        `json:"contentExInfo"`
	ExtendInfo    string        `json:"extendInfo"`
	CpID          string        `json:"cpID"`
	DeviceType    string        `json:"deviceType"`
	Description   string        `json:"description"`
	Category      []struct {
		CategoryID   string `json:"categoryID"`
		CategoryName string `json:"categoryName"`
	} `json:"category"`
	OnlineCount       string   `json:"onlineCount"`
	ReleaseDate       string   `json:"releaseDate"`
	PlayNum           string   `json:"playNum"`
	CategoryNames     []string `json:"categoryNames"`
	ArtistSimpleInfos []struct {
		ContentID   string `json:"contentID"`
		ContentName string `json:"contentName"`
		ContentType string `json:"contentType"`
		Picture     struct {
			BigImgURL    string `json:"bigImgURL"`
			MiddleImgURL string `json:"middleImgURL"`
			SmallImgURL  string `json:"smallImgURL"`
		} `json:"picture"`
		KeyName        string        `json:"keyName"`
		SubTitle       string        `json:"subTitle"`
		Status         string        `json:"status"`
		Times          string        `json:"times"`
		ArtistID       []interface{} `json:"artistID"`
		ArtistName     string        `json:"artistName"`
		CollectTimes   string        `json:"collectTimes"`
		ContentExInfo  string        `json:"contentExInfo"`
		OnlineCount    string        `json:"onlineCount"`
		PlayNum        string        `json:"playNum"`
		OnlineVideoNum string        `json:"onlineVideoNum"`
		OnlineAlbumNum string        `json:"onlineAlbumNum"`
	} `json:"artistSimpleInfos"`
	CommentTimes string `json:"commentTimes"`
}
type MusicListDetail struct {
	ContentID   string `json:"contentID"`
	ContentName string `json:"contentName"`
	ContentType string `json:"contentType"`
	Picture     struct {
		BigImgURL    string `json:"bigImgURL"`
		MiddleImgURL string `json:"middleImgURL"`
		SmallImgURL  string `json:"smallImgURL"`
	} `json:"picture"`
	Times        string        `json:"times"`
	TotalCount   string        `json:"totalCount"`
	ArtistID     []interface{} `json:"artistID"`
	CollectTimes string        `json:"collectTimes"`
	ExtendInfo   string        `json:"extendInfo"`
	QualityType  string        `json:"qualityType"`
	ReviewStatus string        `json:"reviewStatus"`
	Description  string        `json:"description"`
	Category     []struct {
		CategoryID   string `json:"categoryID"`
		CategoryName string `json:"categoryName"`
	} `json:"category"`
	OnlineCount       string   `json:"onlineCount"`
	PlayNum           string   `json:"playNum"`
	CategoryNames     []string `json:"categoryNames"`
	UpdateContentTime string   `json:"updateContentTime"`
	UserIconURL       string   `json:"userIconURL"`
	SnsUserID         string   `json:"snsUserID"`
	NickName          string   `json:"nickName"`
	NotationCode      string   `json:"notationCode"`
	CommentCount      string   `json:"commentCount"`
}

func GetAlbumInfo(id string) (detail Detail, ret int) {
	url := "https://api-drcn.music.dbankcloud.cn/music-odpapp-service/v2/service/album/detail/byalbumcode"
	start := 0
	var Songs []Song.Info
	for {
		params := map[string]any{
			"contentCode":     id,
			"copyrightSwitch": "1",
			"start":           start,
			"limit":           "100",
			"contentType":     "2",
		}
		headers := map[string]string{
			"User-Agent": "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
		}
		resp, err := req.Get(nil, url, requests.RequestOption{
			Params:  params,
			Headers: headers,
		})
		if err != nil || resp.StatusCode() != 200 {
			logger.Error("请求专辑信息失败", "status", resp.StatusCode(), "err", err)
			ret = 1
		}
		_, err = resp.Json(&detail)
		if err != nil {
			logger.Error("解析专辑信息失败", "err", err)
			ret = 1
		}
		if detail.Result.ResultCode != "000000" {
			logger.Error("获取专辑信息失败", "code", detail.Result.ResultCode, "msg", detail.Result.ResultMessage)
			ret = 1
		}
		Songs = append(Songs, detail.AlbumInfoEx.SongSimpleInfos...)
		if len(detail.AlbumInfoEx.SongSimpleInfos) < 100 || len(Songs) >= Conf.MaxCount {
			break
		}
		start += 100
	}
	detail.AlbumInfoEx.SongSimpleInfos = Songs
	return
}

func GetFiles(id string) (files []Song.FileInfo) {
	detail, ret := GetAlbumInfo(id)
	if ret != 0 {
		return
	}
	alb := detail.AlbumInfoEx.AlbumInfo
	AlbumName := alb.ContentName
	AlbumDesc := alb.Description
	Aritist := alb.ArtistSimpleInfos[0].ArtistName
	logger.Info("专辑信息", "专辑名", AlbumName, "歌手", Aritist, "专辑描述", AlbumDesc)
	infos := detail.AlbumInfoEx.SongSimpleInfos
	infos = Song.GetSong(infos)
	for _, info := range infos {
		Infos := Song.GetInfos(info)
		files = append(files, Infos...)
	}
	return
}
