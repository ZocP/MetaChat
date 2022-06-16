package qq

import (
	"MetaChat/app/metaChat/qq/condition"
	"MetaChat/pkg/cq"
	"MetaChat/pkg/qqBot/command"
	"MetaChat/pkg/saucenao"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type CMDHandler func(msg gjson.Result, cmd command.Command)

func (qq *QQ) echo(msg gjson.Result, cmd command.Command) {
	resp, echo := cq.GetCQRespEcho(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, cmd.RawParam))
	qq.SendMessage(resp)
	qq.SetAwaitMessage(echo)
	result := qq.AwaitMessage(echo)
	if result.Get(cq.STATUS).String() != cq.STATUS_OK {
		word := result.Get(cq.WORDING).String()
		qq.log.Info("echo失败", zap.String("echo", echo), zap.String("word", word))
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, result.Get(cq.WORDING).String())))
	}
}

func (qq *QQ) recognize(msg gjson.Result, cmd command.Command) {
	cdn := condition.NewCondition(
		map[string]string{
			cq.USER_ID: msg.Get(cq.USER_ID).String(),
		})
	qq.SetAwaitCondition(cdn)
	qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, "请发送一张图片")))
	result := qq.AwaitConditionResult(cdn)
	if cq.IsImageCQCode(result.Get(cq.RAW_MESSAGE).String()) {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, "识别中...")))
		qq.log.Info("识别到图片", zap.String("raw", result.Get(cq.MESSAGE).String()))
		result := saucenao.Recognize(cq.ParseCQCode(result.Get(cq.MESSAGE).String()).GetParam("url"))
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, result.String())))
	} else {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, "不是一张图片哦")))
	}

}
