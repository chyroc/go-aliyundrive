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
	"time"
)

func (r *FileService) GetFileList(ctx context.Context, request *GetFileListReq) (*GetFileListResp, error) {
	if !request.GetAll {
		return r.getFileList(ctx, request)
	}
	marker := ""
	items := []*File{}
	for {
		request.Marker = marker
		resp, err := r.getFileList(ctx, request)
		if err != nil {
			return nil, err
		}
		items = append(items, resp.Items...)
		marker = resp.NextMarker
		if marker == "" || len(resp.Items) < request.Limit {
			break
		}
	}

	return &GetFileListResp{Items: items, NextMarker: ""}, nil
}

func (r *FileService) getFileList(ctx context.Context, request *GetFileListReq) (*GetFileListResp, error) {
	{
		if request.ParentFileID == "" {
			request.ParentFileID = "root"
		}
		if request.Limit == 0 {
			request.Limit = 100
		}
		request.URLExpireSec = 1600
		request.ImageThumbnailProcess = "image/resize,w_400/format,jpeg"
		request.ImageURLProcess = "image/resize,w_1920/format,jpeg"
		request.VideoThumbnailProcess = "video/snapshot,t_0,f_jpg,ar_auto,w_300"
		request.Fields = "*"
		request.OrderBy = "updated_at"
		request.OrderDirection = "DESC"
	}

	req := &RawRequestReq{
		Scope:  "File",
		API:    "GetFileList",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/adrive/v3/file/list",
		Body:   request,
	}
	resp := new(getFileListResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetFileListResp, nil
}

type GetFileListReq struct {
	GetAll bool `json:"get_all"`

	ShareID               string `json:"share_id"` // drive_id 和 share_id 必选传其中一个
	DriveID               string `json:"drive_id"` // drive_id 和 share_id 必选传其中一个
	ParentFileID          string `json:"parent_file_id"`
	Marker                string `json:"marker"`
	Limit                 int    `json:"limit"`
	All                   bool   `json:"all"`
	URLExpireSec          int    `json:"url_expire_sec"`
	ImageThumbnailProcess string `json:"image_thumbnail_process"`
	ImageURLProcess       string `json:"image_url_process"`
	VideoThumbnailProcess string `json:"video_thumbnail_process"`
	Fields                string `json:"fields"`
	OrderBy               string `json:"order_by"`
	OrderDirection        string `json:"order_direction"`
}

type getFileListResp struct {
	Message string `json:"message"`
	GetFileListResp
}

type GetFileListResp struct {
	Items      []*File `json:"items"`
	NextMarker string  `json:"next_marker"`
}

type File struct {
	DriveID         string    `json:"drive_id"`
	DomainID        string    `json:"domain_id"`
	FileID          string    `json:"file_id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Hidden          bool      `json:"hidden"`
	Starred         bool      `json:"starred"`
	Status          string    `json:"status"`
	UserMeta        string    `json:"user_meta,omitempty"`
	ParentFileID    string    `json:"parent_file_id"`
	EncryptMode     string    `json:"encrypt_mode"`
	ContentType     string    `json:"content_type,omitempty"`
	FileExtension   string    `json:"file_extension,omitempty"`
	MimeType        string    `json:"mime_type,omitempty"`
	MimeExtension   string    `json:"mime_extension,omitempty"`
	Size            int64     `json:"size,omitempty"`
	Crc64Hash       string    `json:"crc64_hash,omitempty"`
	ContentHash     string    `json:"content_hash,omitempty"`
	ContentHashName string    `json:"content_hash_name,omitempty"`
	DownloadURL     string    `json:"download_url,omitempty"`
	URL             string    `json:"url,omitempty"`
	Thumbnail       string    `json:"thumbnail,omitempty"`
	Category        string    `json:"category,omitempty"`
	PunishFlag      int       `json:"punish_flag,omitempty"`
}
