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
package helper_config

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	Home        string
	WorkDir     string
	CookieFile  string
	LogFile     string
	logInitOnce sync.Once
)

func init() {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("get HOME failed: %s", err))
	}

	Home = h
	WorkDir = Home + "/.go-aliyundrive-sdk"
	CookieFile = WorkDir + "/cookie.json"
	LogFile = WorkDir + "/log.log"

	err = os.MkdirAll(WorkDir, 0o777)
	if err != nil {
		panic(fmt.Errorf("create %s dir failed: %s", WorkDir, err))
	}

	InitLogger()
}

func InitLogger() {
	logInitOnce.Do(func() {
		f, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	})
}
