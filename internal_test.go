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
