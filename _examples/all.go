package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chyroc/go-aliyundrive"
)

func main() {
	r := aliyundrive.New()

	ctx := context.Background()

	// 扫码登录
	user, err := r.Auth.LoginByQrcode(ctx, &aliyundrive.LoginByQrcodeReq{
		SmallQrCode: true,
	})
	assert(err)

	fmt.Println("user:", jsonString(user))

	// 获取 box 信息
	box, err := r.File.GetSBox(ctx)
	assert(err)
	fmt.Println("box:", jsonString(box))

	// 获取根目录文件列表
	fileList, err := r.File.GetFileList(ctx, &aliyundrive.GetFileListReq{
		DriveID:      user.DefaultDriveID,
		ParentFileID: aliyundrive.RootFileID,
	})
	assert(err)
	fmt.Println("fileList:", jsonString(fileList))

	// 创建文件夹
	folderResp, err := r.File.CreateFolder(ctx, &aliyundrive.CreateFolderReq{
		DriveID:       user.DefaultDriveID,
		ParentFileID:  aliyundrive.RootFileID,
		Name:          "just-for-example",
		CheckNameMode: "",
		Type:          "",
	})
	assert(err)
	fmt.Println("create folder:", jsonString(folderResp))

	delFile, err := r.File.DeleteFile(ctx, &aliyundrive.DeleteFileReq{
		DriveID: user.DefaultDriveID,
		FileID:  folderResp.FileID,
	})
	assert(err)
	fmt.Println("delete file:", jsonString(delFile))
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func jsonString(v interface{}) string {
	bs, _ := json.MarshalIndent(v, "", "  ")
	return string(bs)
}
