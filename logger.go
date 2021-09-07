package aliyundrive

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Log(ctx context.Context, level LogLevel, msg string, args ...interface{})
}

type LogLevel int

const (
	LogLevelTrace LogLevel = iota + 1 // 只有两个 log req 和 resp 的 文本内容
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (r LogLevel) String() string {
	switch r {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return ""
	}
}

func (r *AliyunDrive) log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	if r.logger != nil && r.logLevel <= level {
		r.logger.Log(ctx, level, "[aliyundrive] "+msg, args...)
	}
}

type LoggerStdout struct{}

func NewLoggerStdout() Logger {
	return &LoggerStdout{}
}

func (l *LoggerStdout) Log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	fmt.Printf("["+level.String()+"] "+msg+"\n", args...)
}

type DefaultLogger struct {
	logger *logrus.Logger
}

func (r *AliyunDrive) newDefaultLogger() Logger {
	logger := logrus.New()
	f, err := os.OpenFile(r.workDir+"/log.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(f)
	return &DefaultLogger{
		logger: logger,
	}
}

func (l *DefaultLogger) Log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	l.logger.Printf("["+level.String()+"] "+msg+"\n", args...)
}
