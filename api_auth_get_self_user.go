package aliyundrive

import (
	"context"
	"net/http"
)

func (r *AuthService) GetSelfUser(ctx context.Context) (*GetSelfUserResp, error) {
	req := &RawRequestReq{
		Scope:   "Auth",
		API:     "GetSelfUser",
		Method:  http.MethodPost,
		URL:     "https://api.aliyundrive.com/v2/user/get",
		Body:    struct{}{},
		IsFile:  false,
		headers: nil,
	}
	resp := new(getSelfUserResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetSelfUserResp, nil
}

type getSelfUserResp struct {
	Message string `json:"message"`
	GetSelfUserResp
}

type GetSelfUserResp struct {
	DomainID                    string      `json:"domain_id"`
	UserID                      string      `json:"user_id"`
	Avatar                      string      `json:"avatar"`
	CreatedAt                   int64       `json:"created_at"`
	UpdatedAt                   int64       `json:"updated_at"`
	Email                       string      `json:"email"`
	NickName                    string      `json:"nick_name"`
	Phone                       string      `json:"phone"`
	Role                        string      `json:"role"`
	Status                      string      `json:"status"`
	UserName                    string      `json:"user_name"`
	Description                 string      `json:"description"`
	DefaultDriveID              string      `json:"default_drive_id"`
	DenyChangePasswordBySelf    bool        `json:"deny_change_password_by_self"`
	NeedChangePasswordNextLogin bool        `json:"need_change_password_next_login"`
	Permission                  interface{} `json:"permission"`
}
