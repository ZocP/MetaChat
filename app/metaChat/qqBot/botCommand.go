package qqBot

import (
	"MetaChat/app/metaChat/qqBot/commands"
	"MetaChat/pkg/cq"
	"MetaChat/pkg/lolicon"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (qq *QQBot) onBotCommand(msg gjson.Result) {
	message := msg.Get(cq.MESSAGE).String()
	cmd := commands.UnpackCommand(message)
	qq.log.Debug("unpack command", zap.Any("cmd", cmd))
	switch cmd.Name {
	case "色图":
		qq.onRandomPics(msg, cmd)
	}
}

func (qq *QQBot) onRandomPics(msg gjson.Result, cmd commands.Command) {
	var (
		result gjson.Result
		err    error
	)
	if !cmd.HasParam {
		result, err = lolicon.GetRandomPictureJSON(lolicon.Param{
			Adult: lolicon.NOTR18,
		})
	} else {
		param, err := lolicon.ParseParam(cmd.Param)
		if err != nil {
			qq.log.Error("error while parsing param", zap.Error(err))
			qq.sendGroupError(msg.Get(cq.GROUP_ID).Int(), ParamError)
			return
		}
		result, err = lolicon.GetRandomPictureJSON(param)
	}
	groupID := msg.Get(cq.GROUP_ID).Int()
	qq.log.Info("getting random pictures", zap.Any("param", cmd.Param), zap.Any("result", result))
	if err != nil {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, "获取涩图失败")))
		return
	}

	if len(result.Get("data").Array()) == 0 {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, "没有找到涩图哦")))
		return
	}
	result.Get("data").ForEach(func(key, value gjson.Result) bool {
		event, echo := cq.GetCQRespEcho(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, cq.GetImageCQCode(value.Get("urls.original").String())))
		qq.RegisterEchoHandler(echo)
		qq.SendMessage(event)
		go func() {
			status := qq.WaitForResult(echo)
			if status.Get(cq.STATUS).String() == cq.STATUS_ERROR {
				qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, "发送涩图失败，也许是太色了，请重试")))
			}
		}()
		return true
	})
}
