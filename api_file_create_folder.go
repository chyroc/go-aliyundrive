package aliyundrive

import (
	"context"
	"net/http"
)

// CreateFolder 创建文件夹
func (r *FileService) CreateFolder(ctx context.Context, request *CreateFolderReq) (*CreateFolderResp, error) {
	request.CheckNameMode = "refuse"
	request.Type = "folder"

	req := &RawRequestReq{
		Scope:  "File",
		API:    "CreateFolder",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/adrive/v2/file/createWithFolders",
		Body:   request,
	}
	resp := new(createFolderResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return &resp.CreateFolderResp, nil
}

type CreateFolderReq struct {
	DriveID       string `json:"drive_id"`
	ParentFileID  string `json:"parent_file_id"`
	Name          string `json:"name"`
	CheckNameMode string `json:"check_name_mode"`
	Type          string `json:"type"`
}

type CreateFolderResp struct {
	ParentFileID string `json:"parent_file_id"`
	Type         string `json:"type"`
	FileID       string `json:"file_id"`
	DomainID     string `json:"domain_id"`
	DriveID      string `json:"drive_id"`
	FileName     string `json:"file_name"`
	EncryptMode  string `json:"encrypt_mode"`
}

type createFolderResp struct {
	Message string `json:"message"`
	CreateFolderResp
}
