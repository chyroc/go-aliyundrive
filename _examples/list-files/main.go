/**
 * Copyright 2022 chyroc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
