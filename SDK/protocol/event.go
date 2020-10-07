// Copyright (c) 2019 Bytedance Inc.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package protocol

const (
	// notification-type
	EventChallenge = "url_verification"
	EventCallback  = "event_callback"
)

const (
	// notification-subtype
	EventTypeApproval             = "approval"                                                         // notification--approval
	EventTypeLeaveApproval        = "leave_approval"                                                   // notification--leave_approval
	EventTypeWorkApproval         = "work_approval"                                                    // notification--work_approval
	EventTypeShiftApproval        = "shift_approval"                                                   // notification--shift_approval
	EventTypeRemedyApproval       = "remedy_approval"                                                  // notification--remedy_approval
	EventTypeTripApproval         = "trip_approval"                                                    // notification--trip_approval
	EventTypeAppOpen              = "app_open"                                                         // notification--app_open
	EventTypeContactUser          = "user_add,user_update,user_leave"                                  // notification--contact_user
	EventTypeContactDept          = "dept_add,dept_update,dept_delete"                                 // notification--contact_Department
	EventTypeContactScope         = "contact_scope_change"                                             // notification--contact_scope
	EventTypeMessage              = "message"                                                          // notification--message
	EventTypeRemoveBot            = "remove_bot"                                                       // notification--remove_bot
	EventTypeAddBot               = "add_bot"                                                          // notification--add_bot
	EventTypeP2PChatCreate        = "p2p_chat_create"                                                  // notification--p2p_chat_create
	EventTypeAppTicket            = "app_ticket"                                                       // notification--app_ticket
	EventTypeAppStatusChange      = "app_status_change"                                                // notification--app_status
	EventTypeUserToChat           = "add_user_to_chat,remove_user_from_chat,revoke_add_user_from_chat" // notification--add_user_to_chat
	EventTypeChatDisband          = "chat_disband"                                                     // notification--chat_disband
	EventTypeGroupSettingUpdate   = "group_setting_update"                                             // notification--group_setting_update
	EventTypeOrderPaid            = "order_paid"                                                       // notification--order_paid
	EventTypeCreateWidgetInstance = "create_widget_instance"                                           // notification--create_widget_instance
	EventTypeDeleteWidgetInstance = "delete_widget_instance"                                           // notification--delete_widget_instance
	EventTypeMessageRead          = "message_read"                                                     // notification--message_read
	EventTypeApprovalInstance     = "approval_instance"                                                // notification--approval_instance
	EventTypeAppUninstall         = "app_uninstalled"                                                  // notification--app_uninstalled
)

const (
	// message-type
	EventMsgTypeText         = "text"          // text
	EventMsgTypePost         = "post"          // rich_text/post
	EventMsgTypeImage        = "image"         // image
	EventMsgTypeMergeForward = "merge_forward" // merge forward
)

// common field
type CallbackBase struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	Ts        string `json:"ts"`
	Uuid      string `json:"uuid"`
}

type BaseEvent struct {
	Type      string `json:"type"`
	AppID     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
}

type ApprovalEvent struct {
	BaseEvent

	DefinitionCode string `json:"definition_code"`
	DefinitionName string `json:"definition_name"`
	InstanceCode   string `json:"instance_code"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	Event          string `json:"event"` // result: approve/reject/cancel
}

type LeaveApprovalEvent struct {
	BaseEvent

	InstanceCode   string `json:"instance_code"`
	EmployeeID     string `json:"employee_id"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	LeaveType      string `json:"leave_type"`
	LeaveUnit      int    `json:"leave_unit"`       // unit 1: harf day, 2: one day
	LeaveStartTime string `json:"leave_start_time"` // YYYY-MM-DD HH:MM:ss
	LeaveEndTime   string `json:"leave_end_time"`   // YYYY-MM-DD HH:MM:ss
	LeaveInterval  int    `json:"leave_interval"`   // unit seconds, 7200: 2 hours
	LeaveReason    string `json:"leave_reason"`
}

type WorkApprovalEvent struct {
	BaseEvent

	InstanceCode  string `json:"instance_code"`
	EmployeeID    string `json:"employee_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	WorkType      string `json:"work_type"`
	WorkStartTime string `json:"work_start_time"`
	WorkEndTime   string `json:"work_end_time"`
	WorkInterval  int    `json:"work_interval"` // unit seconds, 7200: 2 hours
	WorkReason    string `json:"work_reason"`
}

type ShiftApprovalEvent struct {
	BaseEvent

	InstanceCode string `json:"instance_code"`
	EmployeeID   string `json:"employee_id"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	ShiftTime    string `json:"shift_time"`
	ReturnTime   string `json:"return_time"`
	ShiftReason  string `json:"shift_reason"`
}

type RemedyApprovalEvent struct {
	BaseEvent

	InstanceCode string `json:"instance_code"`
	EmployeeID   string `json:"employee_id"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	RemedyTime   string `json:"remedy_time"`
	RemedyReason string `json:"remedy_reason"`
}

type TripApprovalEvent struct {
	BaseEvent

	InstanceCode string     `json:"instance_code"`
	EmployeeID   string     `json:"employee_id"`
	StartTime    int64      `json:"start_time"`
	EndTime      int64      `json:"end_time"`
	Schedules    []Schedule `json:"schedules"`
	TripInterval int        `json:"trip_interval"` // total duration, unit seconds
	TripReason   string     `json:"trip_reason"`
	TripPeers    []string   `json:"trip_peers"`
}

type Schedule struct {
	TripStartTime  string `json:"trip_start_time"` // YYYY-MM-DD HH:MM:ss
	TripEndTime    string `json:"trip_end_time"`   // YYYY-MM-DD HH:MM:ss
	TripInterval   int    `json:"trip_interval"`   // unit seconds
	Departure      string `json:"departure"`
	Destination    string `json:"destination"`
	Transportation string `json:"transportation"`
	TripType       string `json:"trip_type"`
	Remark         string `json:"remark"`
}

// notification--app open
type AppOpenEvent struct {
	BaseEvent

	Applications      []UserIDInfo `json:"applicants"`
	Installer         UserIDInfo   `json:"installer"`
	InstallerEmployee UserIDInfo   `json:"installer_employee"`
}
type UserIDInfo struct {
	OpenID string `json:"open_id"`
	UserID string `json:"user_id"`
}

// notification--contact user
type ContactUserEvent struct {
	BaseEvent

	OpenID     string `json:"open_id"`
	EmployeeID string `json:"employee_id"`
}

type ContactDeptEvent struct {
	BaseEvent

	OpenDepartmentID string `json:"open_department_id"`
}

type ContactScopeEvent struct {
	BaseEvent
}

// notification--message-text
type TextMsgEvent struct {
	BaseEvent

	MsgType          string `json:"msg_type"`
	RootID           string `json:"root_id"`
	ParentID         string `json:"parent_id"`
	OpenChatID       string `json:"open_chat_id"`
	ChatType         string `json:"chat_type"`
	OpenID           string `json:"open_id"`
	OpenMessageID    string `json:"open_message_id"`
	IsMention        bool   `json:"is_mention"`
	Text             string `json:"text"`
	TextWithoutAtBot string `json:"text_without_at_bot"`
}

// notification--message-richtext/post
type PostMsgEvent struct {
	BaseEvent

	MsgType          string   `json:"msg_type"`
	RootID           string   `json:"root_id"`
	ParentID         string   `json:"parent_id"`
	OpenChatID       string   `json:"open_chat_id"`
	ChatType         string   `json:"chat_type"`
	OpenID           string   `json:"open_id"`
	OpenMessageID    string   `json:"open_message_id"`
	IsMention        bool     `json:"is_mention"`
	Text             string   `json:"text"`
	TextWithoutAtBot string   `json:"text_without_at_bot"`
	Title            string   `json:"title"`
	ImageKeys        []string `json:"image_keys"`
}

// notification--message-image
type ImageMsgEvent struct {
	BaseEvent

	MsgType       string `json:"msg_type"`
	RootID        string `json:"root_id"`
	ParentID      string `json:"parent_id"`
	OpenChatID    string `json:"open_chat_id"`
	ChatType      string `json:"chat_type"`
	ImageHeight   string `json:"image_height"`
	ImageWidth    string `json:"image_width"`
	OpenID        string `json:"open_id"`
	OpenMessageID string `json:"open_message_id"`
	IsMention     bool   `json:"is_mention"`
	ImageKey      string `json:"image_key"`
	ImageUrl      string `json:"image_url"`
}

// notification--message-merge_forward
type MergeForwardMsgEvent struct {
	BaseEvent

	MsgType       string        `json:"msg_type"`
	RootID        string        `json:"root_id"`
	ParentID      string        `json:"parent_id"`
	OpenChatID    string        `json:"open_chat_id"`
	OpenID        string        `json:"open_id"`
	OpenMessageID string        `json:"open_message_id"`
	IsMention     bool          `json:"is_mention"`
	ChatType      string        `json:"chat_type"`
	ChartID       string        `json:"chat_id"`
	User          string        `json:"user"`
	MsgList       []interface{} `json:"msg_list"` // TextMsgEvent/PostMsgEvent/ImageMsgEvent
}

type RemoveBotEvent struct {
	BaseEvent

	ChatI18nNames       I18nNames `json:"chat_i18n_names"`
	ChatName            string    `json:"chat_name"`
	ChatOwnerEmployeeID string    `json:"chat_owner_employee_id"`
	ChatOwnerName       string    `json:"chat_owner_name"`
	ChatOwnerOpenID     string    `json:"chat_owner_open_id"`
	OpenChatID          string    `json:"open_chat_id"`
	OperatorEmployeeID  string    `json:"operator_employee_id"`
	OperatorName        string    `json:"operator_name"`
	OperatorOpenID      string    `json:"operator_open_id"`
	OwnerIsBot          bool      `json:"owner_is_bot"`
}

type I18nNames struct {
	EnUs string `json:"en_us"`
	JaJp string `json:"ja_jp"`
	ZhCn string `json:"zh_cn"`
}

type AddBotEvent struct {
	BaseEvent

	ChatI18nNames       I18nNames `json:"chat_i18n_names"`
	ChatName            string    `json:"chat_name"`
	ChatOwnerEmployeeID string    `json:"chat_owner_employee_id"`
	ChatOwnerName       string    `json:"chat_owner_name"`
	ChatOwnerOpenID     string    `json:"chat_owner_open_id"`
	OpenChatID          string    `json:"open_chat_id"`
	OperatorEmployeeID  string    `json:"operator_employee_id"`
	OperatorName        string    `json:"operator_name"`
	OperatorOpenID      string    `json:"operator_open_id"`
	OwnerIsBot          bool      `json:"owner_is_bot"`
}

type P2PChatCreateEvent struct {
	BaseEvent

	ChatID   string     `json:"chat_id"`
	Operator UserIDInfo `json:"operator"`
	User     struct {
		Name string `json:"name"`
		UserIDInfo
	} `json:"user"`
}

// notification--app_ticket
type AppTicketEvent struct {
	Type      string `json:"type"`
	AppID     string `json:"app_id"`
	AppTicket string `json:"app_ticket"`
}

const (
	StartByTenant  = "start_by_tenant"
	StopByTenant   = "stop_by_tenant"
	StopByPlatform = "stop_by_platform"
)

type AppStatusChangeEvent struct {
	BaseEvent

	Status string `json:"status"` // StartByTenant/StopByTenant/StopByPlatform
}

type UserToChatEvent struct {
	BaseEvent

	ChatID   string     `json:"chat_id"`
	Operator UserIDInfo `json:"operator"`
	Users    []struct {
		Name string `json:"name"`
		UserIDInfo
	} `json:"users"`
}

type ChatDisbandEvent struct {
	BaseEvent

	ChatID   string     `json:"chat_id"`
	Operator UserIDInfo `json:"operator"`
}

type GroupSettingUpdateEvent struct {
	BaseEvent

	ChatID   string     `json:"chat_id"`
	Operator UserIDInfo `json:"operator"`

	AfterChange struct {
		MessageNotification bool   `json:"message_notification"`  // case: change message notification
		AddMemberPermission string `json:"add_member_permission"` // case: change add-member permission
		OwnerOpenId         string `json:"owner_open_id"`         // case: change group owner
		OwnerUserId         string `json:"owner_user_id"`         // case: change group owner
	} `json:"after_change"`

	BeforeChange struct {
		MessageNotification bool   `json:"message_notification"`  // case: change message notification
		AddMemberPermission string `json:"add_member_permission"` // case: change add-member permission
		OwnerOpenId         string `json:"owner_open_id"`         // case: change group owner
		OwnerUserId         string `json:"owner_user_id"`         // case: change group owner
	} `json:"before_change"`
}

type OrderPaidEvent struct {
	BaseEvent

	OrderID       string `json:"order_id"`
	PricePlanID   string `json:"price_plan_id"`
	PricePlanType string `json:"price_plan_type"`
	Seats         int    `json:"seats"`
	BuyCount      int    `json:"buy_count"`
	CreateTime    string `json:"create_time"`
	PayTime       string `json:"pay_time"`
	BuyType       string `json:"buy_type"`
	SrcOrderID    string `json:"src_order_id"`
	OrderPayPrice int    `json:"order_pay_price"`
}

type CreateWidgetInstanceEvent struct {
	BaseEvent

	InstanceId []string `json:"instance_id"`
}

type DeleteWidgetInstanceEvent struct {
	BaseEvent

	InstanceId []string `json:"instance_id"`
}

type MessageReadEvent struct {
	BaseEvent

	OpenChatID     string   `json:"open_chat_id"`
	OpenID         string   `json:"open_id"`
	OpenMessageIDs []string `json:"open_message_ids"`
}

type ApprovalInstanceEvent struct {
	BaseEvent

	ApprovalCode string `json:"approval_code"`
	InstanceCode string `json:"instance_code"`
	Status       string `json:"status"`
	OperateTime  string `json:"operate_time"`
}

// url challenge
type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}

// AppUninstallEvent when company no longer uses the application
type AppUninstallEvent struct {
	BaseEvent
}
