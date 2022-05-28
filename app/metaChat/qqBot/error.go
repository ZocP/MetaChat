package qqBot

import "MetaChat/pkg/cq"

const (
	ParamError    = "参数错误，请检查参数"
	NotFoundError = "未找到该指令"
	NotAdminError = "您不是管理员，无法使用该指令"
)

func (qq *QQBot) sendGroupError(group int64, error string) {
	qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(group, error)))
}
