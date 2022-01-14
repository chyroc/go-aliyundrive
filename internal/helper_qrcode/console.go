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

	"github.com/skip2/go-qrcode"
)

type Print interface {
	Print(text string, level RecoveryLevel) error
}

type consolePrint struct{}

func (consolePrint) Print(text string, level RecoveryLevel) (err error) {
	var obj *qrcode.QRCode
	obj, err = qrcode.New(text, level)
	if err != nil {
		return
	}
	fmt.Print(obj.ToSmallString(false))
	return err
}
