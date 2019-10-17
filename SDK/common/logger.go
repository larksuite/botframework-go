// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type CommonLoggerOption struct {
	Level           logrus.Level
	TimestampFormat string
	FullTimestamp   bool
	HighSpeedMode   bool
	OnFile          bool
	FilePathMap     lfshook.PathMap
}

func DefaultOption() *CommonLoggerOption {
	return &CommonLoggerOption{
		Level:           logrus.DebugLevel,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		HighSpeedMode:   false,
		OnFile:          false,
		FilePathMap: lfshook.PathMap{
			logrus.DebugLevel: "./default.log",
			logrus.InfoLevel:  "./default.log",
			logrus.WarnLevel:  "./default.log",
			logrus.ErrorLevel: "./default.log",
		},
	}
}

type CommonLogger struct {
	KeyMap map[string]string
}

func NewCommonLogger() *CommonLogger {
	log := &CommonLogger{}
	log.RegistFieldName("request_id", "request_id")

	return log
}

func (l *CommonLogger) Init(op interface{}) {
	option, ok := op.(*CommonLoggerOption)
	if !ok {
		option = DefaultOption()
	}

	logrus.SetLevel(option.Level)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = option.TimestampFormat
	customFormatter.FullTimestamp = option.FullTimestamp
	logrus.SetFormatter(customFormatter)

	if option.HighSpeedMode {
		logrus.StandardLogger().SetNoLock()
		logrus.SetOutput(ioutil.Discard)
	}

	if option.OnFile {
		fileHook := lfshook.NewHook(option.FilePathMap, customFormatter)
		logrus.AddHook(fileHook)
	}
}

// FlushLogger Flush
func (l *CommonLogger) Flush() {

}

func (l *CommonLogger) GetLogger(ctx context.Context) LogEntry {
	fields := logrus.Fields{}
	if ctx != nil {
		for ctxKey, logKey := range l.KeyMap {
			if contextValue := ctx.Value(ctxKey); contextValue != nil {
				fields[logKey] = fmt.Sprint(contextValue)
			}
		}
	}

	return logrus.WithFields(fields)
}

func (l *CommonLogger) RegistFieldName(ctxFieldName, logFieldName string) (oldLogFieldName string) {
	if l.KeyMap == nil {
		l.KeyMap = make(map[string]string)
	}

	if logFieldName, ok := l.KeyMap[ctxFieldName]; ok {
		return logFieldName
	}

	l.KeyMap[ctxFieldName] = logFieldName
	return ""
}
