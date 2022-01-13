package aliyundrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/chyroc/go-aliyundrive/internal/helper_tool"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

// GetFile 获取文件信息
func (r *FileService) DownloadFile(ctx context.Context, request *DownloadFileReq) error {
	if request.Dist == "" && request.DistDir == "" {
		return fmt.Errorf("must set Dist or DistDir")
	}

	distName := request.Dist
	if request.Dist == "" && request.DistDir != "" {
		res, err := r.GetFile(ctx, &GetFileReq{
			DriveID: request.DriveID,
			FileID:  request.FileID,
		})
		if err != nil {
			return err
		}
		distName = strings.TrimRight(request.DistDir, "/") + "/" + res.Name
	}

	if helper_tool.IsFileExist(distName) {
		if request.ConflictType == DownloadFileConflictTypeError {
			return fmt.Errorf("文件 %q 已存在，无法下载", distName)
		}
		if request.ConflictType == DownloadFileConflictTypeAutoRename {
			distName = helper_tool.AutoRenameFile(distName)
		}
	}

	res, err := r.GetFileDownloadURL(ctx, &GetFileDownloadURLReq{
		DriveID: request.DriveID,
		FileID:  request.FileID,
	})
	if err != nil {
		return err
	}

	err = downloadURL(res.URL, distName, request.ShowProgressBar)
	if err != nil {
		return err
	}
	return nil
}

type DownloadFileReq struct {
	DriveID         string                   `json:"drive_id"`
	FileID          string                   `json:"file_id"`
	Dist            string                   `json:"dist"`              // 如果此值不为空，则下载文件到这个位置
	DistDir         string                   `json:"dist_dir"`          // 如果此值不为空，则下载文件到 `DistDir`/<name> 位置
	ConflictType    DownloadFileConflictType `json:"conflict_type"`     // 如果目标文件已存在，处理方式：报错，覆盖，自动重命名，默认是自动重命名
	ShowProgressBar bool                     `json:"show_progress_bar"` // 展示下载进度条
}

type DownloadFileConflictType int

const (
	DownloadFileConflictTypeAutoRename DownloadFileConflictType = 0
	DownloadFileConflictTypeOverwrite  DownloadFileConflictType = 1
	DownloadFileConflictTypeError      DownloadFileConflictType = 2
)

var downloadHttpClient = http.Client{}

func downloadURL(url string, filename string, showProgressBar bool) error {
	var deleteTemp = true
	var tmp = filename + ".tmp"
	defer func() {
		// 任何的异常退出都会导致临时文件被删除
		if deleteTemp {
			os.Remove(tmp)
		}
	}()
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	f.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.aliyundrive.com/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if showProgressBar {
		fileSize := resp.ContentLength
		progress := mpb.New(mpb.WithWidth(20))
		bar := progress.AddBar(
			fileSize,
			// 进度条前的修饰
			mpb.PrependDecorators(
				decor.Name("[download] "),
				decor.CountersKibiByte("% .2f / % .2f"), // 已下载数量
				decor.Percentage(decor.WCSyncSpace),     // 进度百分比
			),
			// 进度条后的修饰
			mpb.AppendDecorators(
				decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
			),
		)
		reader := bar.ProxyReader(resp.Body)
		defer reader.Close()

		if _, err := io.Copy(f, reader); err != nil {
			return err
		}
	} else {
		if _, err := io.Copy(f, resp.Body); err != nil {
			return err
		}
	}
	if err := os.Rename(tmp, filename); err != nil {
		return err
	}
	deleteTemp = false

	return nil
}
