package aliyundrive

import (
	"context"
	"net/http"
)

func (r *FileService) GetSBox(ctx context.Context) (*GetSBoxResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "GetSBox",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/sbox/get",
		Body:   struct{}{},
	}
	resp := new(getSBoxResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetSBoxResp, nil
}

type GetSBoxResp struct {
	DriveID          string `json:"drive_id"`
	SboxUsedSize     int    `json:"sbox_used_size"`
	SboxRealUsedSize int    `json:"sbox_real_used_size"`
	SboxTotalSize    int64  `json:"sbox_total_size"`
	RecommendVip     string `json:"recommend_vip"`
	PinSetup         bool   `json:"pin_setup"`
	Locked           bool   `json:"locked"`
	InsuranceEnabled bool   `json:"insurance_enabled"`
}

type getSBoxResp struct {
	Message string `json:"message"`
	GetSBoxResp
}
