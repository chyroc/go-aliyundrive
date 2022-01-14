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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chyroc/go-aliyundrive/internal/helper_qrcode"
)

func IsTokenExpired(err error) bool {
	return err != nil && strings.Contains(err.Error(), "AccessTokenExpired")
}

func (r *AuthService) LoginByQrcode(ctx context.Context) (*GetSelfUserResp, error) {
	userInfo, err := r.GetSelfUser(ctx)
	if IsTokenExpired(err) {
		token, err := r.cli.store.Get(ctx, "")
		if err != nil {
			r.cli.log(ctx, LogLevelError, "get user failed, then get store failed: %s", err)
		} else if token.RefreshToken != "" {
			refreshTokenResp, err := r.RefreshToken(ctx, &RefreshTokenReq{RefreshToken: token.RefreshToken})
			if err != nil {
				r.cli.log(ctx, LogLevelError, "get user failed, then refresh token failed: %s", err)
			} else if refreshTokenResp.RefreshToken != "" {
				if err = r.cli.store.Set(ctx, refreshTokenResp.Token()); err != nil {
					r.cli.log(ctx, LogLevelError, "get user failed, then refresh token failed: %s", err)
				}

				userInfo, err = r.GetSelfUser(ctx)
				if err != nil {
					r.cli.log(ctx, LogLevelError, "get user twice failed when token expired: %s", err)
				}
			}
		}
	}
	if userInfo != nil && userInfo.UserID != "" {
		return userInfo, nil
	}
	err = r.cli.store.Set(ctx, nil)
	if err != nil {
		return nil, err
	}
	if err := r.internalLoginByQrcode(ctx); err != nil {
		return nil, err
	}

	return r.GetSelfUser(ctx)
}

func (r *AuthService) internalLoginByQrcode(ctx context.Context) error {
	if err := r.preLogin(ctx); err != nil {
		return err
	}

	qrcode, err := r.getQrCode(ctx)
	if err != nil {
		return err
	}

	err = helper_qrcode.New(true).Print(qrcode.CodeContent, helper_qrcode.Low)
	if err != nil {
		return err
	}
	fmt.Println("请用阿里云盘 App 扫码")

	scaned := false
	for {
		// NEW / SCANED / EXPIRED / CANCELED / CONFIRMED
		res, err := r.queryQrCode(ctx, strconv.FormatInt(qrcode.T, 10), qrcode.Ck)
		if err != nil {
			return err
		}

		switch res.QrCodeStatus {
		case "NEW":

		case "SCANED":
			if !scaned {
				fmt.Println("扫描成功, 请在手机上根据提示确认登录")
			}
			scaned = true
		case "EXPIRED":
			return fmt.Errorf("二维码过期")
		case "CANCELED":
			return fmt.Errorf("取消登录")
		case "CONFIRMED":
			biz := res.BizAction.PdsLoginResult
			if err := r.confirmLogin(ctx, biz.AccessToken); err != nil {
				return err
			}
			if err := r.cli.store.Set(ctx, &Token{AccessToken: biz.AccessToken, RefreshToken: biz.RefreshToken, ExpiredAt: time.Now().Add(time.Second * time.Duration(biz.ExpiresIn))}); err != nil {
				return err
			}
			return nil
		default:
			panic(res)
		}

		time.Sleep(time.Second)
	}
}

func (r *AuthService) preLogin(ctx context.Context) error {
	req := &RawRequestReq{
		Scope:  "Auth",
		API:    "LoginByQrCode",
		Method: http.MethodGet,
		URL: "https://auth.aliyundrive.com/v2/oauth/authorize?client_id=25dzX3vbYqktVxyX" +
			"&redirect_uri=https%3A%2F%2Fwww.aliyundrive.com%2Fsign%2Fcallback" +
			"&response_type=code" +
			"&login_type=custom" +
			"&state=%7B%22origin%22%3A%22https%3A%2F%2Fwww.aliyundrive.com%22%7D",
	}
	_, err := r.cli.RawRequest(ctx, req, nil)
	return err
}

func (r *AuthService) getQrCode(ctx context.Context) (*getQrCodeResp, error) {
	req := &RawRequestReq{
		Scope:  "Auth",
		API:    "getQrCode",
		Method: http.MethodGet,
		URL: "https://passport.aliyundrive.com/newlogin/qrcode/generate.do?" +
			"appName=aliyun_drive" +
			"&fromSite=52" +
			"&appName=aliyun_drive" +
			"&appEntrance=web" +
			"&isMobile=false" +
			"&lang=zh_CN" +
			"&returnUrl=" +
			"&fromSite=52" +
			"&bizParams=" +
			"&_bx-v=2.0.31",
	}
	resp := struct {
		Content struct {
			Data getQrCodeResp `json:"data"`
		} `json:"content"`
	}{}

	_, err := r.cli.RawRequest(ctx, req, &resp)
	if err != nil {
		return nil, err
	} else if resp.Content.Data.TitleMsg != "" {
		return nil, fmt.Errorf(resp.Content.Data.TitleMsg)
	}

	return &resp.Content.Data, nil
}

type getQrCodeResp struct {
	TitleMsg    string `json:"title_msg"`
	T           int64  `json:"t"`
	CodeContent string `json:"codeContent"`
	Ck          string `json:"ck"`
	ResultCode  int    `json:"resultCode"`
}

func (r *AuthService) queryQrCode(ctx context.Context, t, ck string) (*queryQrCodeResp, error) {
	req := &RawRequestReq{
		Scope:  "Auth",
		API:    "queryQrCode",
		Method: http.MethodPost,
		URL:    "https://passport.aliyundrive.com/newlogin/qrcode/query.do?appName=aliyun_drive&fromSite=52&_bx-v=2.0.31",
		Body: queryQrCodeReq{
			T:           t,
			Ck:          ck,
			AppName:     "aliyun_drive",
			AppEntrance: "web",
			IsMobile:    "false",
			Lang:        "zh_CN",
			ReturnURL:   "",
			FromSite:    "52",
			BizParams:   "",
			Navlanguage: "zh-CN",
			NavPlatform: "MacIntel",
		},
		IsFile: false,
		headers: map[string]string{
			"Accept":          "application/json, text/plain",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		},
	}
	resp := struct {
		Content struct {
			Data queryQrCodeResp `json:"data"`
		} `json:"content"`
	}{}

	_, err := r.cli.RawRequest(ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Content.Data.BizExt != "" {
		bs, _ := base64.StdEncoding.DecodeString(resp.Content.Data.BizExt)
		_ = json.Unmarshal(bs, &resp.Content.Data.BizAction)
	}

	return &resp.Content.Data, nil
}

type queryQrCodeReq struct {
	T           string `form-url-encoded:"t"`
	Ck          string `form-url-encoded:"ck"`
	AppName     string `form-url-encoded:"appName"`
	AppEntrance string `form-url-encoded:"appEntrance"`
	IsMobile    string `form-url-encoded:"isMobile"`
	Lang        string `form-url-encoded:"lang"`
	ReturnURL   string `form-url-encoded:"returnUrl"`
	FromSite    string `form-url-encoded:"fromSite"`
	BizParams   string `form-url-encoded:"bizParams"`
	Navlanguage string `form-url-encoded:"navlanguage"`
	NavPlatform string `form-url-encoded:"navPlatform"`
}

type queryQrCodeResp struct {
	QrCodeStatus string `json:"qrCodeStatus"`
	ResultCode   int    `json:"resultCode"`

	LoginResult          string               `json:"loginResult"`
	LoginSucResultAction string               `json:"loginSucResultAction"`
	BizAction            queryQrCodeBizAction `json:"-"`
	St                   string               `json:"st"`
	LoginType            string               `json:"loginType"`
	BizExt               string               `json:"bizExt"`
	LoginScene           string               `json:"loginScene"`
	AppEntrance          string               `json:"appEntrance"`
	Smartlock            bool                 `json:"smartlock"`
}

type queryQrCodeBizAction struct {
	PdsLoginResult struct {
		Role           string        `json:"role"`
		IsFirstLogin   bool          `json:"isFirstLogin"`
		NeedLink       bool          `json:"needLink"`
		LoginType      string        `json:"loginType"`
		NickName       string        `json:"nickName"`
		NeedRpVerify   bool          `json:"needRpVerify"`
		Avatar         string        `json:"avatar"`
		AccessToken    string        `json:"accessToken"`
		UserName       string        `json:"userName"`
		UserID         string        `json:"userId"`
		DefaultDriveID string        `json:"defaultDriveId"`
		ExistLink      []interface{} `json:"existLink"`
		ExpiresIn      int           `json:"expiresIn"`
		ExpireTime     time.Time     `json:"expireTime"`
		RequestID      string        `json:"requestId"`
		DataPinSetup   bool          `json:"dataPinSetup"`
		State          string        `json:"state"`
		TokenType      string        `json:"tokenType"`
		DataPinSaved   bool          `json:"dataPinSaved"`
		RefreshToken   string        `json:"refreshToken"`
		Status         string        `json:"status"`
	} `json:"pds_login_result"`
}

func (r *AuthService) confirmLogin(ctx context.Context, accessToken string) error {
	gotoURL := ""
	{
		req := &RawRequestReq{
			Scope:  "Auth",
			API:    "confirmLogin",
			Method: http.MethodPost,
			URL:    "https://auth.aliyundrive.com/v2/oauth/token_login",
			Body: struct {
				Token string `json:"token"`
			}{
				Token: accessToken,
			},
			IsFile:  false,
			headers: nil,
		}
		resp := new(confirmLoginResp)

		if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
			return err
		}
		gotoURL = resp.Goto
	}

	{
		req := &RawRequestReq{
			Scope:  "Auth",
			API:    "confirmLogin",
			Method: http.MethodGet,
			URL:    gotoURL,
		}

		if _, err := r.cli.RawRequest(ctx, req, nil); err != nil {
			return err
		}
	}
	return nil
}

type confirmLoginResp struct {
	Goto string `json:"goto"`
}
