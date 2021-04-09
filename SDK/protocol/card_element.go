// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

import (
	"encoding/json"
	"net/url"
)

// card element
const (
	IMG_E            = "img"
	BUTTON_E         = "button"
	SELECT_STATIC_E  = "select_static"
	SELECT_PERSON_E  = "select_person"
	OVERFLOW_E       = "overflow"
	PICKERDATE_E     = "picker_date"
	PICKERTIME_E     = "picker_time"
	PICKERDATETIME_E = "picker_datetime"

	PLAIN_TEXT_E = "plain_text"
	LARK_MD_E    = "lark_md"
)

// card block
const (
	DIV_BLOCK    = "div"
	HR_BLOCK     = "hr" //halving line
	IMG_BLOCK    = "img"
	ACTION_BLOCK = "action"
	NOTE_BLOCK   = "note"
)

type ButtonStyle int

const (
	DEFAULT ButtonStyle = iota
	PRIMARY
	DANGER
	UNKNOWN
)

func (t ButtonStyle) String() string {
	return [...]string{"default", "primary", "danger", "unknown"}[t]
}

type I18NForm map[string]string

// interface - <img>, <lark_md>, <plain_text>
type BaseElement interface {
	// 获取基本组件的元素类型
	GetTag() string
}

// interface - <button>, <datepicker>, <select>, <overflow>
type ActionElement interface {
	BaseElement
	SetAction(method string, meta Meta)
	SetSession(session string)
}

type URLForm struct {
	Url        *string `json:"url,omitempty" validate:"omitempty"`
	AndroidUrl *string `json:"android_url,omitempty" validate:"omitempty"`
	IOSUrl     *string `json:"ios_url,omitempty" validate:"omitempty"`
	PCUrl      *string `json:"pc_url,omitempty" validate:"omitempty"`
}

type TextForm struct {
	Tag     string             `json:"tag,omitempty" validate:"omitempty"`
	Content *string            `json:"content,omitempty" validate:"omitempty"`
	Lines   *int               `json:"lines,omitempty" validate:"omitempty"`
	I18N    *I18NForm          `json:"i18n,omitempty" validate:"omitempty"`
	Href    map[string]URLForm `json:"href,omitempty" validate:"omitempty"`
}

func (form *TextForm) GetTag() string {
	return form.Tag
}

type ImageForm struct {
	Tag      string   `json:"tag,omitempty" validate:"omitempty"`
	ImageKey string   `json:"img_key,omitempty" validate:"omitempty"`
	ALT      TextForm `json:"alt,omitempty" validate:"omitempty"`
}

func (form *ImageForm) GetTag() string {
	return form.Tag
}

type ButtonForm struct {
	Tag      string            `json:"tag,omitempty" validate:"omitempty"`
	Text     TextForm          `json:"text,omitempty" validate:"omitempty"`
	URL      *string           `json:"url,omitempty" validate:"omitempty"`
	MultiURL *URLForm          `json:"multi_url,omitempty" validate:"omitempty"`
	Params   map[string]string `json:"value,omitempty" validate:"omitempty"`
	Type     string            `json:"type,omitempty" validate:"omitempty"`
	Confirm  *ConfirmForm      `json:"confirm,omitempty" validate:"omitempty"`
	Webhook  *string           `json:"webhook,omitempty" validate:"omitempty"`
}

func (form *ButtonForm) GetTag() string {
	return form.Tag
}

func (form *ButtonForm) SetAction(method string, meta Meta) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["method"] = method
	data, _ := json.Marshal(meta)
	form.Params["meta"] = string(data)
}

func (form *ButtonForm) SetSession(session string) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["sid"] = session
}

// select menu
type OptionForm struct {
	Text     TextForm `json:"text,omitempty" validate:"omitempty"`
	Value    string   `json:"value,omitempty" validate:"omitempty"`
	URL      *string  `json:"url,omitempty" validate:"omitempty"`
	MultiURL *URLForm `json:"multi_url,omitempty" validate:"omitempty"`
}

// selector
type SelectorForm struct {
	Tag           string            `json:"tag,omitempty" validate:"omitempty"`
	Placeholder   *TextForm         `json:"placeholder,omitempty" validate:"omitempty"`
	Options       []OptionForm      `json:"options,omitempty" validate:"omitempty"`
	InitialOption *string           `json:"initial_option,omitempty" validate:"omitempty"`
	Params        map[string]string `json:"value,omitempty" validate:"omitempty"`
	Confirm       *ConfirmForm      `json:"confirm,omitempty" validate:"omitempty"`
}

func (form *SelectorForm) GetTag() string {
	return form.Tag
}

func (form *SelectorForm) SetAction(method string, meta Meta) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["method"] = method
	data, _ := json.Marshal(meta)
	form.Params["meta"] = string(data)
}

func (form *SelectorForm) SetSession(session string) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["sid"] = session
}

// Picker Date
type PickerDateForm struct {
	Tag         string            `json:"tag,omitempty" validate:"omitempty"`
	Placeholder *TextForm         `json:"placeholder,omitempty" validate:"omitempty"`
	Params      map[string]string `json:"value,omitempty" validate:"omitempty"`
	Confirm     *ConfirmForm      `json:"confirm,omitempty" validate:"omitempty"`
	InitialDate *string           `json:"initial_date,omitempty" validate:"omitempty"`
}

func (form *PickerDateForm) GetTag() string {
	return form.Tag
}

func (form *PickerDateForm) SetAction(method string, meta Meta) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["method"] = method
	data, _ := json.Marshal(meta)
	form.Params["meta"] = string(data)
}

func (form *PickerDateForm) SetSession(session string) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["sid"] = session
}

// Picker Time
type PickerTimeForm struct {
	Tag         string            `json:"tag,omitempty" validate:"omitempty"`
	Placeholder *TextForm         `json:"placeholder,omitempty" validate:"omitempty"`
	Params      map[string]string `json:"value,omitempty" validate:"omitempty"`
	Confirm     *ConfirmForm      `json:"confirm,omitempty" validate:"omitempty"`
	InitialTime *string           `json:"initial_time,omitempty" validate:"omitempty"`
}

func (form *PickerTimeForm) GetTag() string {
	return form.Tag
}

func (form *PickerTimeForm) SetAction(method string, meta Meta) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["method"] = method
	data, _ := json.Marshal(meta)
	form.Params["meta"] = string(data)
}

func (form *PickerTimeForm) SetSession(session string) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["sid"] = session
}

// PickerDatetime控件
type PickerDatetimeForm struct {
	Tag             string            `json:"tag,omitempty" validate:"omitempty"`
	Placeholder     *TextForm         `json:"placeholder,omitempty" validate:"omitempty"`
	Params          map[string]string `json:"value,omitempty" validate:"omitempty"`
	Confirm         *ConfirmForm      `json:"confirm,omitempty" validate:"omitempty"`
	InitialDatetime *string           `json:"initial_datetime,omitempty" validate:"omitempty"`
}

func (form *PickerDatetimeForm) GetTag() string {
	return form.Tag
}

func (form *PickerDatetimeForm) SetAction(method string, meta Meta) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["method"] = method
	data, _ := json.Marshal(meta)
	form.Params["meta"] = string(data)
}

func (form *PickerDatetimeForm) SetSession(session string) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["sid"] = session
}

// overflow
type OverflowForm struct {
	Tag     string            `json:"tag,omitempty" validate:"omitempty"`
	Options []OptionForm      `json:"options,omitempty" validate:"omitempty"`
	Params  map[string]string `json:"value,omitempty" validate:"omitempty"`
	Confirm *ConfirmForm      `json:"confirm,omitempty" validate:"omitempty"`
}

func (form *OverflowForm) GetTag() string {
	return form.Tag
}

func (form *OverflowForm) SetAction(method string, meta Meta) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["method"] = method
	data, _ := json.Marshal(meta)
	form.Params["meta"] = string(data)
}

func (form *OverflowForm) SetSession(session string) {
	if form.Params == nil {
		form.Params = make(map[string]string, 0)
	}
	form.Params["sid"] = session
}

// Field
type FieldForm struct {
	Short bool     `json:"is_short,omitempty" validate:"omitempty"`
	Text  TextForm `json:"text,omitempty" validate:"omitempty"`
}

// toast
type ConfirmForm struct {
	Title   TextForm `json:"title,omitempty" validate:"omitempty"`
	Text    TextForm `json:"text,omitempty" validate:"omitempty"`
	Confirm TextForm `json:"confirm,omitempty" validate:"omitempty"`
	Deny    TextForm `json:"deny,omitempty" validate:"omitempty"`
}

func wrapH5URL(rawURL string, session string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return rawURL, err
	}
	q.Add("sid", session)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// card block - div
type DIVBlockForm struct {
	Tag   string      `json:"tag,omitempty" validate:"omitempty"`
	Text  *TextForm   `json:"text,omitempty" validate:"omitempty"`
	Field []FieldForm `json:"fields,omitempty" validate:"omitempty"`
	Extra interface{} `json:"extra,omitempty" validate:"omitempty"`
}

// card block - hr
type HRBlockForm struct {
	Tag string `json:"tag,omitempty" validate:"omitempty"`
}

// card block - image
type ImageBlockForm struct {
	Tag      string    `json:"tag,omitempty" validate:"omitempty"`
	ImageKey string    `json:"img_key,omitempty" validate:"omitempty"`
	Title    *TextForm `json:"title,omitempty" validate:"omitempty"`
	Alt      TextForm  `json:"alt,omitempty" validate:"omitempty"`
}

// card block - action
type ActionBlockForm struct {
	Tag     string          `json:"tag,omitempty" validate:"omitempty"`
	Actions []ActionElement `json:"actions,omitempty" validate:"omitempty"`
}

// card block - note
type NoteBlockForm struct {
	Tag      string        `json:"tag,omitempty" validate:"omitempty"`
	Elements []BaseElement `json:"elements,omitempty" validate:"omitempty"`
}

// card form
type ConfigForm struct {
	MinVersion     VersionForm `json:"min_version,omitempty" validate:"omitempty"`
	Debug          bool        `json:"debug,omitempty" validate:"omitempty"`
	WideScreenMode bool        `json:"wide_screen_mode"`
	EnableForward  bool        `json:"enable_forward,omitempty" validate:"omitempty"`
}

type VersionForm struct {
	Version        string `json:"version,omitempty" validate:"omitempty"`
	AndroidVersion string `json:"android_version,omitempty" validate:"omitempty"`
	IOSVersion     string `json:"ios_version,omitempty" validate:"omitempty"`
	PCVersion      string `json:"pc_version,omitempty" validate:"omitempty"`
}

type CardHeaderForm struct {
	Title    TextForm `json:"title,omitempty" validate:"omitempty"`
	Template string   `json:"template,omitempty" validate:"omitempty"`
}

type CardForm struct {
	OpenIDs      []string                 `json:"open_ids,omitempty" validate:"omitempty"`
	Config       *ConfigForm              `json:"config,omitempty" validate:"omitempty"`
	CardLink     *URLForm                 `json:"card_link,omitempty" validate:"omitempty"`
	Header       *CardHeaderForm          `json:"header,omitempty" validate:"omitempty"`
	Elements     []interface{}            `json:"elements,omitempty" validate:"omitempty"`
	I18NElements map[string][]interface{} `json:"i18n_elements,omitempty" validate:"omitempty"`
}

type ToastTips struct {
	Content string            `json:"content,omitempty" validate:"omitempty"`
	I18N    map[string]string `json:"i18n,omitempty" validate:"omitempty"`
}
