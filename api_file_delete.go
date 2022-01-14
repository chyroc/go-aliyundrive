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

// DeleteFile 删除文件
func (r *FileService) DeleteFile(ctx context.Context, request *DeleteFileReq) (*DeleteFileResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "DeleteFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/recyclebin/trash",
		Body:   request,
	}
	resp := new(deleteFileResp)

	if _, err := r.cli.RawRequest(ctx, req, nil); err != nil {
		return nil, err
	}
	return &resp.DeleteFileResp, nil
}

type DeleteFileReq struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type DeleteFileResp struct {
	DomainID    string `json:"domain_id"`
	DriveID     string `json:"drive_id"`
	FileID      string `json:"file_id"`
	AsyncTaskID string `json:"async_task_id"`
}

type deleteFileResp struct {
	Message string `json:"message"`
	DeleteFileResp
}
