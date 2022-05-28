package cq

import (
	"strconv"
	"time"
)

const (
	ACTION_SEND_MESSAGE    = "send_msg"
	ACTION_GET_LOGIN_INFO  = "get_login_info"
	ACTION_GET_GROUP_LIST  = "get_group_list"
	ACTION_GET_FRIEND_LIST = "get_friend_list"
)

type CQResp struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

type CQNormalMessage struct {
	MessageType string `json:"message_type"`
	UserID      int64  `json:"user_id"`
	GroupID     int64  `json:"group_id"`
	Message     string `json:"message"`
	AutoEscape  bool   `json:"auto_escape"`
}

type CQGetGroupInfoMessage struct {
	GroupID int64 `json:"group_id"`
	NoCache bool  `json:"no_cache"`
}

func GetCQResp(action string, param interface{}) CQResp {
	return CQResp{
		Action: action,
		Params: param,
		//current time to string
		Echo: time.Now().Format("2006-01-02 15:04:05"),
	}
}

func GetCQRespEcho(action string, param interface{}) (CQResp, string) {
	t := strconv.FormatInt(time.Now().UnixMicro(), 10)
	return CQResp{
		Action: action,
		Params: param,
		//current time to string
		Echo: t,
	}, t
}

func GetNormalMessage(messageType string, userID int64, groupID int64, message string, autoEscape bool) CQNormalMessage {
	return CQNormalMessage{
		MessageType: messageType,
		UserID:      userID,
		GroupID:     groupID,
		Message:     message,
		AutoEscape:  autoEscape,
	}
}

func GetPrivateMessage(userID int64, message string) CQNormalMessage {
	return CQNormalMessage{
		MessageType: MESSAGE_TYPE_PRIVATE,
		UserID:      userID,
		Message:     message,
		AutoEscape:  false,
	}
}

func GetGroupMessage(groupID int64, message string) CQNormalMessage {
	return CQNormalMessage{
		MessageType: MESSAGE_TYPE_GROUP,
		GroupID:     groupID,
		Message:     message,
		AutoEscape:  false,
	}
}

func GetGroupInfoMessage(groupID int64) CQGetGroupInfoMessage {
	return CQGetGroupInfoMessage{
		GroupID: groupID,
		NoCache: true,
	}
}

func GetUserInfo(userID int64) CQResp {
	//TODO get user info
	//return GetCQResp(ACTION_GET_LOGIN_INFO, (userID))
	panic("not implemented")
}
