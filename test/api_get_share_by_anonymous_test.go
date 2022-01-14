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
package test

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/chyroc/go-aliyundrive"
)

func Test_GetShare(t *testing.T) {
	t.Skip()

	as := assert.New(t)
	ctx := context.TODO()

	ins := aliyundrive.New(aliyundrive.WithLogger(aliyundrive.NewLoggerStdout(), aliyundrive.LogLevelTrace))

	res, err := ins.ShareLink.GetShareByAnonymous(ctx, &aliyundrive.GetShareByAnonymousReq{
		ShareID: "YBJMK2np1uc",
	})
	as.Nil(err)
	as.NotNil(res)
	spew.Dump(res)

	as.Equal("0613badbb0d94b01b96ad8cdaa7184b6", res.CreatorID)
	as.Equal("图灵程序丛书-302本", res.ShareName)
	as.Equal(1, res.FileCount)
	as.Len(res.FileInfos, 1)
	as.Equal("folder", res.FileInfos[0].Type)
	as.Equal("图灵程序丛书-302本", res.FileInfos[0].FileName)
}
