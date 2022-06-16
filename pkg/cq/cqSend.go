package cq

import (
	"github.com/tidwall/gjson"
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

func GetMessageQuick(msg gjson.Result, message string) CQNormalMessage {
	switch msg.Get(MESSAGE_TYPE).String() {
	case MESSAGE_TYPE_GROUP:
		return GetGroupMessage(msg.Get(GROUP_ID).Int(), message)
	case MESSAGE_TYPE_PRIVATE:
		return GetPrivateMessage(msg.Get(USER_ID).Int(), message)
	}
	panic("not implemented")
}

func GetMessageAt(id int64, message string, at string) CQNormalMessage {
	switch at {
	case MESSAGE_TYPE_GROUP:
		return GetGroupMessage(id, message)
	case MESSAGE_TYPE_PRIVATE:
		return GetPrivateMessage(id, message)
	}
	panic("not implemented")
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
