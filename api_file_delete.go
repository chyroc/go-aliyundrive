package aliyundrive

import (
	"context"
	"net/http"
)

// DeleteFile 删除文件
func (r *FileService) DeleteFile(ctx context.Context, request *DeleteFileReq) (*DeleteFileResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "DeleteFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/recyclebin/trash",
		Body:   request,
	}
	resp := new(deleteFileResp)

	if _, err := r.cli.RawRequest(ctx, req, nil); err != nil {
		return nil, err
	}
	return &resp.DeleteFileResp, nil
}

type DeleteFileReq struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type DeleteFileResp struct {
	DomainID    string `json:"domain_id"`
	DriveID     string `json:"drive_id"`
	FileID      string `json:"file_id"`
	AsyncTaskID string `json:"async_task_id"`
}

type deleteFileResp struct {
	Message string `json:"message"`
	DeleteFileResp
}
