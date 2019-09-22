// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"runtime"
	"runtime/debug"
)

func RecoverPanic(ctx context.Context) {
	if err := recover(); err != nil {
		pc, file, line, _ := runtime.Caller(3)
		f := runtime.FuncForPC(pc)
		Logger(ctx).Errorf("Recover: funcName=%v, file=%s, line=%d, panic:%v, stack info:%v",
			f.Name(), file, line, err, string(debug.Stack()))
	}
}

func GetMd5(src io.Reader) string {
	hash := md5.New()
	io.Copy(hash, src)
	return hex.EncodeToString(hash.Sum(nil))
}

func GetMd5ByBytes(src []byte) string {
	hash := md5.New()
	hash.Write(src)
	return hex.EncodeToString(hash.Sum(nil))
}
