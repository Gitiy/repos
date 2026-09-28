package Song

import (
	"fmt"
	"hwyy/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func DownloadFiles(files []FileInfo) {

	for _, file := range files {
		if file.Url == "" {
			continue
		}
		artist := NormalizeFilename(file.ArtistName)
		title := NormalizeFilename(file.SongName)
		album := NormalizeFilename(file.AlbumName)
		subTitle := NormalizeFilename(file.SubTitle)
		output := Conf.Output
		playlistName := NormalizeFilename(file.ListName)
		fileFormat := strings.ToLower(file.MediaFileType)
		filesize := file.GetFileSize()
		duration := file.GetDuration()
		rate := file.SampleRate
		BitRate := file.BitRate
		ss := map[string]string{
			"artist":        artist,
			"title":         title,
			"album":         album,
			"sub_title":     subTitle,
			"output":        output,
			"playlist_name": playlistName,
			"fileFormat":    fileFormat,
			"filesize":      filesize,
			"duration":      duration,
			"rate":          BitRate + "-" + rate,
		}
		outPath := config.Format
		for k, v := range ss {
			outPath = strings.ReplaceAll(outPath, "{"+k+"}", v)
		}
		filename := Conf.FileName
		for k, v := range ss {
			filename = strings.ReplaceAll(filename, "{"+k+"}", v)
		}
		filename = NormalizeFilename(filename)
		logger.Debug("下载文件", "文件名", filename, "路径", outPath)
		if Conf.Lyric {
			lyric := file.Lyrics
			if lyric != "" {
				lyricPath := outPath + title + ".lrc"
				logger.Info("下载歌词", "路径", lyricPath)
				if err := DownloadFile(lyric, outPath, title+".lrc"); err != nil {
					logger.Error("下载歌词失败", "err", err)
				} else {
					logger.Info("下载歌词成功", "路径", lyricPath)
				}

			}
		}
		if Conf.Cover {
			cover := file.Picture
			if cover != "" {
				coverPath := outPath + filename + ".jpg"
				logger.Info("下载封面", "路径", coverPath)
				if err := DownloadFile(cover, outPath, album+".jpg"); err != nil {
					logger.Error("下载封面失败", "err", err)
				} else {
					logger.Info("下载封面成功", "路径", coverPath)
				}
			}
		}
		logger.Debug("下载音频中", "路径", outPath)
		if err := DownloadFile(file.Url, outPath, filename+"."+fileFormat); err != nil {
			logger.Error("下载文件失败", "err", err)
		}
		logger.Debug("下载音频成功", "路径", outPath)
		iv := file.IV
		SecretKey := file.SecretKey
		drm := file.Drm
		if iv != "" && SecretKey != "" {
			switch drm {
			case "11":
				logger.Info("检测到音频为drm11，开始解密")
				err := Decrypt11(outPath+filename+"."+fileFormat, iv, SecretKey)
				if err != nil {
					logger.Error("解密失败", "err", err)
					continue
				}
				logger.Info("解密成功", "文件名", filename)
			default:
				logger.Error("未知drm类型", "drm", drm)
			}
		}

	}
}
func DownloadFile(url string, path string, filename string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
	}

	filePath := filepath.Join(path, filename)
	tmpFilePath := filePath + ".tmp"

	if _, err := os.Stat(filePath); err == nil {
		logger.Debug("文件已存在，跳过下载", "filename", filename)
		return nil
	}

	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusMethodNotAllowed {
			return singleThreadDownload(url, tmpFilePath, filePath)
		} else if resp.StatusCode == http.StatusForbidden { //可能不支持head请求 ，少部分会出现
			return singleThreadDownload(url, tmpFilePath, filePath)
		}
		return fmt.Errorf("无法获取文件信息: %s", resp.Status)
	}

	if resp.Header.Get("Accept-Ranges") != "bytes" {
		fmt.Println("服务器不支持分块下载，使用单线程下载。")
		return singleThreadDownload(url, tmpFilePath, filePath)
	}

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return fmt.Errorf("文件大小无效: %v", err)
	}

	numThreads := Conf.NumThreads
	return multiThreadDownload(url, tmpFilePath, filePath, contentLength, numThreads)
}

// singleThreadDownload 执行单线程下载，如果服务器不支持分块下载
func singleThreadDownload(url, tmpFilePath, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("文件下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("文件下载失败: %s", resp.Status)
	}

	out, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return os.Rename(tmpFilePath, filePath)
}

// multiThreadDownload 使用多线程下载文件
func multiThreadDownload(url, tmpFilePath, filePath string, contentLength, numThreads int) error {
	partSize := contentLength / numThreads
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * partSize
			end := start + partSize - 1
			if i == numThreads-1 {
				end = contentLength - 1
			}

			err := downloadPart(url, tmpFilePath+fmt.Sprintf(".part%d", i), start, end)
			if err != nil {
				fmt.Printf("下载文件的第 %d 部分失败: %v\n", i, err)
			}
		}(i)
	}

	wg.Wait()

	return combineParts(tmpFilePath, filePath, numThreads)
}

// downloadPart 下载文件的指定字节范围，并保存到临时文件
func downloadPart(url, partPath string, start, end int) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("服务器不支持部分内容")
	}

	out, err := os.Create(partPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// combineParts 合并所有部分文件到最终文件，并删除临时部分文件
func combineParts(tmpFilePath, filePath string, numParts int) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	for i := 0; i < numParts; i++ {
		partPath := tmpFilePath + fmt.Sprintf(".part%d", i)
		partFile, err := os.Open(partPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, partFile)
		partFile.Close()
		if err != nil {
			return err
		}

		if err := os.Remove(partPath); err != nil {
			return err
		}
	}

	return nil
}

func NormalizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, "/", " ")
	filename = strings.ReplaceAll(filename, "&", " and ")
	filename = strings.ReplaceAll(filename, "$", "S")
	filename = regexp.MustCompile(`[;:\s]`).ReplaceAllString(filename, " ")
	filename = regexp.MustCompile(`[\\*!?¿,'\"()<>|#]`).ReplaceAllString(filename, " ")
	filename = strings.TrimRight(regexp.MustCompile(`[.]{2,}`).ReplaceAllString(filename, "."), ".")
	filename = strings.TrimRight(regexp.MustCompile(`[ ]{2,}`).ReplaceAllString(filename, " "), "")
	return filename
}
