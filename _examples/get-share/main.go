package main

import (
	"context"
	"flag"

	"github.com/davecgh/go-spew/spew"

	"github.com/chyroc/go-aliyundrive"
)

func main() {
	ins := aliyundrive.New()
	ctx := context.TODO()

	shareID := ""
	flag.StringVar(&shareID, "share", "", "share id")
	flag.Parse()

	sharedInfo, err := ins.ShareLink.GetShareByAnonymous(ctx, &aliyundrive.GetShareByAnonymousReq{
		ShareID: shareID,
	})
	if err != nil {
		panic(err)
	}
	spew.Dump(sharedInfo)
}
