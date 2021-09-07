package aliyundrive

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Req struct {
	Name   string `json:"name"`
	Path   string `path:"path"`
	Action string `query:"action"`
	Sender string `json:"sender" default:"bob"`
}

func Test_parseReq(t *testing.T) {
	t.SkipNow()
	as := assert.New(t)
	res, err := parseRequestParam(&RawRequestReq{
		Scope:  "scope",
		API:    "api",
		Method: http.MethodPost,
		URL:    "https://url.com/:path/do",
		Body: Req{
			Name:   "name1",
			Path:   "path2",
			Action: "action3",
		},
		IsFile:  false,
		headers: nil,
	})
	as.Nil(err)
	as.Equal(http.MethodPost, res.Method)
	as.Equal("https://url.com/path2/do", res.URL)
	as.Equal(http.MethodPost, res.RawBody)
}
