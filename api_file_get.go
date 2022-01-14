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
	"net/http"
)

// GetFile 获取文件信息
func (r *FileService) GetFile(ctx context.Context, request *GetFileReq) (*GetFileResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "GetFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/file/get",
		Body:   request,
	}
	resp := new(getFileResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetFileResp, nil
}

type GetFileReq struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type GetFileResp struct {
	File
	Trashed bool `json:"trashed"`
}

type getFileResp struct {
	Message string `json:"message"`
	GetFileResp
}
