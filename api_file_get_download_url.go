package aliyundrive

import (
	"context"
	"net/http"
	"time"
)

// GetFile 获取文件信息
func (r *FileService) GetFileDownloadURL(ctx context.Context, request *GetFileDownloadURLReq) (*GetFileDownloadURLResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "GetFileDownloadURL",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/file/get_download_url",
		Body:   request,
	}
	resp := new(getFileDownloadURLResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.GetFileDownloadURLResp, nil
}

type GetFileDownloadURLReq struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type GetFileDownloadURLResp struct {
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	InternalURL string    `json:"internal_url"`
	CdnURL      string    `json:"cdn_url"`
	Expiration  time.Time `json:"expiration"`
	Size        int       `json:"size"`
	RateLimit   struct {
		PartSpeed int `json:"part_speed"`
		PartSize  int `json:"part_size"`
	} `json:"ratelimit"`
}

type getFileDownloadURLResp struct {
	Message string `json:"message"`
	GetFileDownloadURLResp
}
