package main

import (
	"context"

	"github.com/chyroc/go-aliyundrive"
	"github.com/davecgh/go-spew/spew"
)

// 扫码登录
func main() {
	ins := aliyundrive.New(
	// aliyundrive.WithLogger(aliyundrive.NewLoggerStdout(), aliyundrive.LogLevelTrace),
	)
	ctx := context.TODO()

	user, err := ins.Auth.LoginByQrcode(ctx)
	if err != nil {
		panic(err)
	}

	spew.Dump("user", user)
}
