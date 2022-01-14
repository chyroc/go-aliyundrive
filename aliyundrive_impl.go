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
	"fmt"
	"os"

	"github.com/chyroc/gorequests"

	"github.com/chyroc/go-aliyundrive/internal/helper_config"
)

type AliyunDrive struct {
	// logger
	logger   Logger
	logLevel LogLevel

	// config
	workDir string // defalut: ~/.go-aliyundrive-sdk
	store   Store

	// session
	session *gorequests.Session

	// service
	ShareLink *ShareLinkService
	Auth      *AuthService
	File      *FileService
}

func newClient(options []ClientOptionFunc) *AliyunDrive {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("get HOME failed: %s", err))
	}

	r := &AliyunDrive{
		// logger
		logLevel: LogLevelTrace,

		// timeout:      time.Second * 3,
		session: gorequests.NewSession(helper_config.CookieFile),

		// config
		workDir: home + "/.go-aliyundrive-sdk",
	}
	for _, v := range options {
		if v != nil {
			v(r)
		}
	}

	_ = os.MkdirAll(r.workDir, 0o777)
	r.initService()

	if r.logger == nil {
		r.logger = r.newDefaultLogger()
	}
	if r.store == nil {
		r.store = NewFileStore(r.workDir + "/token.json")
	}

	return r
}

func (r *AliyunDrive) initService() {
	r.ShareLink = &ShareLinkService{cli: r}
	r.Auth = &AuthService{cli: r}
	r.File = &FileService{cli: r}
}
