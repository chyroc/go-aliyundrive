package helper_qrcode

import "testing"

func TestPrintQRSmall(t *testing.T) {
	(&consolePrint{}).Print("123", Low)
}
