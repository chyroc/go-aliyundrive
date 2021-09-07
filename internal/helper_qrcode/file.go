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
