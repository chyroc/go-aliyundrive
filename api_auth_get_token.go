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

func (r *AuthService) getToken(ctx context.Context, request *getTokenReq) (*RefreshTokenResp, error) {
	req := &RawRequestReq{
		Scope:  "Auth",
		API:    "GetToken",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/token/get",
		Body: &getTokenReq{
			Code:      request.Code,
			LoginType: "normal",
			DeviceId:  "aliyundrive",
		},
	}
	resp := new(RefreshTokenResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil // r.cli.token.Refresh(resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Second*time.Duration(resp.ExpiresIn)))
}

type getTokenReq struct {
	Code      string `json:"code"`
	LoginType string `json:"loginType"`
	DeviceId  string `json:"deviceId"`
}
