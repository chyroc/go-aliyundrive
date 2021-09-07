package aliyundrive

import (
	"context"
	"net/http"
)

// RenameFile 重命名
func (r *FileService) MoveFile(ctx context.Context, request *MoveFileReq) (*MoveFileResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "MoveFile",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v3/file/move",
		Body:   request,
	}
	resp := new(moveFileResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.MoveFileResp, nil
}

type MoveFileReq struct {
	DriveID        string `json:"drive_id"`
	FileID         string `json:"file_id"`
	ToDriveID      string `json:"to_drive_id"`
	ToParentFileID string `json:"to_parent_file_id"`
}

type MoveFileResp struct {
	DomainID string `json:"domain_id"`
	DriveID  string `json:"drive_id"`
	FileID   string `json:"file_id"`
}

type moveFileResp struct {
	Message string `json:"message"`
	MoveFileResp
}
