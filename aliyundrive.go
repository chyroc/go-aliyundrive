package aliyundrive

import (
	"path/filepath"
)

func New(options ...ClientOptionFunc) *AliyunDrive {
	return newClient(options)
}

type ClientOptionFunc func(*AliyunDrive)

func WithLogger(logger Logger, level LogLevel) ClientOptionFunc {
	return func(ins *AliyunDrive) {
		ins.logger = logger
		ins.logLevel = level
	}
}

func WithWorkDir(dir string) ClientOptionFunc {
	return func(ins *AliyunDrive) {
		path, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}
		ins.workDir = path
	}
}

func WithStore(store Store) ClientOptionFunc {
	return func(ins *AliyunDrive) {
		ins.store = store
	}
}

type AuthService struct {
	cli *AliyunDrive
}

type FileService struct {
	cli *AliyunDrive
}

type ShareLinkService struct {
	cli *AliyunDrive
}
