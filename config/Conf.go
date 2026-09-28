package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gospider007/requests"
	"github.com/spf13/viper"
	"os"
)

var (
	Conf   = NewConfig()
	IsLoad = false
	Req, _ = requests.NewClient(nil)
	Format = Conf.SingleFormat
)

type Config struct {
	Authorization      string `yaml:"authorization"`
	ExpireTime         string `yaml:"expiretime"`
	AccessToken        string `yaml:"accesstoken"`
	FileName           string `yaml:"file_name" mapstructure:"file_name"`
	SingleFormat       string `yaml:"single_format" mapstructure:"single_format"`
	AlbumFormat        string `yaml:"album_format" mapstructure:"album_format"`
	PlaylistFormat     string `yaml:"playlist_format" mapstructure:"playlist_format"`
	ArtistSingleFormat string `yaml:"artist_single_format" mapstructure:"artist_single_format"`
	ArtistAlbumFormat  string `yaml:"artist_album_format" mapstructure:"artist_album_format"`
	Lyric              bool   `yaml:"lyric"`
	Cover              bool   `yaml:"cover"`
	CoverSize          string `yaml:"cover_size" mapstructure:"cover_size"`
	Output             string `yaml:"output"`
	Download           bool   `yaml:"download"`
	MaxCount           int    `yaml:"max_count" mapstructure:"max_count"`
	ArtistType         string `yaml:"artist_type" mapstructure:"artist_type"`
	Quality            string `yaml:"quality"`
	Range              string `yaml:"range"`
	AblumRange         string `yaml:"ablum_range" mapstructure:"ablum_range"`
	NumThreads         int    `yaml:"num_threads"`
}

func NewConfig() *Config {
	//
	return &Config{
		Lyric:              true,
		Cover:              true,
		CoverSize:          "big",
		Output:             "./o",
		SingleFormat:       "{output}/{artist}/",
		AlbumFormat:        "{output}/{album}/",
		PlaylistFormat:     "{output}/{playlist_name}/",
		ArtistSingleFormat: "{output}/{artist}/",
		ArtistAlbumFormat:  "{output}/{artist}/{album}/",
		FileName:           "{title}-{filesize}-{duration}-{rate}",
		Quality:            "best",
		Range:              "1-3",
		ArtistType:         "s",
		MaxCount:           100,
		AblumRange:         "1-3",
		Download:           false,
		NumThreads:         10,
	}
}
func UpdateConfig(key string, value interface{}) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("写入配置文件失败：%s\n", err)
	}
}
func InitConfig() {
	if IsLoad {
		return
	}
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取目录失败：%s\n", err))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(workDir + "/config")
	err = viper.ReadInConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已更新")
		err := viper.Unmarshal(Conf)
		if err != nil {
			fmt.Printf("解析配置文件失败：%s\n", err)
		}
	})
	viper.WatchConfig()
	if err != nil {
		fmt.Printf("读取配置文件失败：%s\n", err)
		err = viper.SafeWriteConfigAs("config.yml")
		if err != nil {
			fmt.Printf("写入配置文件失败：%s\n", err)
		}
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Printf("再次读取配置文件失败：%s\n", err)
			os.Exit(1)
		}
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		fmt.Printf("解析配置文件失败：%s\n", err)
		os.Exit(1)
	}
	IsLoad = true
}
func init() {
	InitConfig()
}
