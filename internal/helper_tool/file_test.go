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
