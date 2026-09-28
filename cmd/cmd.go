package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"hwyy/config"
	"os"
)

var (
	Long = ` 
                    本程序仅供学习使用，请勿用于商业用途，下载的音乐版权归华为音乐所有。请在下载后24小时内删除。
▄         ▄  ▄         ▄  ▄         ▄  ▄         ▄       ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄            ▄▄▄▄▄▄▄▄▄▄▄ 
▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌     ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌          ▐░░░░░░░░░░░▌
▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌      ▀▀▀▀█░█▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌▐░▌          ▐░█▀▀▀▀▀▀▀▀▀ 
▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌          ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░▌          
▐░█▄▄▄▄▄▄▄█░▌▐░▌   ▄   ▐░▌▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌          ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░█▄▄▄▄▄▄▄▄▄ 
▐░░░░░░░░░░░▌▐░▌  ▐░▌  ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌          ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░░░░░░░░░░░▌
▐░█▀▀▀▀▀▀▀█░▌▐░▌ ▐░▌░▌ ▐░▌ ▀▀▀▀█░█▀▀▀▀  ▀▀▀▀█░█▀▀▀▀           ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌           ▀▀▀▀▀▀▀▀▀█░▌
▐░▌       ▐░▌▐░▌▐░▌ ▐░▌▐░▌     ▐░▌          ▐░▌               ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌                    ▐░▌
▐░▌       ▐░▌▐░▌░▌   ▐░▐░▌     ▐░▌          ▐░▌               ▐░▌     ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄█░▌
▐░▌       ▐░▌▐░░▌     ▐░░▌     ▐░▌          ▐░▌               ▐░▌     ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀         ▀  ▀▀       ▀▀       ▀            ▀                 ▀       ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀ 
                                                                                                                          `
)
var rootCmd = &cobra.Command{
	Use:   "hwyy [flags] url",
	Short: "华为音乐下载",
	Long:  Long,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Long)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Conf.Lyric, "lyric", "l", config.Conf.Lyric, "Download lyrics")
	rootCmd.PersistentFlags().BoolVarP(&config.Conf.Cover, "cover", "c", config.Conf.Cover, "Download cover image")
	rootCmd.PersistentFlags().StringVarP(&config.Conf.CoverSize, "cover_size", "s", config.Conf.CoverSize, "Cover size: big (1000x1000), mid (600x600), small (320x320)")
	rootCmd.PersistentFlags().StringVarP(&config.Conf.Output, "output", "o", config.Conf.Output, "Output path")
	rootCmd.PersistentFlags().StringVarP(&config.Conf.Quality, "quality", "q", config.Conf.Quality, "Quality: 1 (standard), 2 (HQ), 3 (SQ), 5 (Hi-Res), 13 (Audio Vivid), 15 (multi-track), all, best")
	rootCmd.PersistentFlags().IntVarP(&config.Conf.MaxCount, "max_count", "m", config.Conf.MaxCount, "Maximum number of songs to parse at once")
	rootCmd.PersistentFlags().StringVarP(&config.Conf.ArtistType, "artist_type", "a", config.Conf.ArtistType, "Artist type: s (single), a (album)")
	rootCmd.PersistentFlags().StringVarP(&config.Conf.Range, "range", "r", config.Conf.Range, "Range: 1-10,13,20-30, all")
	rootCmd.PersistentFlags().StringVarP(&config.Conf.AblumRange, "ablum_range", "z", config.Conf.AblumRange, "Ablum range: 1-10,13,20-30, all")
	rootCmd.PersistentFlags().BoolVarP(&config.Conf.Download, "download", "d", config.Conf.Download, "Download songs")
	rootCmd.PersistentFlags().BoolVarP(&config.Conf.Download, "num_threads", "t", config.Conf.Download, "Download Threadss")

}
