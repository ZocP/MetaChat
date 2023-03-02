package qq

import (
	"MetaChat/pkg/cqhttp/command"
	"MetaChat/pkg/gpt"
	"MetaChat/pkg/util/cq"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (qq *QQ) handleMessage(msg gjson.Result) {
	raw := msg.Get(cq.RAW_MESSAGE).String()

	if command.IsCommand(raw) {
		qq.log.Debug("收到一条命令：", zap.Any("command", command.UnpackCommand(raw)),
			zap.Any("type", msg.Get(cq.MESSAGE_TYPE).String()),
			zap.Any("from user", msg.Get(cq.USER_ID).String()),
		)
		qq.cmdHandler(msg, command.UnpackCommand(raw))
		return
	}
	if !qq.chat {
		return
	}
	if !cq.IsAtMe(msg.Get(cq.RAW_MESSAGE).String(), "3047430597") && msg.Get(cq.MESSAGE_TYPE).String() == cq.MESSAGE_TYPE_GROUP {
		return
	}

	req, err := gpt.SendReq(msg.Get(cq.RAW_MESSAGE).String())
	if err == nil && req != nil && len(req.Choices) > 0 {
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, req.Choices[0].Message.Content)))
	} else {
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, "聊天失败")))
	}

}

func (qq *QQ) cmdHandler(msg gjson.Result, cmd command.Command) {
	if handlers, ok := qq.cmdHandlers[cmd.Name]; ok {
		for _, handler := range handlers {
			handler(qq, msg, cmd)
			return
		}
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, "没有找到命令")))
	}
}
