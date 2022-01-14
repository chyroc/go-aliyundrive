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

// GetFile 获取文件信息
func (r *FileService) SearchFile(ctx context.Context, request *SearchFileReq) (*GetFileListResp, error) {
	//  "drive_id": "20167741",
	//  "query": "name match \"/文件夹2\" and type = \"folder\"",
	request.Limit = 100
	request.ImageThumbnailProcess = "image/resize,w_200/format,jpeg"
	request.ImageURLProcess = "image/resize,w_1920/format,jpeg"
	request.VideoThumbnailProcess = "video/snapshot,t_0,f_jpg,ar_auto,w_300"
	request.OrderBy = "updated_at DESC"
	req := &RawRequestReq{
		Scope:  "File",
		API:    "SearchFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/adrive/v3/file/search",
		Body:   request,
	}
	resp := new(getFileListResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	if resp.Message != "" {
		return nil, fmt.Errorf(resp.Message)
	}
	return &resp.GetFileListResp, nil
}

type SearchFileReq struct {
	DriveID               string `json:"drive_id"`
	Limit                 int    `json:"limit"`
	Query                 string `json:"query"`
	ImageThumbnailProcess string `json:"image_thumbnail_process"`
	ImageURLProcess       string `json:"image_url_process"`
	VideoThumbnailProcess string `json:"video_thumbnail_process"`
	OrderBy               string `json:"order_by"`
}
