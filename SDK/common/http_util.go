// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/larksuite/botframework-go/SDK/protocol"
)

// DoHttpPostOApi open platform POST http
func DoHttpPostOApi(path protocol.OpenApiPath, headers map[string]string, data interface{}) ([]byte, error) {
	reqBody := new(bytes.Buffer)
	err := json.NewEncoder(reqBody).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("jsonEncodeError[%v]", err)
	}

	reqURL := GetOpenPlatformHost() + string(path)

	return DoHttp("POST", reqURL, headers, reqBody)
}

// DoHttpGetOApi open platform GET http
func DoHttpGetOApi(path protocol.OpenApiPath, headers map[string]string, params map[string]string) ([]byte, error) {
	// http URL
	reqURL := GetOpenPlatformHost() + string(path)

	if params != nil && len(params) > 0 {
		m := make(url.Values)
		for k, v := range params {
			m.Set(k, v)
		}

		reqURL = reqURL + "?" + m.Encode()
	}

	reqBody := new(bytes.Buffer)
	return DoHttp("GET", reqURL, headers, reqBody)
}

func DoHttp(method string, url string, headers map[string]string, body *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("httpNewRequestError[%v]", err)
	}

	// http header
	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("httpDoError[%v]", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("httpRetError statusCode[%d]", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readRespBodyError[%v]", err)
	}

	return respBody, nil
}

func NewHeaderToken(accessToken string) map[string]string {
	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	header["Content-Type"] = "application/json"
	return header
}
