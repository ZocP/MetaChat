package qqBot

import (
	"MetaChat/app/metaChat/qqBot/commands"
	"MetaChat/pkg/cq"
	"MetaChat/pkg/network"
	"github.com/tidwall/gjson"
)

func (qq *QQBot) onBotCommand(msg gjson.Result) {
	message := msg.Get(cq.MESSAGE).String()
	cmd := commands.UnpackCommand(message)
	switch cmd.Name {
	case "色图":
		qq.onRandomPics(msg, cmd)
	}
}

func (qq *QQBot) onRandomPics(msg gjson.Result, cmd commands.Command) {
	groupID := msg.Get(cq.GROUP_ID).Int()
	qq.log.Info("on random pic")
	result, err := network.GetFromUrlJSON("https://api.lolicon.app/setu/v2", map[string]string{"r18": "0"})
	if err != nil {
		qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, "获取涩图失败")))
		return
	}
	var (
		echo  string
		event cq.CQResp
	)
	result.Get("data").ForEach(func(key, value gjson.Result) bool {
		event, echo = cq.GetCQRespEcho(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, cq.GetImageCQCode(value.Get("urls.original").String())))
		qq.RegisterEchoHandler(echo)
		qq.SendMessage(event)
		return true
	})

	go func() {
		status := qq.WaitForResult(echo)
		if status.Get(cq.STATUS).String() == cq.STATUS_ERROR {
			qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(groupID, "发送涩图失败，也许是太色了，请重试")))
		}
	}()
}
