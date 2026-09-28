package main

import (
	"fmt"
	"github.com/gospider007/requests"
	"hwyy/Ablum"
	"hwyy/Artist"
	"hwyy/Playlist"
	"hwyy/Song"
	_ "hwyy/auth"
	"hwyy/cmd"
	"hwyy/config"
	"hwyy/log"
	"hwyy/search"
	"os"
	"regexp"
	"strings"
)

var (
	logger = log.Logger
)

func parseUrl(u string) (ty string, id string) {
	//https://url.cloud.huawei.com/mBzrI9cAqQ
	//https://portal-drcn.music.dbankcloud.cn/music-apph5-service/h5/index.html#/musicShare?songid=68897221&shareChannel=copyLink
	//https://url.cloud.huawei.com/lPX0K9pHEs
	//https://portal-drcn.music.dbankcloud.cn/music-apph5-service/h5/index.html#/musicshare/3?id=35430&listid=&shareChannel=copyLink
	//https://url.cloud.huawei.com/lWHUkdocnK
	//https://portal-drcn.music.dbankcloud.cn/music-apph5-service/h5/index.html#/musicshare/4?id=43032&listid=&shareChannel=copyLink
	//https://url.cloud.huawei.com/sDBhvCQxbO
	//https://portal-drcn.music.dbankcloud.cn/music-apph5-service/h5/index.html#/musicshare/2?id=5710928&listid=&shareChannel=copyLink
	//https://url.cloud.huawei.com/sRxeqRiQnu 索尼精选源
	//https://portal-drcn.music.dbankcloud.cn/music-apph5-service/h5/index.html#/hifisongdetail2/O9mTFxFFMUQXVKIYs?name=%E6%B3%A2%E5%8A%A8%20%28Ondulation%29&type=3&columnId=M35YcAtlCBruWutoW&columnName=%E4%B8%93%E8%BE%91%E6%8E%A8%E8%8D%90&pageType=32&isMongoliaPage=true&share=true&shareChannel=copyLink
	//https://portal-drcn.music.dbankcloud.cn/music-apph5-service/h5/index.html#/musicshare/hifi?id=NeI_inwrYTTrBQ0yO&listid=&shareChannel=copyLink
	//如果不是http开始，搜索接口
	if !strings.HasPrefix(u, "http") {
		// 如果不是以 http 开头，返回搜索类型
		return "5", u
	}
	rsp, err := requests.Get(nil, u)
	if err != nil {
		return
	}
	ul := rsp.Url().String()
	re := regexp.MustCompile(`share/(\d+).*?id=(.*?)&`)
	ul = strings.Replace(ul, "musicShare?", "musicshare/1?", -1)
	ul = strings.Replace(ul, "hifisongdetail2/", "musicshare/2?id=", -1)
	ul = strings.Replace(ul, "hifi", "1", -1)
	ul = strings.Replace(ul, "?", "&", -1)

	matches := re.FindStringSubmatch(ul)
	ty = "1"
	id = ""
	if len(matches) > 2 {
		if matches[1] != "" {
			ty = matches[1]
		}
		id = matches[2]
	} else if len(matches) == 2 {
		id = matches[1]
	}
	return

}
func init() {
	if len(os.Args) > 1 {
		cmd.Execute()
	}
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("程序异常", "异常信息", err)
		}
	}()
	var url string
	if len(os.Args) > 1 {
		url = os.Args[len(os.Args)-1]
	}
	if !strings.HasPrefix(url, "http") {
		url = ""
		fmt.Println("请输入分享链接:")
		_, err := fmt.Scanln(&url)
		if err != nil {
			return
		}
	}
	ty, id := parseUrl(url)
	var fileInfos []Song.FileInfo
	switch ty {
	case "1":
		logger.Info("获取单曲", "歌曲ID", id)
		fileInfos = Song.Download(id)
	case "2":
		config.Format = config.Conf.AlbumFormat
		logger.Info("获取专辑信息", "专辑ID", id)
		fileInfos = Ablum.GetFiles(id)
	case "3":
		logger.Info("获取歌手信息", "歌手ID", id)
		fileInfos = Artist.GetFiles(id)
	case "4":
		config.Format = config.Conf.PlaylistFormat
		logger.Info("获取歌单信息", "歌单ID", id)
		fileInfos = Playlist.GetFiles(id)
	case "5":
		logger.Info("尝试搜索")
		fileInfos = search.Search(id)
	}
	if len(fileInfos) == 0 {
		return
	}
	logger.Info("开始下载", "共计", len(fileInfos))
	Song.DownloadFiles(fileInfos)
}
