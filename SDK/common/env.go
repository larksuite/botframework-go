// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

var (
	isFeishuOrLark bool = true

	hostFeishu string = "https://open.feishu.cn"
	hostLark   string = "https://open.larksuite.com"
)

func SetFeishu() {
	isFeishuOrLark = true
}

func SetLark() {
	isFeishuOrLark = false
}

func ReplaceFeishuHost(host string) {
	hostFeishu = host
}

func ReplaceLarkHost(host string) {
	hostLark = host
}

// GetOpenPlatformHost outside mainland China use https://open.larksuite.com , in Mainland China use https://open.feishu.cn
func GetOpenPlatformHost() string {
	if isFeishuOrLark {
		return hostFeishu
	}

	return hostLark
}
