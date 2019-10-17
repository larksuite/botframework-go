// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

// card meta
const (
	VER = "1.0.0"
)

// card meta info
type Meta struct {
	SDKVersion string `json:"sdk_version,omitempty" validate:"omitempty"`
}

func NewMeta() *Meta {
	meta := &Meta{}
	meta.SDKVersion = VER
	return meta
}

// card callback
type CardChallenge struct {
	AppID     string `json:"appid,omitempty" validate:"omitempty"`
	Challenge string `json:"challenge,omitempty" validate:"omitempty"`
	Token     string `json:"token,omitempty" validate:"omitempty"`
	Type      string `json:"type,omitempty" validate:"omitempty"`
}

type CardCallbackBase struct {
	OpenID        string  `json:"open_id,omitempty" validate:"omitempty"`
	UserID        string  `json:"user_id,omitempty" validate:"omitempty"`
	OpenMessageID string  `json:"open_message_id,omitempty" validate:"omitempty"`
	TenantKey     string  `json:"tenant_key,omitempty" validate:"omitempty"`
	Token         *string `json:"token,omitempty" validate:"omitempty"`
}

type CardCallbackForm struct {
	CardCallbackBase
	Action CardCallBackAction `json:"action,omitempty" validate:"omitempty"`
}

type CardCallBackAction struct {
	Value  map[string]string `json:"value,omitempty" validate:"omitempty"`
	Tag    *string           `json:"tag,omitempty" validate:"omitempty"`
	Option *string           `json:"option,omitempty" validate:"omitempty"`
}
