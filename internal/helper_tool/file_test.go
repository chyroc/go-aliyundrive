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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tool(t *testing.T) {
	as := assert.New(t)

	IsFileExist = func(file string) bool {
		return file == "1(2).txt" || file == "1.txt(2).txt"
	}

	as.Equal("2(2)", AutoRenameFile("2"))
	as.Equal("2(2).txt", AutoRenameFile("2.txt"))
	as.Equal("2.txt(2).txt", AutoRenameFile("2.txt.txt"))

	as.Equal("1(3).txt", AutoRenameFile("1.txt"))
	as.Equal("1.txt(3).txt", AutoRenameFile("1.txt.txt"))
}
