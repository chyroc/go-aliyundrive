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

// GetShareByAnonymous 通过分享 ID 获取分享内容
func (r *ShareLinkService) GetShareByAnonymous(ctx context.Context, request *GetShareByAnonymousReq) (*GetShareByAnonymousResp, error) {
	req := &RawRequestReq{
		Scope:   "ShareLink",
		API:     "GetShareByAnonymous",
		Method:  http.MethodPost,
		URL:     "https://api.aliyundrive.com/adrive/v3/share_link/get_share_by_anonymous",
		Body:    request,
		IsFile:  false,
		headers: nil,
	}
	resp := new(getShareByAnonymousResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetShareByAnonymousResp, nil
}

type GetShareByAnonymousReq struct {
	ShareID string `json:"share_id"`
}

type GetShareByAnonymousResp struct {
	CreatorID    string    `json:"creator_id"`
	CreatorName  string    `json:"creator_name"`
	CreatorPhone string    `json:"creator_phone"`
	Expiration   string    `json:"expiration"`
	UpdatedAt    time.Time `json:"updated_at"`
	Vip          string    `json:"vip"`
	Avatar       string    `json:"avatar"`
	ShareName    string    `json:"share_name"`
	FileCount    int       `json:"file_count"`
	FileInfos    []struct {
		Type     string `json:"type"`
		FileID   string `json:"file_id"`
		FileName string `json:"file_name"`
	} `json:"file_infos"`
}

type getShareByAnonymousResp struct {
	Message string `json:"message"`
	GetShareByAnonymousResp
}
