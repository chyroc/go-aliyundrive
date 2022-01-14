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

// CreateFolder 创建文件夹
func (r *FileService) CreateFolder(ctx context.Context, request *CreateFolderReq) (*CreateFolderResp, error) {
	request.CheckNameMode = "refuse"
	request.Type = "folder"

	req := &RawRequestReq{
		Scope:  "File",
		API:    "CreateFolder",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/adrive/v2/file/createWithFolders",
		Body:   request,
	}
	resp := new(createFolderResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.CreateFolderResp, nil
}

type CreateFolderReq struct {
	DriveID       string   `json:"drive_id"`
	ParentFileID  string   `json:"parent_file_id"`
	Name          string   `json:"name"`
	CheckNameMode string   `json:"check_name_mode"`
	Type          FileType `json:"type"`
}

type CreateFolderResp struct {
	ParentFileID string   `json:"parent_file_id"`
	Type         FileType `json:"type"`
	FileID       string   `json:"file_id"`
	DomainID     string   `json:"domain_id"`
	DriveID      string   `json:"drive_id"`
	FileName     string   `json:"file_name"`
	EncryptMode  string   `json:"encrypt_mode"`
}

type createFolderResp struct {
	Message string `json:"message"`
	CreateFolderResp
}
