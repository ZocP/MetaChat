package qq

import (
	"MetaChat/pkg/cqhttp/command"
	"MetaChat/pkg/cqhttp/commute"
	"MetaChat/pkg/util/cq"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (qq *QQ) handleMessage(ctx commute.Context, msg gjson.Result) {
	raw := msg.Get(cq.RAW_MESSAGE).String()
	if command.IsCommand(raw) {
		qq.log.Debug("收到一条命令：", zap.Any("command", command.UnpackCommand(raw)),
			zap.Any("type", msg.Get(cq.MESSAGE_TYPE).String()),
			zap.Any("from user", msg.Get(cq.USER_ID).String()),
		)
		qq.cmdHandler(ctx, msg, command.UnpackCommand(raw))
	}
}

func (qq *QQ) handleRequest(ctx commute.Context, msg gjson.Result) {

}

func (qq *QQ) cmdHandler(ctx commute.Context, msg gjson.Result, cmd command.Command) {
	if handlers, ok := qq.cmdHandlers[cmd.Name]; ok {
		for _, handler := range handlers {
			handler(msg, cmd)
		}
		return
	}
	qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, "没有找到命令")))
}
