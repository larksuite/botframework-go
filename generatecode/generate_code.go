// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/larksuite/botframework-go/SDK/common"
)

func GenerateCode(ctx context.Context, tplName, tplStr, path, file string, tpl interface{}, isForceUpdate bool) error {
	fileName := path + file

	isExist, err := fileExists(fileName) // check if file exists
	if err != nil {
		return fmt.Errorf("Error checkFileIsExist file[%s]error[%v]", file, err)
	}

	if !isForceUpdate {
		if isExist {
			common.Logger(ctx).Infof("file[.%s] fileIsExist, do nothing", file)
			return nil
		} else {
			common.Logger(ctx).Infof("file[.%s] will create", file)
		}
	} else { // Force Update
		if isExist {
			common.Logger(ctx).Infof("file[.%s] will force update", file)
		} else {
			common.Logger(ctx).Infof("file[.%s] will create", file)
		}
	}

	tmpl, err := template.New(tplName).Parse(tplStr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tpl)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("write to file err: %v", err.Error())
	}

	return nil
}

func InitPath(path string) error {
	isExist, err := fileExists(path)
	if err != nil {
		return fmt.Errorf("checkPathExistError path[%s]error[%v]\n", path, err)
	}
	if !isExist {
		err = pathMkdir(path)
		if err != nil {
			return fmt.Errorf("mkdirPathError path[%s]error[%v]\n", path, err)
		}
	}

	path = path + "/handler_event"
	isExist, err = fileExists(path)
	if err != nil {
		return fmt.Errorf("checkPathExistError path[%s]error[%v]\n", path, err)
	}
	if !isExist {
		err = pathMkdir(path)
		if err != nil {
			return fmt.Errorf("mkdirPathError path[%s]error[%v]\n", path, err)
		}
	}

	return nil
}

func FormatFuncName(s string) string {
	if len(s) < 1 {
		return s
	}
	strArr := []byte(s)
	if strArr[0] >= 'a' && strArr[0] <= 'z' {
		strArr[0] -= 'a' - 'A'
	}
	return string(strArr)
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func pathMkdir(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdirFailed path[%s]err[%v]", path, err)
	}

	return nil
}
