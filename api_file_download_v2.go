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
	"net/http"
)

// GetFileDownloadURLV2 使用 DownloadFile API 获取下载链接（返回 302 跳转）
func (r *FileService) GetFileDownloadURLV2(ctx context.Context, request *GetFileDownloadURLReq) (string, error) {
	url := fmt.Sprintf("https://api.aliyundrive.com/v2/file/download?drive_id=%s&file_id=%s",
		request.DriveID, request.FileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	token, _ := r.cli.store.Get(ctx, "")
	if token != nil && token.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.aliyundrive.com/")

	// 禁用自动重定向，获取 302 跳转 URL
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusTemporaryRedirect {
		location := resp.Header.Get("Location")
		if location != "" {
			return location, nil
		}
	}

	return "", fmt.Errorf("无法获取下载链接，状态码: %d", resp.StatusCode)
}
