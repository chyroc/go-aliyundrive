package aliyundrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func (r *FileService) UploadFile(ctx context.Context, request *UploadFileReq) (*UploadFileResp, error) {
	return r.uploadFile(ctx, request.DriveID, request.ParentID, request.FilePath)
}

type UploadFileReq struct {
	DriveID  string
	ParentID string
	FilePath string
}

type UploadFileResp struct {
	DriveID            string    `json:"drive_id"`
	DomainID           string    `json:"domain_id"`
	FileID             string    `json:"file_id"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	ContentType        string    `json:"content_type"` // application/oct-stream
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	FileExtension      string    `json:"file_extension"`
	Hidden             bool      `json:"hidden"`
	Size               int       `json:"size"`
	Starred            bool      `json:"starred"`
	Status             string    `json:"status"` // available
	UploadID           string    `json:"upload_id"`
	ParentFileID       string    `json:"parent_file_id"`
	Crc64Hash          string    `json:"crc64_hash"`
	ContentHash        string    `json:"content_hash"`
	ContentHashName    string    `json:"content_hash_name"` // sha1
	Category           string    `json:"category"`
	EncryptMode        string    `json:"encrypt_mode"`
	ImageMediaMetadata struct {
		ImageQuality struct{} `json:"image_quality"`
	} `json:"image_media_metadata"`
	Location string `json:"location"`
}

func (r *FileService) uploadFile(ctx context.Context, driveID, parentID string, filepath string) (*UploadFileResp, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	} else if fileInfo.IsDir() {
		// TODO：支持文件夹
		return nil, fmt.Errorf("unsupport dir upload")
	}

	proofResp, err := r.createFileWithProof(ctx, &createFileWithProofReq{
		DriveID:       driveID,
		PartInfoList:  makePartInfoList(fileInfo.Size()),
		ParentFileID:  parentID,
		Name:          path.Base(fileInfo.Name()),
		Type:          "file",
		CheckNameMode: "auto_rename",
		Size:          fileInfo.Size(),
		PreHash:       "",
	})
	if err != nil {
		return nil, err
	}

	for _, part := range proofResp.PartInfoList {
		// TODO: 并发？
		if err := r.uploadPart(ctx, part.UploadURL, io.LimitReader(file, maxPartSize)); err != nil {
			return nil, err
		}
	}

	return r.completeUpload(ctx, &completeUploadReq{
		DriveID:  driveID,
		UploadID: proofResp.UploadID,
		FileID:   proofResp.FileID,
	})
}

func makePartInfoList(size int64) []*partInfo {
	partInfoNum := int(size / maxPartSize)
	if size%maxPartSize > 0 {
		partInfoNum += 1
	}
	res := []*partInfo{}
	for i := 0; i < partInfoNum; i++ {
		res = append(res, &partInfo{PartNumber: i + 1})
	}
	return res
}

// == create with proof ==

func (r *FileService) createFileWithProof(ctx context.Context, request *createFileWithProofReq) (*createFileWithProofResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "createFileWithProof",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/file/create_with_proof",
		Body:   request,
	}
	resp := new(createFileWithProofResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type createFileWithProofReq struct {
	DriveID       string      `json:"drive_id"`
	PartInfoList  []*partInfo `json:"part_info_list"`
	ParentFileID  string      `json:"parent_file_id"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	CheckNameMode string      `json:"check_name_mode"`
	Size          int64       `json:"size"`
	PreHash       string      `json:"pre_hash"`
}

type partInfo struct {
	PartNumber int    `json:"part_number"`
	UploadURL  string `json:"upload_url"`
}

type createFileWithProofResp struct {
	UploadID     string      `json:"upload_id"`
	FileID       string      `json:"file_id"`
	PartInfoList []*partInfo `json:"part_info_list"`
}

// == create with proof ==

// == upload part ==

func (r *FileService) uploadPart(ctx context.Context, uri string, reader io.Reader) error {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "uploadPart",
		Method: http.MethodPut,
		URL:    uri,
		Body:   reader,
	}

	response, err := r.cli.RawRequest(ctx, req, nil)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("upload_part failed, status: %d, resp: %s", response.StatusCode, response.Body)
}

// == upload part ==

// == complete upload ==

func (r *FileService) completeUpload(ctx context.Context, request *completeUploadReq) (*UploadFileResp, error) {
	req := &RawRequestReq{
		Scope:  "File",
		API:    "completeUpload",
		Method: http.MethodPost,
		URL:    "https://api.aliyundrive.com/v2/file/complete",
		Body:   request,
	}
	resp := new(UploadFileResp)

	if _, err := r.cli.RawRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type completeUploadReq struct {
	DriveID  string `json:"drive_id"`
	UploadID string `json:"upload_id"`
	FileID   string `json:"file_id"`
}

// == complete upload ==

const maxPartSize = 1024 * 1024 * 1024 // 每个分片的大小
