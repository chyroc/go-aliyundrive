package test

import (
	"context"
	"testing"

	"github.com/chyroc/go-aliyundrive"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_GetShare(t *testing.T) {
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
