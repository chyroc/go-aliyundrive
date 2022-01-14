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
	"path/filepath"
)

func New(options ...ClientOptionFunc) *AliyunDrive {
	return newClient(options)
}

type ClientOptionFunc func(*AliyunDrive)

func WithLogger(logger Logger, level LogLevel) ClientOptionFunc {
	return func(ins *AliyunDrive) {
		ins.logger = logger
		ins.logLevel = level
	}
}

func WithWorkDir(dir string) ClientOptionFunc {
	return func(ins *AliyunDrive) {
		path, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}
		ins.workDir = path
	}
}

func WithStore(store Store) ClientOptionFunc {
	return func(ins *AliyunDrive) {
		ins.store = store
	}
}

type AuthService struct {
	cli *AliyunDrive
}

type FileService struct {
	cli *AliyunDrive
}

type ShareLinkService struct {
	cli *AliyunDrive
}
