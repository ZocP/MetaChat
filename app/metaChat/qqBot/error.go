package qqBot

import "MetaChat/pkg/cq"

const (
	ParamError    = "参数错误，请检查参数"
	NotFoundError = "未找到该指令"
	NotAdminError = "您不是管理员，无法使用该指令"
)

func (qq *QQBot) sendErrorAt(id int64, error string, at string) {
	switch at {
	case cq.MESSAGE_TYPE_GROUP:
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(id, error)))
	case cq.MESSAGE_TYPE_PRIVATE:
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetPrivateMessage(id, error)))
	}

}
