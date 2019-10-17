// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

var (
	aesEncryptKey = ""
)

// to encrypt/decrypt information more conveniently

func InitAESEncryptKey(encryptKey string) {
	aesEncryptKey = encryptKey
}

func GetAESEncryptKey() string {
	return aesEncryptKey
}

func AESEncrypt(origData string) string {
	if aesEncryptKey == "" {
		return origData
	}

	data, _ := EncryptAESCBCBase64(origData, aesEncryptKey)
	return data
}

func AESDecrypt(encryptedData string) string {
	if aesEncryptKey == "" {
		return encryptedData
	}

	data, _ := DecryptAESCBCBase64(encryptedData, aesEncryptKey)
	return data
}
