// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package message

import (
	"github.com/larksuite/botframework-go/SDK/protocol"
)

func NewRichTextForm(title *string, content *protocol.RichTextContent) *protocol.RichTextForm {
	richForm := &protocol.RichTextForm{}
	richForm.Title = *title
	richForm.Content = content
	return richForm
}

func NewRichTextContent() *protocol.RichTextContent {
	return &protocol.RichTextContent{}
}

func NewRichTextElementForm() *protocol.RichTextElementForm {
	return &protocol.RichTextElementForm{}
}

func NewTextTag(text string, unURLEncode bool, lines int32) *protocol.RichTextElementForm {
	var textTag protocol.RichTextElementForm
	textTag.Tag = "text"
	textTag.Text = &text
	textTag.UnEscape = unURLEncode
	textTag.Lines = &lines
	return &textTag
}

func NewATag(text string, unURLEncode bool, href string) *protocol.RichTextElementForm {
	var textTag protocol.RichTextElementForm
	textTag.Tag = "a"
	textTag.Text = &text
	textTag.UnEscape = unURLEncode
	textTag.Href = href
	return &textTag
}

func NewAtTag(text, userID string) *protocol.RichTextElementForm {
	var atTag protocol.RichTextElementForm
	atTag.Tag = "at"
	atTag.Text = &text
	atTag.UserID = userID
	return &atTag
}

// go can get imageKey by call the function "GetImageKey"
func NewImageTag(imageKey string, height int32, width int32) *protocol.RichTextElementForm {

	var imageTag protocol.RichTextElementForm
	imageTag.Tag = "img"
	imageTag.ImageKey = imageKey
	imageTag.Height = height
	imageTag.Width = width
	return &imageTag
}

func checkPostContent(content map[protocol.Language]*protocol.RichTextForm) bool {
	for _, v := range content {
		if !checkRichTextContent(*v) {
			return false
		}
	}
	return true
}

func checkRichTextContent(richTextForm protocol.RichTextForm) bool {
	for _, RichTextElementForms := range *richTextForm.Content {
		for _, richTextElementForm := range RichTextElementForms {
			if richTextElementForm.Tag == "img" && richTextElementForm.ImageKey == "" {
				return false
			}
		}
	}
	return true
}
