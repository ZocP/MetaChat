package metaChat

import (
	"MetaChat/app/metaChat/cq/group"
	"MetaChat/app/metaChat/eventBridge/request"
	"MetaChat/app/metaChat/eventBridge/response"
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
			go meta.handleCQPostMsg(msg)
		case request.POST_TYPE_REQUEST:
			//to be implemented
		case request.POST_TYPE_NOTICE:
			//to be implemented
		}
	}
}

func (meta *MetaChat) handleCQPostMsg(msg gjson.Result) {
	// use for handling group and private message
	switch msg.Get(request.MESSAGE_TYPE).String() {
	case request.MESSAGE_TYPE_GROUP:
		meta.handleCQPostMsgGroup(msg)
	case request.MESSAGE_TYPE_PRIVATE:
		meta.handleCQPostMsgPrivate(msg)
	}
}

func (meta *MetaChat) handleCQPostMsgPrivate(msg gjson.Result) {
	meta.log.Info("receive private message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
}

func (meta *MetaChat) handleCQPostMsgGroup(msg gjson.Result) {
	meta.log.Info("receive group message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
	groupid := msg.Get(request.GROUP_ID).Int()
	//mode := meta.cqHandler.GetGroupMode(groupid)
	//if mode == group.MODE_REPEAT{
	//	meta.cqreplych <- response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetNormalMessage(
	//		request.MESSAGE_TYPE_GROUP,
	//		0,
	//		msg.Get(request.GROUP_ID).Int(),
	//		msg.Get(request.MESSAGE).String(),
	//		false,
	//	))
	//}
	user := msg.Get(request.USER_ID).Int()
	group := meta.qqBot.GetGroup(groupid)
	message := msg.Get(request.MESSAGE).String()
	compiler, err := regexp.Compile("^//")
	if err != nil {
		panic("compiler error")
	}
	if compiler.MatchString(message) {
		meta.handleCommand(msg, user, group)
		return
	}

	compiler2, err := regexp.Compile("^/")
	if err != nil {
		panic("compiler2 error")
	}
	if compiler2.MatchString(message) {
		meta.handleTransfer(msg, user, group)
		return
	}
}

//处理用户给机器人的指令

func (meta *MetaChat) handleCommand(msg gjson.Result, user int64, group *group.Group) {

}

//处理命令转发

func (meta *MetaChat) handleTransfer(msg gjson.Result, user int64, group *group.Group) {
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

func (meta *MetaChat) handleMsgTransfer(msg gjson.Result, user int64, group *group.Group) {

}
