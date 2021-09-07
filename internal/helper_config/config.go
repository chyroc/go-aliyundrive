package helper_config

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	Home        string
	WorkDir     string
	CookieFile  string
	LogFile     string
	logInitOnce sync.Once
)

func init() {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("get HOME failed: %s", err))
	}

	Home = h
	WorkDir = Home + "/.go-aliyundrive-sdk"
	CookieFile = WorkDir + "/cookie.json"
	LogFile = WorkDir + "/log.log"

	err = os.MkdirAll(WorkDir, 0o777)
	if err != nil {
		panic(fmt.Errorf("create %s dir failed: %s", WorkDir, err))
	}

	InitLogger()
}

func InitLogger() {
	logInitOnce.Do(func() {
		f, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	})
}
