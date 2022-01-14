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

// RenameFile 重命名
func (r *FileService) RenameFile(ctx context.Context, request *RenameFileReq) (*RenameFileResp, error) {
	request.CheckNameMode = "refuse"

	req := &RawRequestReq{
		Scope:  "File",
		API:    "RenameFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v3/file/update",
		Body:   request,
	}
	resp := new(renameFileResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.RenameFileResp, nil
}

type RenameFileReq struct {
	DriveID       string `json:"drive_id"`
	FileID        string `json:"file_id"`
	Name          string `json:"name"`
	CheckNameMode string `json:"check_name_mode"`
}

type RenameFileResp struct {
	DriveID          string `json:"drive_id"`
	SboxUsedSize     int    `json:"sbox_used_size"`
	SboxRealUsedSize int    `json:"sbox_real_used_size"`
	SboxTotalSize    int64  `json:"sbox_total_size"`
	RecommendVip     string `json:"recommend_vip"`
	PinSetup         bool   `json:"pin_setup"`
	Locked           bool   `json:"locked"`
	InsuranceEnabled bool   `json:"insurance_enabled"`
}

type renameFileResp struct {
	Message string `json:"message"`
	RenameFileResp
}
