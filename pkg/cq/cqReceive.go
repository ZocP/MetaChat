package cq

const (
	TIME    = "time"
	SELF_ID = "self_id"

	POST_TYPE         = "post_type"
	POST_TYPE_MESSAGE = "message"
	POST_TYPE_REQUEST = "request"
	POST_TYPE_NOTICE  = "notice"

	META_EVENT_TYPE           = "meta_event_type"
	META_EVENT_TYPE_HEARTBEAT = "heartbeat"
	META_EVENT_TYPE_LIFECYCLE = "lifecycle"

	MESSAGE_TYPE         = "message_type"
	MESSAGE_TYPE_GROUP   = "group"
	MESSAGE_TYPE_PRIVATE = "private"

	MESSAGE     = "message"
	RAW_MESSAGE = "raw_message"

	DATA     = "data"
	NICKNAME = "nickname"

	GROUP_ID               = "group_id"
	GROUP_NAME             = "group_name"
	GROUP_MEMO             = "group_memo"
	GROUP_CREATE_TIME      = "group_create_time"
	GROUP_MEMBER_COUNT     = "member_count"
	GROUP_MAX_MEMBER_COUNT = "group_max_member_count"
	USER_ID                = "user_id"
	ECHO                   = "echo"
	STATUS                 = "status"
	STATUS_OK              = "ok"
	STATUS_ERROR           = "failed"

	SUB_TYPE         = "sub_type"
	SUB_TYPE_kiCK_ME = "kick_me"
	SUB_TYPE_APPROVE = "approve"

	NOTICE_TYPE                = "notice_type"
	NOTICE_TYPE_GROUP_INCREASE = "group_increase"
	NOTICE_TYPE_GROUP_DECREASE = "group_reduce"

	WORDING = "wording"

	SENDER        = "sender"
	SENDER_USERID = SENDER + "." + USER_ID
)
