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
package helper_tool

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var IsFileExist = func(file string) bool {
	f, err := os.Stat(file)
	return f != nil || (err != nil && os.IsExist(err))
}

func AutoRenameFile(file string) string {
	for {
		file = autoRenameFile(file)
		if !IsFileExist(file) {
			return file
		}
	}
}

func autoRenameFile(file string) string {
	ext := filepath.Ext(file)
	if ext != file && ext != "" {
		file = file[:len(file)-len(ext)]
	}
	i := int64(1)
	if match := regEndWithNumber.FindStringSubmatch(file); len(match) == 2 {
		i, _ = strconv.ParseInt(match[1], 10, 64)
		file = file[:len(file)-(len(match[1])+2)]
	}
	return fmt.Sprintf("%s(%d)%s", file, i+1, ext)
}

var regEndWithNumber = regexp.MustCompile(`.*?\((\d+)\)$`)
