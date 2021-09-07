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
