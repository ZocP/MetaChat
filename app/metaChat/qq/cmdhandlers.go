package qq

import (
	"MetaChat/pkg/cqhttp/command"
	"MetaChat/pkg/gpt"
	"MetaChat/pkg/util/saucenao"

	"MetaChat/pkg/util/cq"
	"MetaChat/pkg/util/cq/condition"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type CMDHandler func(qq *QQ, msg gjson.Result, cmd command.Command)

func Echo(qq *QQ, msg gjson.Result, cmd command.Command) {
	resp, echo := cq.CQRespEcho(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, cmd.RawParam))
	qq.SendMessage(resp)
	qq.SetAwaitMessage(echo)
	result, _ := qq.AwaitMessage(echo)
	if result.Get(cq.STATUS).String() != cq.STATUS_OK {
		word := result.Get(cq.WORDING).String()
		qq.log.Info("echo失败", zap.String("echo", echo), zap.String("word", word))
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, result.Get(cq.WORDING).String())))
	}
}

func Recognize(qq *QQ, msg gjson.Result, cmd command.Command) {
	cdn := condition.NewCondition(
		map[string]string{
			cq.USER_ID: msg.Get(cq.USER_ID).String(),
		})
	qq.SetAwaitCondition(cdn)
	qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, "请发送一张图片")))
	result, isTimeout := qq.AwaitConditionResult(cdn)
	if isTimeout {
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, msg.Get(cq.USER_ID).String()+"的操作超时")))
		return
	}
	if cq.IsImageCQCode(result.Get(cq.RAW_MESSAGE).String()) {
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, "识别中...")))
		qq.log.Info("识别到图片", zap.String("raw", result.Get(cq.MESSAGE).String()))
		result := saucenao.Recognize(cq.ParseCQCode(result.Get(cq.MESSAGE).String()).GetParam("url"))
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, saucenao.GetQQSauceNaoResult(result))))
	} else {
		qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, "不是一张图片哦")))
	}
}

func Chat(qq *QQ, msg gjson.Result, cmd command.Command) {
	qq.Lock()
	qq.chat = !qq.chat
	lang := ""
	if qq.chat {
		lang = "聊天已开启"
	} else {
		lang = "聊天已关闭"
	}
	gpt.SetConfig(qq.viper)
	qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.CQMessageQuick(msg, lang)))
	qq.Unlock()
}
