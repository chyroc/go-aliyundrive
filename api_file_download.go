/**
 * Copyright 2022 chyroc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package aliyundrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/schollz/progressbar/v3"

	"github.com/chyroc/go-aliyundrive/internal/helper_tool"
	runewidth "github.com/mattn/go-runewidth"
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

	// 使用新的 DownloadFile API 获取下载链接（302 跳转）
	url, err := r.GetFileDownloadURLV2(ctx, &GetFileDownloadURLReq{
		DriveID: request.DriveID,
		FileID:  request.FileID,
	})
	if err != nil {
		return err
	}

	err = downloadURL(ctx, url, distName, request.ShowProgressBar)
	if err != nil {
		return err
	}
	return nil
}

// DownloadFileStream 获取文件流
func (r *FileService) DownloadFileStream(ctx context.Context, driveID, fileID string) (io.ReadCloser, error) {
	// 使用新的 DownloadFile API 获取下载链接（302 跳转）
	downloadURL, err := r.GetFileDownloadURLV2(ctx, &GetFileDownloadURLReq{DriveID: driveID, FileID: fileID})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.aliyundrive.com/")
	// 不设置 Accept-Encoding，让 Go 自动处理 gzip 解压缩
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect status code: %v", resp.Status)
	}
	return resp.Body, nil
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

func downloadURL(ctx context.Context, url string, filename string, showProgressBar bool) error {
	deleteTemp := true
	tmp := filename + ".tmp"
	defer func() {
		if deleteTemp {
			os.Remove(tmp)
		}
	}()

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建带超时的 HTTP 客户端，便于中断
	client := &http.Client{
		Timeout: 0, // 不设置总超时，但使用 context 控制
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.aliyundrive.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 使用通道和 goroutine 来支持中断
	type result struct {
		n   int64
		err error
	}

	done := make(chan result, 1)

	go func() {
		var n int64
		var err error

		if showProgressBar {
			bar := progressbar.NewOptions(
				int(resp.ContentLength),
				progressbar.OptionSetWriter(os.Stdout),
				progressbar.OptionSetDescription(runewidth.FillRight(path.Base(filename), 40)),
				progressbar.OptionShowBytes(true),
				progressbar.OptionSetPredictTime(true),
				progressbar.OptionFullWidth(),
				progressbar.OptionSetRenderBlankState(true),
				progressbar.OptionOnCompletion(func() {
					fmt.Printf("\n")
				}),
			)
			n, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
		} else {
			n, err = io.Copy(f, resp.Body)
		}

		done <- result{n: n, err: err}
	}()

	select {
	case <-ctx.Done():
		// context 被取消，关闭连接
		resp.Body.Close()
		return ctx.Err()
	case res := <-done:
		if res.err != nil {
			return res.err
		}
	}

	f.Close()
	if err := os.Rename(tmp, filename); err != nil {
		return err
	}
	deleteTemp = false

	return nil
}
