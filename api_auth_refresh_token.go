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

func (r *AuthService) RefreshToken(ctx context.Context, request *RefreshTokenReq) (*RefreshTokenResp, error) {
	req := &RawRequestReq{
		Scope:  "Auth",
		API:    "RefreshToken",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/token/refresh",
		Body:   request,
	}
	resp := new(RefreshTokenResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil // r.cli.token.Refresh(resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Second*time.Duration(resp.ExpiresIn)))
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResp struct {
	DefaultSboxDriveID string    `json:"default_sbox_drive_id"`
	Role               string    `json:"role"`
	DeviceID           string    `json:"device_id"`
	UserName           string    `json:"user_name"`
	NeedLink           bool      `json:"need_link"`
	ExpireTime         time.Time `json:"expire_time"`
	PinSetup           bool      `json:"pin_setup"`
	NeedRpVerify       bool      `json:"need_rp_verify"`
	Avatar             string    `json:"avatar"`
	UserData           struct {
		// DingDingRobotURL string `json:"DingDingRobotUrl"`
		// EncourageDesc    string `json:"EncourageDesc"`
		// FeedBackSwitch   bool   `json:"FeedBackSwitch"`
		// FollowingDesc    string `json:"FollowingDesc"`
		DingDingRobotURL string `json:"ding_ding_robot_url"`
		EncourageDesc    string `json:"encourage_desc"`
		FeedBackSwitch   bool   `json:"feed_back_switch"`
		FollowingDesc    string `json:"following_desc"`
	} `json:"user_data"`
	TokenType      string        `json:"token_type"`
	AccessToken    string        `json:"access_token"`
	DefaultDriveID string        `json:"default_drive_id"`
	RefreshToken   string        `json:"refresh_token"`
	IsFirstLogin   bool          `json:"is_first_login"`
	UserID         string        `json:"user_id"`
	NickName       string        `json:"nick_name"`
	ExistLink      []interface{} `json:"exist_link"`
	State          string        `json:"state"`
	ExpiresIn      int           `json:"expires_in"`
	Status         string        `json:"status"`
}

func (r *RefreshTokenResp) Token() *Token {
	return &Token{
		AccessToken:  r.AccessToken,
		ExpiredAt:    time.Now().Add(time.Second * time.Duration(r.ExpiresIn)),
		RefreshToken: r.RefreshToken,
	}
}
