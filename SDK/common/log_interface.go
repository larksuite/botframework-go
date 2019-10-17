package common

import (
	"context"
)

type LogInterface interface {
	Init(option interface{})
	Flush()
	GetLogger(ctx context.Context) LogEntry
}

type LogEntry interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
}

var (
	logger LogInterface
)

func InitLogger(log LogInterface, option interface{}) {
	logger = log
	logger.Init(option)
}

func FlushLogger() {
	if logger == nil {
		logger = NewCommonLogger()
		logger.Init(DefaultOption())
	}
	logger.Flush()
}

func Logger(ctx context.Context) LogEntry {
	if logger == nil {
		logger = NewCommonLogger()
		logger.Init(DefaultOption())
	}

	return logger.GetLogger(ctx)
}
