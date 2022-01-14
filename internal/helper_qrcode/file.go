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
package helper_qrcode

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/skip2/go-qrcode"
)

type filePrint struct{}

func (filePrint) Print(text string, level RecoveryLevel) error {
	file, err := ioutil.TempFile("", "qrcode-*.png")
	if err != nil {
		return fmt.Errorf("create qrcode failed: %w", err)
	}

	if err := qrcode.WriteFile(text, level, 256, file.Name()); err != nil {
		_ = file.Close()
		return fmt.Errorf("write qrcode failed: %w", err)
	}

	openFile(file.Name())

	return nil
}

func openFile(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
