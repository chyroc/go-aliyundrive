package main

import (
	"context"
	"fmt"

	"github.com/chyroc/go-aliyundrive"
)

func main() {
	// options := aliyundrive.WithLogger(aliyundrive.NewLoggerStdout(), aliyundrive.LogLevelTrace)
	ins := aliyundrive.New()
	ctx := context.TODO()

	user, err := ins.Auth.LoginByQrcode(ctx)
	if err != nil {
		panic(err)
	}

	liftCount := 3
	listAllFile(ins, "", user.DefaultDriveID, "", &liftCount)
}

func listAllFile(ins *aliyundrive.AliyunDrive, prefix string, driveID, parentID string, liftCount *int) {
	if *liftCount < 0 {
		return
	}
	next := ""
	for {
		resp, err := ins.File.GetFileList(context.TODO(), &aliyundrive.GetFileListReq{
			DriveID:      driveID,
			ParentFileID: parentID,
			Marker:       next,
		})
		if err != nil {
			panic(err)
		}
		next = resp.NextMarker
		for _, v := range resp.Items {
			fmt.Println(prefix, v.Name, v.Type)
			*liftCount--
			if v.Type == "folder" {
				listAllFile(ins, prefix+"  ", driveID, v.FileID, liftCount)
			}
		}
		if next == "" {
			break
		}
	}
}
