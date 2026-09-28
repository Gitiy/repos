package auth

import (
	"fmt"
	"github.com/gospider007/requests"
	"hwyy/config"
	"hwyy/log"
	"math/rand"
	"strings"
	"time"
)

type LoginEx struct {
	Result struct {
		ResultCode    string `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
	} `json:"result"`
	Gender        string `json:"gender"`
	SnsUserID     string `json:"snsUserID"`
	RealNameLevel string `json:"realNameLevel"`
	UserID        string `json:"userID"`
	UserHC        string `json:"userHC"`
	AgeLevel      string `json:"ageLevel"`
	Token         string `json:"token"`
	ExpireTime    string `json:"expireTime"`
	RealNameAuth  string `json:"realNameAuth"`
}

var (
	logger = log.Logger
	Conf   = config.Conf
	req    = config.Req
)

// getUTCTime
func getUTCTime() string {
	now := time.Now()
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ",
		now.Year(), int(now.Month()), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

// createRandomString
func createRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var builder strings.Builder
	rand.Seed(time.Now().UnixNano()) // 设置随机种子

	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset)) // 随机选择字符集中的一个字符
		builder.WriteByte(charset[index])
	}

	return builder.String()
}

// getAuthorization
func getAuthorization() string {
	//Digest username=app-music-baseline, nonce=1bL69F4549l8J2z0q1vl6Wc608lAH879, created=2024-10-25T16:51:34Z
	u := "https://api-drcn.music.dbankcloud.cn/music-odpapp-service/v5/service/userinfo/loginex"
	data := map[string]interface{}{
		"accessToken": Conf.AccessToken,
	}
	Authorization := fmt.Sprintf("Digest username=%s, nonce=%s, created=%s", "app-music-baseline", createRandomString(32), getUTCTime())
	headers := map[string]string{
		"User-Agent":    "model=MRO-W00,brand=huawei,rom=,emui=,os=12,apilevel=32,manufacturer=huawei,useBrandCust=0,extChannel=,cpucore=4,memory=6.0G,srceenHeight=2160,screenWidth=3840,harmonyApiLevel=,huaweiOsBrand=",
		"authorization": Authorization,
	}
	resp, err := req.Post(nil, u, requests.RequestOption{
		Headers: headers,
		Json:    data,
	})
	if err != nil {
		logger.Error("获取Authorization失败", "错误信息", err)
	}
	var loginEx LoginEx
	if _, err := resp.Json(&loginEx); err != nil {
		logger.Error("解析Authorization失败", "错误信息", err)
	}
	if loginEx.Result.ResultCode != "000000" {
		logger.Error("获取Authorization失败", "错误信息", loginEx.Result.ResultMessage)
		return ""
	}
	Authorization = fmt.Sprintf("%s, token=%s", Authorization, loginEx.Token)
	ExpireTime := loginEx.ExpireTime
	config.UpdateConfig("authorization", Authorization)
	config.UpdateConfig("expire_time", ExpireTime)
	return ""
}

// init
func init() {
	//如果有AccessToken则获取Authorization
	if Conf.AccessToken == "" {
		return
	}
	expiryTime, err := time.Parse("20060102150405", Conf.ExpireTime)
	if err != nil {
		getAuthorization()
		return
	}
	if Conf.ExpireTime == "0" || expiryTime.Unix() < time.Now().Unix() {
		getAuthorization()
	}

}
