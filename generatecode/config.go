// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generatecode

type GenCodeInfo struct {
	ServiceInfo    ServiceInfoConf `yaml:"ServiceInfo"`
	EventList      []EventConf     `yaml:"EventList"`
	CommandList    []CommandConf   `yaml:"CommandList"`
	CardActionList []CardConf      `yaml:"CardActionList"`
}

type ServiceInfoConf struct {
	Path         string `yaml:"Path"`
	GenCodePath  string `yaml:"GenCodePath"`
	EventWebhook string `yaml:"EventWebhook"`
	CardWebhook  string `yaml:"CardWebhook"`
	AppID        string `yaml:"AppID"`
	Description  string `yaml:"Description"`
	IsISVApp     bool   `yaml:"IsISVApp"`
}

type EventConf struct {
	EventName string `yaml:"EventName"`
}

type CommandConf struct {
	Cmd         string `yaml:"Cmd"`
	Description string `yaml:"Description"`
}

type CardConf struct {
	MethodName string `yaml:"MethodName"`
}
