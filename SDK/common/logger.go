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

var (
	// context key --> log key
	ctxKeyMap = map[string]string{"request_id": "request_id"}
)

type LogOption struct {
	Level           logrus.Level
	TimestampFormat string
	FullTimestamp   bool
	HighSpeedMode   bool
	OnFile          bool
	FilePathMap     lfshook.PathMap
}

func DefaultOption() LogOption {
	return LogOption{
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

func InitLogger(option LogOption) {
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
func FlushLogger() {

}

func Logger(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}
	if ctx != nil {
		for ctxKey, logKey := range ctxKeyMap {
			if contextValue := ctx.Value(ctxKey); contextValue != nil {
				fields[logKey] = fmt.Sprint(contextValue)
			}
		}
	}

	return logrus.WithFields(fields)
}

func RegistFieldName(ctxFieldName, logFieldName string) (oldLogFieldName string) {
	if logFieldName, ok := ctxKeyMap[ctxFieldName]; ok {
		return logFieldName
	}
	ctxKeyMap[ctxFieldName] = logFieldName
	return ""
}
