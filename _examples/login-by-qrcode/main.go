package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"

	"github.com/chyroc/go-aliyundrive"
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
