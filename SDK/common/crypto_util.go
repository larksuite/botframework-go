// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// EncryptAESCBCBase64 AES CBC base64URLEncoding
func EncryptAESCBCBase64(originalData string, key string) (string, error) {
	hashKey := getHashKey(key)
	data := []byte(originalData)
	block, err := aes.NewCipher(hashKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, hashKey[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, []byte(data))
	return base64.URLEncoding.EncodeToString(crypted), nil
}

// DecryptAESCBCBase64 AES CBC base64URLEncoding
func DecryptAESCBCBase64(encryptedData string, key string) (string, error) {
	hashKey := getHashKey(key)
	crypted, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(hashKey)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, hashKey)
	originalData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originalData, crypted)
	originalData, err = PKCS5UnPadding(originalData)
	if err != nil {
		return "", err
	}
	return string(originalData), nil
}

// EncryptAes aes encrypt
func EncryptAes(original interface{}, key string) (string, error) {
	bEncrypt, err := json.Marshal(original)
	if err != nil {
		return "", fmt.Errorf("original Marshal failed. err:%v", err)
	}

	return EncryptAESCBCBase64(string(bEncrypt), key)
}

// DecryptAes aes decypt
func DecryptAes(encryptedData string, key string, original interface{}) error {
	rv := reflect.ValueOf(original)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("original isnot Ptr, or is nil")
	}

	originalData, err := DecryptAESCBCBase64(encryptedData, key)
	if err != nil {
		return fmt.Errorf("aes decypt failed. err:%v", err)
	}

	err = json.Unmarshal([]byte(originalData), &original)
	if err != nil {
		return fmt.Errorf("encryptedData Unmarshal failed. err:%v", err)
	}

	return nil
}

// PKCS5Padding PKCS5
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding PKCS5
func PKCS5UnPadding(originalData []byte) ([]byte, error) {
	length := len(originalData)
	if length <= 0 {
		return originalData, errors.New("length error")
	}
	unpadding := int(originalData[length-1])
	if length < unpadding {
		return originalData, errors.New("length lower unpadding")
	}
	return originalData[:(length - unpadding)], nil
}

func getHashKey(input string) []byte {
	h := sha256.New()
	h.Write([]byte(input))
	buf := h.Sum(nil)
	return buf[:16]
}
