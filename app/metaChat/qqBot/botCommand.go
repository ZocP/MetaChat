package qqBot

import (
	"MetaChat/app/metaChat/qqBot/commands"
	"MetaChat/pkg/cq"
	"MetaChat/pkg/lolicon"
	"fmt"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (qq *QQBot) onBotCommand(msg gjson.Result, at string) {
	message := msg.Get(cq.MESSAGE).String()
	cmd := commands.UnpackCommand(message)
	qq.log.Debug("unpack command", zap.Any("cmd", cmd))
	//switch cmd.Name{
	//	case "色图":
	//	qq.onRandomPics(msg, cmd, at)
	//}
	switch at {
	case cq.MESSAGE_TYPE_GROUP:
		switch cmd.Name {
		case "色图":
			qq.onRandomPics(msg, cmd, at)

		}
	case cq.MESSAGE_TYPE_PRIVATE:
		switch cmd.Name {
		case "色图":
			qq.onRandomPics(msg, cmd, at)
		}

	}

}

func (qq *QQBot) onRandomPics(msg gjson.Result, cmd commands.Command, at string) {
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
		if qq.IsAdmin(msg.Get(cq.USER_ID).Int()) {
			param, err = lolicon.ParseParamAll(cmd.Param)
		}
		if err != nil {
			qq.log.Error("error while parsing param", zap.Error(err))
			qq.sendErrorAt(msg.Get(cq.GROUP_ID).Int(), ParamError, at)
			return
		}
		result, err = lolicon.GetRandomPictureJSON(param)
	}
	var ID int64
	if at == cq.MESSAGE_TYPE_GROUP {
		ID = msg.Get(cq.GROUP_ID).Int()
	} else {
		ID = msg.Get(cq.USER_ID).Int()
	}
	qq.log.Info("getting random pictures", zap.Any("param", cmd.Param), zap.Any("result", result))
	if err != nil {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageAt(ID, "获取涩图失败", at)))
		return
	}

	if len(result.Get("data").Array()) == 0 {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageAt(ID, "没有找到涩图哦", at)))
		return
	}
	result.Get("data").ForEach(func(key, value gjson.Result) bool {
		link := value.Get("urls.original").String()
		cqCode := cq.GetImageCQCode(link)
		event, echo := cq.GetCQRespEcho(cq.ACTION_SEND_MESSAGE, cq.GetMessageAt(ID, cqCode, at))
		qq.RegisterEchoHandler(echo)
		qq.SendMessage(event)
		qq.log.Debug("cq code sent", zap.Any("cqCode", cqCode))
		status := qq.WaitForResult(echo)
		if status.Get(cq.STATUS).String() == cq.STATUS_ERROR {
			fmtmsg := fmt.Sprintf("发送涩图失败，也许是太色了，但是可以从下面的链接访问\n%s", link)
			qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageAt(ID, fmtmsg, at)))
			qq.log.Debug("error while sending message", zap.Any("status", status))
		}
		return true
	})
}
