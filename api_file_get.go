package aliyundrive

import (
	"context"
	"net/http"
)

// GetFile 获取文件信息
func (r *FileService) GetFile(ctx context.Context, request *GetFileReq) (*GetFileResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "GetFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/file/get",
		Body:   request,
	}
	resp := new(getFileResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetFileResp, nil
}

type GetFileReq struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type GetFileResp struct {
	File
	Trashed bool `json:"trashed"`
}

type getFileResp struct {
	Message string `json:"message"`
	GetFileResp
}
