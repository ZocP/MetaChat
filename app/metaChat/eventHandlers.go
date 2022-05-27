package metaChat

import (
	"MetaChat/app/metaChat/cq/group"
	"MetaChat/app/metaChat/eventBridge/request"
	"MetaChat/app/metaChat/eventBridge/response"
	"MetaChat/pkg/network"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"regexp"
)

func (meta *MetaChat) handleCQMessage(msg gjson.Result) {
	postType := msg.Get(request.POST_TYPE).String()
	// 优先处理有编号的信息
	f, ok := meta.echoHandlerCh[msg.Get(request.ECHO).String()]
	if ok {
		f <- msg
		return
	}
	if meta.botIsReady {
		switch postType {
		case request.POST_TYPE_MESSAGE:
			meta.onCQPostMsg(msg)
		case request.POST_TYPE_REQUEST:
			//to be implemented
		case request.POST_TYPE_NOTICE:
			meta.onCQPostNotice(msg)
		}
	}
}

func (meta *MetaChat) onCQPostMsg(msg gjson.Result) {
	// use for handling group and private message
	switch msg.Get(request.MESSAGE_TYPE).String() {
	case request.MESSAGE_TYPE_GROUP:
		meta.onCQPostMsgGroup(msg)
	case request.MESSAGE_TYPE_PRIVATE:
		meta.onCQPostMsgPrivate(msg)
	}
}

//处理私聊消息
func (meta *MetaChat) onCQPostMsgPrivate(msg gjson.Result) {
	meta.log.Info("receive private message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
}

//处理群消息
func (meta *MetaChat) onCQPostMsgGroup(msg gjson.Result) {
	meta.log.Info("receive group message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
	groupid := msg.Get(request.GROUP_ID).Int()
	user := msg.Get(request.USER_ID).Int()
	group := meta.qqBot.GetGroup(groupid)
	message := msg.Get(request.MESSAGE).String()
	compiler, err := regexp.Compile("^//")
	if err != nil {
		panic("compiler error")
	}
	if compiler.MatchString(message) {
		meta.onCommand(msg, user, group)
		return
	}

	compiler2, err := regexp.Compile("^/")
	if err != nil {
		panic("compiler2 error")
	}
	if compiler2.MatchString(message) {
		meta.onTransfer(msg, user, group)
		return
	}
}

//处理用户给机器人的指令
func (meta *MetaChat) onCommand(msg gjson.Result, user int64, group *group.Group) {
	switch msg.Get(request.MESSAGE).String() {
	case "//色图":
		meta.onRandomPic(msg, user, group)
	}
}

//处理命令转发
func (meta *MetaChat) onTransfer(msg gjson.Result, user int64, group *group.Group) {
	if !meta.cqHandler.IsAdmin(user) {
		meta.cqreplych <- response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetGroupMessage(
			group.GetID(),
			"你没有权限使用该命令,要添加管理员，请让管理员输入//addadmin [qq]",
		))
	}
	meta.mcreplych <- response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetGroupMessage(
		group.GetID(),
		msg.Get(request.MESSAGE).String(),
	))
}

func (meta *MetaChat) onMsgTransfer(msg gjson.Result, user int64, group *group.Group) {

}

//处理色图
func (meta *MetaChat) onRandomPic(msg gjson.Result, user int64, group *group.Group) {
	meta.log.Info("on random pic")
	result, err := network.GetFromUrlJSON("https://api.lolicon.app/setu/v2", map[string]string{"r18": "0"})
	if err != nil {
		meta.SendToQQ(response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetGroupMessage(group.GetID(), "获取涩图失败")))
		return
	}
	var (
		echo  string
		event response.CQResp
	)
	result.Get("data").ForEach(func(key, value gjson.Result) bool {
		event, echo = response.GetCQRespEcho(response.ACTION_SEND_MESSAGE, response.GetGroupMessage(group.GetID(), response.GetImageCQCode(value.Get("urls.original").String())))
		meta.registerEchoHandler(echo)
		meta.SendToQQ(event)
		return true
	})

	go func() {
		status := meta.waitForResult(echo)
		if status.Get(request.STATUS).String() == request.STATUS_ERROR {
			meta.SendToQQ(response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetGroupMessage(group.GetID(), "发送涩图失败，也许是太色了，请重试")))
		}
	}()
	//meta.SendToQQ(response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetGroupMessage(group.GetID(), response.GetImageCQCode(url))))
}

func (meta *MetaChat) onCQPostNotice(msg gjson.Result) {
	switch msg.Get(request.NOTICE_TYPE).String() {
	case request.NOTICE_TYPE_GROUP_INCREASE:
		meta.onGroupIncrease(msg)
	case request.NOTICE_TYPE_GROUP_DECREASE:
		meta.onGroupDecrease(msg)

	}
}

func (meta *MetaChat) onGroupIncrease(msg gjson.Result) {
	if msg.Get(request.SUB_TYPE).String() != request.SUB_TYPE_APPROVE {
		return
	}
	send, echo := response.GetCQRespEcho(response.ACTION_SEND_MESSAGE, response.GetGroupInfoMessage(msg.Get(request.GROUP_ID).Int()))
	meta.SendToQQ(send)
	meta.registerEchoHandler(echo)
	result := meta.waitForResult(echo)
	group := &group.Group{
		GroupName: result.Get(request.GROUP_NAME).String(),
		GroupID:   result.Get(request.GROUP_ID).Int(),
	}
	meta.qqBot.AddGroup(group)
	meta.log.Info("on group increase")
}

func (meta *MetaChat) onGroupDecrease(msg gjson.Result) {
	meta.log.Info("on group decrease")
}
