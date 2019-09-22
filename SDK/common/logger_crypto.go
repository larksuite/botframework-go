// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

var (
	logEncryptKey = ""
)

// to encrypt/decrypt information more conveniently, when printing logs.
func InitEncryptKey(encryptKey string) {
	logEncryptKey = encryptKey
}

func LogEncrypt(origData string) string {
	if logEncryptKey == "" {
		return origData
	}

	data, _ := EncryptDESCBCBase64(origData, logEncryptKey)
	return data
}

func LogDecrypt(encryptedData string) string {
	if logEncryptKey == "" {
		return encryptedData
	}

	data, _ := DecryptDESCBCBase64(encryptedData, logEncryptKey)
	return data
}
