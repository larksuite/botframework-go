// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package message

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/larksuite/botframework-go/SDK/protocol"
)

type CardBuilder struct {
	blocks        []interface{} // block list: DIVBlock/HRBlock/ImageBlock/ActionBlock/NoteBlock
	i18nBlocks    map[string][]interface{}
	currentLocale string
	Header        *protocol.CardHeaderForm
	Conf          *protocol.ConfigForm
	OpenIDs       []string
	Session       *string
}

// switch to the i18n locale
func (builder *CardBuilder) SwitchLocale(locale protocol.Language) *CardBuilder {
	if builder.i18nBlocks == nil {
		builder.i18nBlocks = make(map[string][]interface{}, 0)
	}
	builder.currentLocale = locale.String()

	return builder
}

// set ConfigForm
func (builder *CardBuilder) SetConfig(config protocol.ConfigForm) *CardBuilder {
	builder.Conf = &config

	return builder
}

// set updating user group
func (builder *CardBuilder) SetUpdatingUserGroup(users []string) *CardBuilder {
	builder.OpenIDs = users
	return builder
}

// add header
// @param title : TextForm, tag must be "plain_text". Can call the NewPlainText function to get it
// @param template : Reserved field. Fill in the empty string
func (builder *CardBuilder) AddHeader(title protocol.TextForm, template string) *CardBuilder {
	header := &protocol.CardHeaderForm{}
	header.Title = title
	header.Template = template
	builder.Header = header

	return builder
}

// add div block
func (builder *CardBuilder) AddDIVBlock(text *protocol.TextForm, field []protocol.FieldForm, extra interface{}) *CardBuilder {
	// generate a new session
	if builder.Session == nil {
		sid := uuid.New().String()
		builder.Session = &sid
	}

	div := &protocol.DIVBlockForm{}
	div.Tag = protocol.DIV_BLOCK
	div.Text = text
	div.Field = field

	if act, ok := extra.(protocol.ActionElement); ok {
		div.Extra = act
		act.SetSession(*builder.Session)
	} else if base, ok := extra.(protocol.BaseElement); ok {
		div.Extra = base
	}

	if builder.currentLocale != "" {
		builder.i18nBlocks[builder.currentLocale] = append(builder.i18nBlocks[builder.currentLocale], div)
	} else {
		builder.blocks = append(builder.blocks, div)
	}

	return builder
}

// add hr block
func (builder *CardBuilder) AddHRBlock() *CardBuilder {
	hr := &protocol.DIVBlockForm{}
	hr.Tag = protocol.HR_BLOCK

	if builder.currentLocale != "" {
		builder.i18nBlocks[builder.currentLocale] = append(builder.i18nBlocks[builder.currentLocale], hr)
	} else {
		builder.blocks = append(builder.blocks, hr)
	}
	return builder
}

// add image block
func (builder *CardBuilder) AddImageBlock(title *protocol.TextForm, alt protocol.TextForm, imageKey string) *CardBuilder {
	img := &protocol.ImageBlockForm{}
	img.Tag = protocol.IMG_BLOCK
	img.Title = title
	img.Alt = alt
	img.ImageKey = imageKey

	if builder.currentLocale != "" {
		builder.i18nBlocks[builder.currentLocale] = append(builder.i18nBlocks[builder.currentLocale], img)
	} else {
		builder.blocks = append(builder.blocks, img)
	}
	return builder
}

// add action block
func (builder *CardBuilder) AddActionBlock(actions []protocol.ActionElement) *CardBuilder {
	// generate a new session
	if builder.Session == nil {
		sid := uuid.New().String()
		builder.Session = &sid
	}

	act := &protocol.ActionBlockForm{}
	act.Tag = protocol.ACTION_BLOCK
	for _, a := range actions {
		a.SetSession(*builder.Session)
	}
	act.Actions = actions

	if builder.currentLocale != "" {
		builder.i18nBlocks[builder.currentLocale] = append(builder.i18nBlocks[builder.currentLocale], act)
	} else {
		builder.blocks = append(builder.blocks, act)
	}
	return builder
}

// add note block
func (builder *CardBuilder) AddNoteBlock(elements []protocol.BaseElement) *CardBuilder {
	note := &protocol.NoteBlockForm{}
	note.Tag = protocol.NOTE_BLOCK
	note.Elements = elements

	if builder.currentLocale != "" {
		builder.i18nBlocks[builder.currentLocale] = append(builder.i18nBlocks[builder.currentLocale], note)
	} else {
		builder.blocks = append(builder.blocks, note)
	}
	return builder
}

// build card
func (builder *CardBuilder) Build() (data []byte, err error) {
	// build card json data
	card := &protocol.CardForm{}
	card.Elements = builder.blocks
	card.I18NElements = builder.i18nBlocks
	card.Header = builder.Header
	card.Config = builder.Conf

	data, err = json.Marshal(card)
	return
}

// build card
func (builder *CardBuilder) BuildForm() (card *protocol.CardForm, err error) {
	// build card json form
	card = &protocol.CardForm{}
	card.Elements = builder.blocks
	card.I18NElements = builder.i18nBlocks
	card.Header = builder.Header
	card.Config = builder.Conf
	card.OpenIDs = builder.OpenIDs

	return
}

// new element plain_text
func NewPlainText(content *string, i18n *protocol.I18NForm, lines *int) *protocol.TextForm {
	form := &protocol.TextForm{}
	form.Tag = protocol.PLAIN_TEXT_E
	form.Content = content
	form.Lines = lines
	form.I18N = i18n

	return form
}

// new element lark_md
func NewMDText(content string, lines *int, i18n *protocol.I18NForm, href map[string]protocol.URLForm) *protocol.TextForm {
	form := &protocol.TextForm{}
	form.Tag = protocol.LARK_MD_E
	form.Content = &content
	form.Lines = lines
	form.I18N = i18n
	form.Href = href

	return form
}

// new element img
func NewImage(alt *protocol.TextForm, imageKey string) *protocol.ImageForm {
	form := &protocol.ImageForm{}
	form.Tag = protocol.IMG_E
	form.ImageKey = imageKey
	form.ALT = *alt

	return form
}

// new element button
func NewButton(text *protocol.TextForm, url *string, multiURL *protocol.URLForm, params map[string]string,
	category protocol.ButtonStyle, confirm *protocol.ConfirmForm, method string) *protocol.ButtonForm {

	form := &protocol.ButtonForm{}
	form.Tag = protocol.BUTTON_E
	form.Text = *text
	form.URL = url
	form.MultiURL = multiURL
	form.Params = params
	form.Type = category.String()
	form.Confirm = confirm
	form.SetAction(method, *protocol.NewMeta())

	return form
}

// new element button
func NewJumpButton(text *protocol.TextForm, url *string, multiURL *protocol.URLForm,
	category protocol.ButtonStyle) *protocol.ButtonForm {

	form := &protocol.ButtonForm{}
	form.Tag = protocol.BUTTON_E
	form.Text = *text
	form.URL = url
	form.MultiURL = multiURL
	form.Type = category.String()

	return form
}

// new element button
func NewActionButton(text *protocol.TextForm, params map[string]string,
	category protocol.ButtonStyle, confirm *protocol.ConfirmForm, method string) *protocol.ButtonForm {

	form := &protocol.ButtonForm{}
	form.Tag = protocol.BUTTON_E
	form.Text = *text
	form.Params = params
	form.Type = category.String()
	form.Confirm = confirm
	form.SetAction(method, *protocol.NewMeta())

	return form
}

func NewOption(text protocol.TextForm, value string) protocol.OptionForm {
	form := protocol.OptionForm{}
	form.Text = text
	form.Value = value

	return form
}

// only be used for overflow elements
func NewJumpOption(text protocol.TextForm, url *string, multiURL *protocol.URLForm) protocol.OptionForm {
	form := protocol.OptionForm{}
	form.Text = text
	form.URL = url
	form.MultiURL = multiURL
	return form
}

func NewSelectStaticMenu(placeHolder *protocol.TextForm, params map[string]string,
	options []protocol.OptionForm, initOption *string, confirm *protocol.ConfirmForm, method string) *protocol.SelectorForm {

	form := &protocol.SelectorForm{}
	form.Tag = protocol.SELECT_STATIC_E
	form.Placeholder = placeHolder
	form.Options = options
	form.InitialOption = initOption
	form.Params = params
	form.Confirm = confirm
	form.SetAction(method, *protocol.NewMeta())

	return form
}

func NewSelectPersonMenu(placeHolder *protocol.TextForm, params map[string]string,
	options []protocol.OptionForm, initOption *string, confirm *protocol.ConfirmForm, method string) *protocol.SelectorForm {

	form := &protocol.SelectorForm{}
	form.Tag = protocol.SELECT_PERSON_E
	form.Placeholder = placeHolder
	form.Options = options
	form.InitialOption = initOption
	form.Params = params
	form.Confirm = confirm
	form.SetAction(method, *protocol.NewMeta())

	return form
}

// new element picker date
func NewPickerDate(placeHolder *protocol.TextForm, params map[string]string,
	confirm *protocol.ConfirmForm, initialDate *string, method string) *protocol.PickerDateForm {

	form := &protocol.PickerDateForm{}
	form.Tag = protocol.PICKERDATE_E
	form.Placeholder = placeHolder
	form.Params = params
	form.Confirm = confirm
	form.InitialDate = initialDate
	form.SetAction(method, *protocol.NewMeta())

	return form
}

// new element picker time
func NewPickerTime(placeHolder *protocol.TextForm, params map[string]string,
	confirm *protocol.ConfirmForm, initialTime *string, method string) *protocol.PickerTimeForm {
	form := &protocol.PickerTimeForm{}
	form.Tag = protocol.PICKERTIME_E
	form.Placeholder = placeHolder
	form.Params = params
	form.Confirm = confirm
	form.InitialTime = initialTime
	form.SetAction(method, *protocol.NewMeta())
	return form
}

// new element picker datetime
func NewPickerDatetime(placeHolder *protocol.TextForm, params map[string]string,
	confirm *protocol.ConfirmForm, initialDatetime *string, method string) *protocol.PickerDatetimeForm {
	form := &protocol.PickerDatetimeForm{}
	form.Tag = protocol.PICKERDATETIME_E
	form.Placeholder = placeHolder
	form.Params = params
	form.InitialDatetime = initialDatetime
	form.Confirm = confirm
	form.SetAction(method, *protocol.NewMeta())
	return form
}

func NewOverflowMenu(params map[string]string, options []protocol.OptionForm, confirm *protocol.ConfirmForm,
	method string) *protocol.OverflowForm {
	form := &protocol.OverflowForm{}
	form.Tag = protocol.OVERFLOW_E
	form.Params = params
	form.Options = options
	form.Confirm = confirm
	form.SetAction(method, *protocol.NewMeta())

	return form
}

func NewMultiPlatformURL(Url *string, AndroidUrl *string, IOSUrl *string, PCUrl *string) *protocol.URLForm {
	form := &protocol.URLForm{}
	form.Url = Url
	form.PCUrl = PCUrl
	form.IOSUrl = IOSUrl
	form.AndroidUrl = AndroidUrl

	return form
}

func NewField(short bool, text *protocol.TextForm) *protocol.FieldForm {
	form := &protocol.FieldForm{}
	form.Short = short
	form.Text = *text

	return form
}
