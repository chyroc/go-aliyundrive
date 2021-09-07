package aliyundrive

import (
	"context"
	"net/http"
)

// GetFilePath 重命名
func (r *FileService) GetFilePath(ctx context.Context, request *GetFilePathReq) (*GetFilePathResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "GetFilePath",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/adrive/v1/file/get_path",
		Body:   request,
	}
	resp := new(getFilePathResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetFilePathResp, nil
}

type GetFilePathReq struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type GetFilePathResp struct {
	Items []*File `json:"items"`
}

type getFilePathResp struct {
	Message string `json:"message"`
	GetFilePathResp
}
