package Artist

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

type Detail struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	ArtistInfoEx struct {
		Info             Info        `json:"artistInfo"`
		SongSimpleInfos  []Song.Info `json:"songSimpleInfos"`
		DigitalAlbumInfo struct {
			AlbumCode string `json:"albumCode"`
			AlbumName string `json:"albumName"`
			Price     string `json:"price"`
			BuyCount  string `json:"buyCount"`
		} `json:"digitalAlbumInfo"`
		AlbumSimpleInfos []Ablum.AlbumInfo `json:"albumSimpleInfos"`
		VideoSimpleInfos []struct {
			ContentID   string `json:"contentID"`
			ContentName string `json:"contentName"`
			ContentType string `json:"contentType"`
			Picture     struct {
				BigImgURL    string `json:"bigImgURL"`
				MiddleImgURL string `json:"middleImgURL"`
				SmallImgURL  string `json:"smallImgURL"`
			} `json:"picture"`
			KeyName        string   `json:"keyName"`
			SubTitle       string   `json:"subTitle"`
			Status         string   `json:"status"`
			Quality        string   `json:"quality"`
			Times          string   `json:"times"`
			AlbumID        string   `json:"albumID"`
			ArtistID       []string `json:"artistID"`
			AlbumName      string   `json:"albumName"`
			ArtistName     string   `json:"artistName"`
			CollectTimes   string   `json:"collectTimes"`
			RelationShips  string   `json:"relationShips,omitempty"`
			ContentExInfo  string   `json:"contentExInfo"`
			ExtendInfo     string   `json:"extendInfo"`
			CpID           string   `json:"cpID"`
			DeviceType     string   `json:"deviceType"`
			SecondSubTitle string   `json:"secondSubTitle"`
			Category       []struct {
				CategoryID string `json:"categoryID"`
			} `json:"category"`
			ReleaseDate   string `json:"releaseDate"`
			DiversionFlag string `json:"diversionFlag"`
			SongType      string `json:"songType"`
		} `json:"videoSimpleInfos"`
	} `json:"artistInfoEx"`
	SongSimpleInfos []Song.Info `json:"songSimpleInfos"`
}
type Info struct {
	ContentID   string `json:"contentID"`
	ContentName string `json:"contentName"`
	ContentType string `json:"contentType"`
	Picture     struct {
		BigImgURL    string `json:"bigImgURL"`
		MiddleImgURL string `json:"middleImgURL"`
		SmallImgURL  string `json:"smallImgURL"`
	} `json:"picture"`
	SubTitle      string        `json:"subTitle"`
	Times         string        `json:"times"`
	ArtistID      []interface{} `json:"artistID"`
	CollectTimes  string        `json:"collectTimes"`
	Description   string        `json:"description"`
	OnlineCount   string        `json:"onlineCount"`
	PlayNum       string        `json:"playNum"`
	CategoryNames []string      `json:"categoryNames"`
	SongNum       string        `json:"songNum"`
	AlbumNum      string        `json:"albumNum"`
	VideoNum      string        `json:"videoNum"`
}

func GetArtistInfo(id string) (detail Detail, ret int) {
	url := "https://api-drcn.music.dbankcloud.cn/music-metacontent-service/v2/service/artist/detail/byartistcode"
	start := 0
	var Songs []Song.Info
	for {
		data := map[string]interface{}{
			"albumNum":        "150",
			"albumType":       "",
			"artistCode":      id,
			"contentSubType":  "0|2",
			"copyrightSwitch": "1",
			"mvType":          "",
			"songNum":         "100",
			"songOrderBy":     "0", //1:最新 0:最热
			"videoNum":        "10",
		}
		if Conf.ArtistType == "s" {
			data = map[string]interface{}{
				"artistCode":      id,
				"copyrightSwitch": "1",
				"limit":           "100",
				"songOrderBy":     "0", //0:热度 1:时间
				"songType":        "0|2",
				"start":           start,
			}
			url = "https://api-drcn.music.dbankcloud.cn/music-metacontent-service/v2/service/song/simple/byartistcode"
		}
		headers := map[string]string{
			"User-Agent": "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
		}
		resp, err := req.Post(nil, url, requests.RequestOption{
			Json:    data,
			Headers: headers,
		})
		if err != nil || resp.StatusCode() != 200 {
			logger.Error("请求歌手信息失败", "status", resp.StatusCode(), "err", err)
			ret = 1
			return
		}
		_, err = resp.Json(&detail)
		if err != nil {
			logger.Error("解析歌手信息失败", "err", err)
			ret = 1
			return
		}
		if detail.Result.ResultCode != "000000" {
			logger.Error("获取歌手信息失败", "code", detail.Result.ResultCode, "msg", detail.Result.ResultMessage)
			ret = 1
			return
		}
		Songs = append(Songs, detail.SongSimpleInfos...)
		if len(detail.SongSimpleInfos) < 100 || len(Songs) >= Conf.MaxCount || Conf.ArtistType == "a" {
			break
		}
	}
	if detail.ArtistInfoEx.Info.ContentName == "" {
		url = "https://api-drcn.music.dbankcloud.cn/music-metacontent-service/v2/service/artist/detail/byartistcode"
		data := map[string]interface{}{
			"albumNum":        "150",
			"albumType":       "",
			"artistCode":      id,
			"contentSubType":  "0|2",
			"copyrightSwitch": "1",
			"mvType":          "",
			"songNum":         "100",
			"songOrderBy":     "0", //1:最新 0:最热
			"videoNum":        "10",
		}
		headers := map[string]string{
			"User-Agent": "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
		}
		resp, err := req.Post(nil, url, requests.RequestOption{
			Json:    data,
			Headers: headers,
		})
		if err != nil || resp.StatusCode() != 200 {
			logger.Error("请求歌手信息失败", "status", resp.StatusCode(), "err", err)
			ret = 1
			return
		}
		_, err = resp.Json(&detail)
		if err != nil {
			logger.Error("解析歌手信息失败", "err", err)
			ret = 1
			return
		}

	}
	detail.ArtistInfoEx.SongSimpleInfos = Songs

	return
}
func GetFiles(id string) (files []Song.FileInfo) {
	detail, ret := GetArtistInfo(id)
	if ret != 0 {
		return
	}
	ArtistInfoEx := detail.ArtistInfoEx.Info
	ArtistName := ArtistInfoEx.ContentName
	AblumNum := ArtistInfoEx.AlbumNum
	logger.Info("歌手信息", "歌手名", ArtistName, "专辑数", AblumNum)
	infos := detail.ArtistInfoEx.SongSimpleInfos
	if Conf.ArtistType == "s" {
		config.Format = config.Conf.ArtistSingleFormat
		Songs := Song.GetSong(infos)
		for _, info := range Songs {
			Infos := Song.GetInfos(info)
			files = append(files, Infos...)
		}
	} else {
		config.Format = config.Conf.ArtistAlbumFormat
		Albums := detail.ArtistInfoEx.AlbumSimpleInfos
		Albums = Song.GetRange(Albums, Conf.AblumRange)
		logger.Info("根据歌手专辑获取歌曲", "专辑数", len(Albums))
		for _, album := range Albums {
			file := Ablum.GetFiles(album.ContentID)
			files = append(files, file...)
		}
	}
	return
}
