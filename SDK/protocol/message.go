// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

type UserInfo struct {
	ID   string
	Type UserType
}

type UserType int64

const (
	UserTypeUnknown UserType = iota
	UserTypeOpenID
	UserTypeUserID
	UserTypeEmail
	UserTypeChatID
)

type Language int64

const (
	LanguageUnknown Language = iota
	ZhCN
	EnUS
	JaJP
)

func (i Language) String() string {
	switch i {
	case ZhCN:
		return "zh_cn"
	case EnUS:
		return "en_us"
	case JaJP:
		return "ja_jp"
	default:
		return "UNKNOWN"
	}
}

type MessageType string

const (
	TextMsgType      MessageType = "text"
	ImageMsgType     MessageType = "image"
	PostMsgType      MessageType = "post"
	ShareChatMsgType MessageType = "share_chat"
	CardMsgType      MessageType = "interactive"
)

// message base info
type BaseInfo struct {
	// chat_id > open_id > user_id > email

	// Send message to the user
	OpenID string `json:"open_id,omitempty" validate:"omitempty"`
	UserID string `json:"user_id,omitempty" validate:"omitempty"`
	Email  string `json:"email,omitempty" validate:"omitempty"`

	// send group message
	ChatID string `json:"chat_id,omitempty" validate:"omitempty"`

	// replay message
	RootID string `json:"root_id,omitempty" validate:"omitempty"`
}

func (b *BaseInfo) SetBaseInfo(userInfo *UserInfo, rootID string) {
	b.RootID = rootID

	switch userInfo.Type {
	case UserTypeOpenID:
		b.OpenID = userInfo.ID
	case UserTypeUserID:
		b.UserID = userInfo.ID
	case UserTypeEmail:
		b.Email = userInfo.ID
	case UserTypeChatID:
		b.ChatID = userInfo.ID
	}
}

// message batch base info
type BatchBaseInfo struct {
	DepartmentIDs []string `json:"department_ids,omitempty" validate:"omitempty"`
	OpenIDs       []string `json:"open_ids,omitempty" validate:"omitempty"`
	UserIDs       []string `json:"user_ids,omitempty" validate:"omitempty"`
}

type MessageContent struct {
	Text            string                   `json:"text,omitempty" validate:"omitempty"`
	Title           string                   `json:"title,omitempty" validate:"omitempty"`
	ImageKey        string                   `json:"image_key,omitempty" validate:"omitempty"`
	Post            map[string]*RichTextForm `json:"post,omitempty" validate:"omitempty"`
	ShareOpenChatID string                   `json:"share_open_chat_id,omitempty" validate:"omitempty"`
}

type RichTextForm struct {
	Title   string           `json:"title" validate:"omitempty"`
	Content *RichTextContent `json:"content" validate:"omitempty"`
}

type RichTextContent [][]RichTextElementForm

func (r *RichTextContent) AddElementBlock(elements ...*RichTextElementForm) *RichTextContent {

	var richTextElementForms []RichTextElementForm
	for _, v := range elements {
		richTextElementForms = append(richTextElementForms, *v)
	}

	*r = append(*r, richTextElementForms)
	return r
}

type RichTextElementForm struct {
	Tag      string  `json:"tag,omitempty" validate:"omitempty"`
	Text     *string `json:"text,omitempty" validate:"omitempty"`
	Lines    *int32  `json:"lines,omitempty" validate:"omitempty,min=1,max=100"`
	UnEscape bool    `json:"un_escape,omitempty" validate:"omitempty"`
	//--------------------A-----------------------
	Href string `json:"href,omitempty" validate:"omitempty"`
	//--------------------Image-----------------------
	ImageKey string `json:"image_key,omitempty" validate:"omitempty"`
	Height   int32  `json:"height,omitempty" validate:"omitempty"`
	Width    int32  `json:"width,omitempty" validate:"omitempty"`
	//--------------------At------------------------
	UserID string `json:"user_id,omitempty" validate:"omitempty"`
}

type UpLoadImageResponse struct {
	BaseResponse
	Data struct {
		ImageKey string `json:"image_key,omitempty" validate:"omitempty"`
		Url      string `json:"url,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type SendMsgRequest struct {
	BaseInfo

	MsgType string         `json:"msg_type,omitempty" validate:"required"`
	Content MessageContent `json:"content,omitempty" validate:"required"`
}

type SendMsgResponse struct {
	BaseResponse

	Data struct {
		MessageID string `json:"message_id,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type SendMsgBatchRequest struct {
	BatchBaseInfo

	MsgType string         `json:"msg_type,omitempty" validate:"required"`
	Content MessageContent `json:"content,omitempty" validate:"required"`
}

type SendMsgBatchResponse struct {
	BaseResponse

	Data struct {
		MessageID            string   `json:"message_id,omitempty" validate:"omitempty"`
		InvalidDepartmentIDs []string `json:"invalid_department_ids,omitempty" validate:"omitempty"`
		InvalidOpenIDs       []string `json:"invalid_open_ids,omitempty" validate:"omitempty"`
		InvalidUserIDs       []string `json:"invalid_user_ids,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type SendCardMsgRequest struct {
	BaseInfo

	UUID        string   `json:"uuid,omitempty" validate:"omitempty"`
	MsgType     string   `json:"msg_type,omitempty" validate:"required"`
	Card        CardForm `json:"card,omitempty" validate:"omitempty"`
	UpdateMulti bool     `json:"update_multi" validate:"omitempty"`
}

func (s *SendCardMsgRequest) SetUUID(UUID string) {
	s.UUID = UUID
}

type SendCardMsgResponse struct {
	BaseResponse

	Data struct {
		MessageID string `json:"message_id,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

type SendCardMsgBatchRequest struct {
	BatchBaseInfo

	MsgType     string   `json:"msg_type,omitempty" validate:"required"`
	Card        CardForm `json:"card,omitempty" validate:"omitempty"`
	UpdateMulti bool     `json:"update_multi" validate:"omitempty"`
}

type SendCardMsgBatchResponse struct {
	BaseResponse

	Data struct {
		MessageID            string   `json:"message_id,omitempty" validate:"omitempty"`
		InvalidDepartmentIDs []string `json:"invalid_department_ids,omitempty" validate:"omitempty"`
		InvalidOpenIDs       []string `json:"invalid_open_ids,omitempty" validate:"omitempty"`
		InvalidUserIDs       []string `json:"invalid_user_ids,omitempty" validate:"omitempty"`
	} `json:"data,omitempty" validate:"omitempty"`
}

func NewTextMsgReq(user *UserInfo, rootID string, text string) *SendMsgRequest {
	request := &SendMsgRequest{
		MsgType: string(TextMsgType),
		Content: MessageContent{Text: text},
	}
	request.SetBaseInfo(user, rootID)

	return request
}

func NewImageMsgReq(user *UserInfo, rootID string, imageKey string) *SendMsgRequest {
	request := &SendMsgRequest{
		MsgType: string(ImageMsgType),
		Content: MessageContent{ImageKey: imageKey},
	}
	request.SetBaseInfo(user, rootID)

	return request
}

func NewPostMsgReq(user *UserInfo, rootID string, post map[string]*RichTextForm) *SendMsgRequest {
	request := &SendMsgRequest{
		MsgType: string(PostMsgType),
		Content: MessageContent{Post: post},
	}
	request.SetBaseInfo(user, rootID)

	return request
}

func NewShareChatMsgReq(user *UserInfo, rootID string, shareChatID string) *SendMsgRequest {
	request := &SendMsgRequest{
		MsgType: string(ShareChatMsgType),
		Content: MessageContent{ShareOpenChatID: shareChatID},
	}
	request.SetBaseInfo(user, rootID)

	return request
}

func NewCardMsgReq(user *UserInfo, rootID string, card CardForm, updateMulti bool) *SendCardMsgRequest {
	request := &SendCardMsgRequest{
		MsgType:     string(CardMsgType),
		Card:        card,
		UpdateMulti: updateMulti,
	}
	request.SetBaseInfo(user, rootID)

	return request
}

func NewBatchTextMsgReq(info *BatchBaseInfo, rootID string, text string) *SendMsgBatchRequest {
	request := &SendMsgBatchRequest{
		MsgType: string(TextMsgType),
		Content: MessageContent{Text: text},
	}

	request.DepartmentIDs = info.DepartmentIDs
	request.OpenIDs = info.OpenIDs
	request.UserIDs = info.UserIDs

	return request
}

func NewBatchImageMsgReq(info *BatchBaseInfo, rootID string, imageKey string) *SendMsgBatchRequest {
	request := &SendMsgBatchRequest{
		MsgType: string(ImageMsgType),
		Content: MessageContent{ImageKey: imageKey},
	}

	request.DepartmentIDs = info.DepartmentIDs
	request.OpenIDs = info.OpenIDs
	request.UserIDs = info.UserIDs

	return request
}

func NewBatchPostMsgReq(info *BatchBaseInfo, rootID string, post map[string]*RichTextForm) *SendMsgBatchRequest {
	request := &SendMsgBatchRequest{
		MsgType: string(PostMsgType),
		Content: MessageContent{Post: post},
	}

	request.DepartmentIDs = info.DepartmentIDs
	request.OpenIDs = info.OpenIDs
	request.UserIDs = info.UserIDs

	return request
}

func NewBatchShareChatMsgReq(info *BatchBaseInfo, rootID string, shareChatID string) *SendMsgBatchRequest {
	request := &SendMsgBatchRequest{
		MsgType: string(ShareChatMsgType),
		Content: MessageContent{ShareOpenChatID: shareChatID},
	}

	request.DepartmentIDs = info.DepartmentIDs
	request.OpenIDs = info.OpenIDs
	request.UserIDs = info.UserIDs

	return request
}

func NewBatchCardMsgReq(info *BatchBaseInfo, rootID string, card CardForm, updateMulti bool) *SendCardMsgBatchRequest {
	request := &SendCardMsgBatchRequest{
		MsgType:     string(CardMsgType),
		Card:        card,
		UpdateMulti: updateMulti,
	}

	request.DepartmentIDs = info.DepartmentIDs
	request.OpenIDs = info.OpenIDs
	request.UserIDs = info.UserIDs

	return request
}
