package helper_qrcode

import "github.com/skip2/go-qrcode"

type RecoveryLevel = qrcode.RecoveryLevel

const (
	Low RecoveryLevel = iota
	Medium
	High
	Highest
)

func New(isConsole bool) Print {
	if isConsole {
		return consolePrint{}
	}
	return filePrint{}
}
