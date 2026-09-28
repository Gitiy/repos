package Song

import (
	"encoding/json"
	"fmt"
	"github.com/gospider007/requests"
	"github.com/olekukonko/tablewriter"
	"hwyy/KeyStore"
	"hwyy/config"
	"hwyy/log"
	"os"
	"sort"
	"strings"
)

var (
	qualityMap = map[string]string{
		"1":  "标准品质",
		"2":  "HQ 高品质",
		"3":  "SQ 无损品质",
		"7":  "SQ 无损品质",
		"4":  "索尼精选 hifi",
		"5":  "Hi-Res 品质",
		"13": "Audio Vivid",
		"15": "多轨道",
		"17": "多轨道",
	}
	logger = log.Logger
	Conf   = config.Conf
	req    = config.Req
)

type mp3File struct {
	BakFileURls   []interface{} `json:"bakFileURls"`
	CopyrightType string        `json:"copyrightType"`
	FileCode      string        `json:"fileCode"`
	FileInfos     []struct {
		CopyrightType string `json:"copyrightType"`
		FileCode      string `json:"fileCode"`
		FileSize      string `json:"fileSize"`
		FileType      string `json:"fileType"`
		FileURL       string `json:"fileURL"`
		MediaFileType string `json:"mediaFileType"`
	} `json:"fileInfos"`
	FileURL     string `json:"fileURL"`
	PlayQuality string `json:"playQuality"`
	Result      struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	SongInfo  Info   `json:"songInfo"`
	Type      string `json:"type"`
	SecretKey string `json:"secretKey"`
	Iv        string `json:"iv"`
}
type Picture struct {
	BigImgURL    string `json:"bigImgURL"`
	MiddleImgURL string `json:"middleImgURL"`
	SmallImgURL  string `json:"smallImgURL"`
}
type Info struct {
	AlbumID       string `json:"albumID"`
	AlbumNamePair struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"albumNamePair"`
	ArtistID        []string `json:"artistID"`
	ArtistNamePairs []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"artistNamePairs"`
	Category []struct {
		CategoryID string `json:"categoryID"`
	} `json:"category"`
	CollectTimes     string  `json:"collectTimes"`
	CommentTimes     string  `json:"commentTimes"`
	ContentExInfo    string  `json:"contentExInfo"`
	ContentID        string  `json:"contentID"`
	ContentName      string  `json:"contentName"`
	ContentType      string  `json:"contentType"`
	ExtendInfo       string  `json:"extendInfo"`
	LyricAddress     string  `json:"lyricAddress"`
	Lyrics           string  `json:"lyrics"`
	Picture          Picture `json:"picture"`
	Quality          string  `json:"quality"`
	QualityNamePairs []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"qualityNamePairs"`
	SecondSubTitle string `json:"secondSubTitle"`
	SongSubName    string `json:"songSubName"`
	Status         string `json:"status"`
	SubTitle       string `json:"subTitle"`
	Times          string `json:"times"`
	VisitControl   string `json:"visitControl"`
	ListName       string `json:"listName"`
}
type ContentExInfo struct {
	AccompanySongId string   `json:"accompanySongId"`
	AlbumCode       string   `json:"albumCode"`
	AlbumName       string   `json:"albumName"`
	ArtistCodes     []string `json:"artistCodes"`
	ArtistNames     string   `json:"artistNames"`
	Associations    string   `json:"associations"`
	AuthorInfos     []struct {
		Code       string `json:"code"`
		Name       string `json:"name"`
		PictureURL string `json:"pictureURL"`
		Type       string `json:"type"`
	} `json:"authorInfos"`
	CompleteFileInfos json.RawMessage `json:"completeFileInfos"`
	Composer          string          `json:"composer"`
	DeviceTypes       []string        `json:"deviceTypes"`
	DimmedField       []interface{}   `json:"dimmedField"`
	Exclusivity       string          `json:"exclusivity"`
	Explicit          string          `json:"explicit"`
	FileInfos         FileInfos       `json:"fileInfos"`
	FragmentInfo      struct {
		AuditionBegin    string `json:"auditionBegin"`
		AuditionEnd      string `json:"auditionEnd"`
		AuditionTime     string `json:"auditionTime"`
		AuditionUrl      string `json:"auditionUrl"`
		FragmentUrlInfos []struct {
			AuditionBegin string `json:"auditionBegin"`
			AuditionEnd   string `json:"auditionEnd"`
			AuditionTime  string `json:"auditionTime"`
			AuditionUrl   string `json:"auditionUrl"`
			QualityType   string `json:"qualityType"`
			AuditionDesc  string `json:"auditionDesc,omitempty"`
		} `json:"fragmentUrlInfos"`
		HasFragement string `json:"hasFragement"`
		MaxQuality   string `json:"maxQuality"`
	} `json:"fragmentInfo"`
	HasSample      string `json:"hasSample"`
	HasVividSample string `json:"hasVividSample"`
	LogoType       string `json:"logoType"`
	LyricAddres    string `json:"lyricAddres"`
	Lyrics         string `json:"lyrics"`
	MinPrivilege   struct {
		Download  string `json:"download"`
		Streaming string `json:"streaming"`
	} `json:"minPrivilege"`
	MvID              string    `json:"mvID"`
	Order             string    `json:"order"`
	Original          string    `json:"original"`
	OuterSongCode     string    `json:"outerSongCode"`
	OuterSongCodetype string    `json:"outerSongCodetype"`
	PlayFragmentFlag  string    `json:"playFragmentFlag"`
	Quality           string    `json:"quality"`
	QualityV2         string    `json:"qualityV2"`
	RingtoneEndTime   string    `json:"ringtoneEndTime"`
	RingtoneStartTime string    `json:"ringtoneStartTime"`
	SampleUrl         string    `json:"sampleUrl"`
	SciFiFileInfos    FileInfos `json:"sciFiFileInfos"`
	SongStyleCodes    string    `json:"songStyleCodes"`
	SongType          string    `json:"songType"`
	SoundEffects      []string  `json:"soundEffects"`
	Substatus         string    `json:"substatus"`
	TempType          string    `json:"tempType"`
	Usagetype         string    `json:"usagetype"`
	Vip               string    `json:"vip"`
	VividFileInfos    FileInfos `json:"vividFileInfos"`
	VividSampleUrl    string    `json:"vividSampleUrl"`
}
type FileInfo struct {
	BitRate       string `json:"bitRate"`
	Download      string `json:"download"`
	Drm           string `json:"drm"`
	Duration      string `json:"duration"`
	FileSize      string `json:"fileSize"`
	FileType      string `json:"fileType"`
	Format        string `json:"format"`
	MediaFileType string `json:"mediaFileType"`
	PreviewURL    string `json:"previewURL,omitempty"`
	Quality       string `json:"quality"`
	SampleRate    string `json:"sampleRate"`
	Streaming     string `json:"streaming"`
	Url           string `json:"url"`
	AlbumName     string `json:"albumName"`
	SongName      string `json:"songName"`
	SubTitle      string `json:"subTitle"`
	ArtistName    string `json:"artistName"`
	Lyrics        string `json:"lyrics"`
	Picture       string `json:"picture"`
	ListName      string `json:"listName"`
	IV            string `json:"iv"`
	SecretKey     string
}

func (f *FileInfo) GetFileSize() string {
	size := f.FileSize
	if size == "" {
		return "0 MB"
	}
	intSize := 0
	_, _ = fmt.Sscanf(size, "%d", &intSize)
	return fmt.Sprintf("%.2f MB", float64(intSize)/1024/1024)

}

func (f *FileInfo) GetDuration() string {
	duration := f.Duration
	if duration == "" {
		return "00：00"
	}
	intDuration := 0
	_, _ = fmt.Sscanf(duration, "%d", &intDuration)
	hour := intDuration / 3600
	minute := (intDuration - hour*3600) / 60
	second := intDuration - hour*3600 - minute*60
	if hour > 0 {
		return fmt.Sprintf("%02d：%02d：%02d", hour, minute, second)
	}
	return fmt.Sprintf("%02d：%02d", minute, second)

}

type FileInfos []FileInfo

func (f FileInfos) Len() int {
	return len(f)
}
func (f FileInfos) Less(i, j int) bool {
	sizeI := 0
	sizeJ := 0
	_, _ = fmt.Sscanf(f[i].FileSize, "%d", &sizeI)
	_, _ = fmt.Sscanf(f[j].FileSize, "%d", &sizeJ)
	return sizeI > sizeJ
}
func (f FileInfos) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
func (f FileInfos) FindQuality(q string) (int, bool) {
	for k, v := range f {
		if v.Quality == q {
			return k, true
		}
	}
	return 0, false

}
func parseCompleteFileInfos(raw json.RawMessage) (FileInfos, error) {
	var mapData map[string]FileInfo
	if err := json.Unmarshal(raw, &mapData); err == nil {
		var infos FileInfos
		for _, v := range mapData {
			infos = append(infos, v)
		}
		return infos, nil
	}

	var sliceData []FileInfo
	if err := json.Unmarshal(raw, &sliceData); err == nil {
		return sliceData, nil
	}
	return nil, fmt.Errorf("invalid format for CompleteFileInfos")
}
func getQualityNamePairs(quality string) string {
	if v, ok := qualityMap[quality]; ok {
		return v
	}
	return quality
}
func getFile(contentCode string, qualityType string) (mp3file mp3File, ret int) {
	url := "https://api-drcn.music.dbankcloud.cn/music-play-service/v3/service/file/bycontentcode"
	headers := map[string]string{
		"User-Agent":    "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
		"authorization": Conf.Authorization,
	}
	publicKey := KeyStore.ExportPublicKey()
	data := map[string]interface{}{
		"algType":     "1",
		"contentCode": contentCode,
		"contentType": "1",
		"publicKey":   publicKey,
		"qualityType": qualityType,
	}
	resp, err := req.Post(nil, url, requests.RequestOption{
		Headers: headers,
		Json:    data,
	})
	if err != nil || resp.StatusCode() != 200 {
		logger.Error("请求音频文件失败", "status", resp.StatusCode(), "err", err)
		ret = 1
	}
	_, err = resp.Json(&mp3file)
	if err != nil {
		logger.Error("解析音频文件失败", "err", err)
		ret = 1
	}
	if mp3file.Result.ResultCode != "000000" {
		logger.Error("获取音频文件失败", "SongID", contentCode, "ResultMessage", mp3file.Result.ResultMessage)
		ret = 1
	}
	return
}
func GetInfos(songInfo Info) (fileInfos FileInfos) {
	contentCode := songInfo.ContentID
	var contentExInfo ContentExInfo
	err := json.Unmarshal([]byte(songInfo.ContentExInfo), &contentExInfo)
	if err != nil {
		logger.Error("解析音频文件失败", "err", err)
		return
	}
	SubTitle := songInfo.SubTitle
	Title := songInfo.ContentName
	AlbumName := contentExInfo.AlbumName
	Lyrics := songInfo.LyricAddress
	logger.Info("歌曲信息", "SongID", contentExInfo.OuterSongCode, "Title", Title, "ArtistNames", contentExInfo.ArtistNames, "AlbumName", AlbumName)
	infos, err := parseCompleteFileInfos(contentExInfo.CompleteFileInfos)
	if err != nil {
		logger.Error("解析音频文件失败", "err", err)
		return
	}
	if len(infos) == 0 {
		logger.Error("音频文件信息为空", "SongID", contentExInfo.OuterSongCode)
		return
	}
	sort.Sort(infos)
	q := Conf.Quality
	if q == "all" {
		fileInfos = infos
	} else if q == "best" {
		fileInfos = append(fileInfos, infos[0])
	} else if qInt, ok := infos.FindQuality(q); ok {
		fileInfos = append(fileInfos, infos[qInt])
	}
	if len(fileInfos) == 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Quality", "FileSize", "FileType"})
		for k, v := range infos {
			table.Append([]string{fmt.Sprintf("%d", k+1), getQualityNamePairs(v.Quality), v.GetFileSize(), v.MediaFileType})
		}
		table.SetRowLine(true)
		table.SetRowSeparator("-")
		table.SetCenterSeparator("*")
		table.SetColumnSeparator("|")
		table.SetBorder(true)
		table.Render()
		var id string
		fmt.Print("请输入要下载的音频文件ID,多个ID用逗号分隔:")
		_, _ = fmt.Scan(&id)
		ids := strings.Split(id, ",")
		for _, v := range ids {
			intV := 0
			_, _ = fmt.Sscanf(v, "%d", &intV)
			fileInfos = append(fileInfos, infos[intV-1])
		}
	}
	for i, fileInfo := range fileInfos {
		if fileInfo.FileType == "16" {
			continue
		}
		if fileInfo.GetFileSize() == "0 MB" {
			continue
		}
		logger.Info("已选择音频文件信息", "Quality", getQualityNamePairs(fileInfo.Quality), "FileSize", fileInfo.GetFileSize(), "FileType", fileInfo.MediaFileType)
		mp3file, code := getFile(contentCode, fileInfo.Quality)
		if code != 0 {
			continue
		}
		if mp3file.FileURL == "" {
			logger.Error("获取音频文件失败", "SongID", songInfo.ContentID, "Quality", getQualityNamePairs(fileInfo.Quality))
			continue
		}
		fileInfos[i].Url = mp3file.FileURL
		fileInfos[i].AlbumName = mp3file.SongInfo.AlbumNamePair.Name
		fileInfos[i].SongName = Title
		fileInfos[i].SubTitle = SubTitle
		fileInfos[i].ArtistName = contentExInfo.ArtistNames
		fileInfos[i].Lyrics = Lyrics
		fileInfos[i].ListName = songInfo.ListName
		fileInfos[i].IV = mp3file.Iv
		fileInfos[i].SecretKey = mp3file.SecretKey
		switch Conf.CoverSize {
		case "middle":
			fileInfos[i].Picture = songInfo.Picture.MiddleImgURL
		case "small":
			fileInfos[i].Picture = songInfo.Picture.SmallImgURL
		default:
			fileInfos[i].Picture = songInfo.Picture.BigImgURL
		}

	}
	if Conf.Download {
		DownloadFiles(fileInfos)
		return nil
	}
	return

}
func Download(contentCode string) (fileInfos []FileInfo) {
	mp3file, code := getFile(contentCode, "1")
	if code != 0 {
		return
	}
	fileInfos = GetInfos(mp3file.SongInfo)
	return
}

func GetSong(data []Info) []Info {
	r := Conf.Range
	if r == "all" || r == "" {
		return GetRange(data, r)
	}
	r = strings.ReplaceAll(r, "，", ",")
	if r == "0" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Title", "SubTitle"})
		for k, v := range data {
			table.Append([]string{fmt.Sprintf("%d", k+1), fmt.Sprintf("%v", v.ContentName), fmt.Sprintf("%v", v.SubTitle)})
		}
		table.SetRowLine(true)
		table.SetRowSeparator("-")
		table.SetCenterSeparator("*")
		table.SetColumnSeparator("|")
		table.SetBorder(true)
		table.Render()
		fmt.Print("请输入要下载的歌曲ID,多个ID用逗号分隔:")
		_, _ = fmt.Scan(&r)

	}

	return GetRange(data, r)
}
func GetRange[T any](data []T, rs string) []T {
	if rs == "all" || rs == "" {
		return data
	}
	rs = strings.ReplaceAll(rs, "，", ",")

	ranges := strings.Split(rs, ",")
	var ret []T
	var rangeRet []int
	for _, v := range ranges {
		if strings.Contains(v, "-") {
			startEnd := strings.Split(v, "-")
			start := 0
			end := 0
			_, _ = fmt.Sscanf(startEnd[0], "%d", &start)
			_, _ = fmt.Sscanf(startEnd[1], "%d", &end)
			if start > end {
				start, end = end, start
			}
			if end > len(data) {
				end = len(data) + 1
			}
			if start < 1 {
				start = 1
			}

			for i := start; i <= end; i++ {
				rangeRet = append(rangeRet, i)
			}
		} else {
			intV := 0
			_, _ = fmt.Sscanf(v, "%d", &intV)
			rangeRet = append(rangeRet, intV)
		}
	}
	uniqueMap := make(map[int]struct{})
	var uniqueRangeRet []int
	for _, val := range rangeRet {
		if _, exists := uniqueMap[val]; !exists {
			uniqueMap[val] = struct{}{}
			uniqueRangeRet = append(uniqueRangeRet, val)
		}
	}
	sort.Ints(uniqueRangeRet)
	for _, v := range uniqueRangeRet {
		ret = append(ret, data[v-1])
	}
	return ret
}
