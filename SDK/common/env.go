// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

var (
	isOutsideMainlandChina bool = false
	isStaging              bool = false // Used only for internal testing
)

// outside mainland China, need set "OutsideMainlandChina" flag
func SetOutsideMainlandChina(flag bool) {
	isOutsideMainlandChina = flag
}

// Used only for internal testing
func SetStaging(flag bool) {
	isStaging = flag
}

// GetOpenPlatformHost outside mainland China use https://open.larksuite.com , in Mainland China use https://open.feishu.cn
func GetOpenPlatformHost() string {
	if isOutsideMainlandChina {
		if isStaging {
			return "https://open.larksuite-staging.com"
		}
		return "https://open.larksuite.com"
	}

	if isStaging {
		return "https://open.feishu-staging.cn"
	}
	return "https://open.feishu.cn"
}
