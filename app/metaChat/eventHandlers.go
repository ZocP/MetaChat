package metaChat

import (
	"MetaChat/app/metaChat/eventBridge/request"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (meta *MetaChat) handleCQMessage(msg gjson.Result) {

	postType := msg.Get(request.POST_TYPE).String()
	switch postType {
	case request.POST_TYPE_MESSAGE:
		meta.handleCQPostMsg(msg)
	case request.POST_TYPE_REQUEST:
		//to be implemented
	case request.POST_TYPE_NOTICE:
		//to be implemented
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
	group := meta.cqHandler.GetGroup(groupid)
	message := msg.Get(request.MESSAGE).String()

}

func (meta *MetaChat) handleCQPostMsgPrivate(msg gjson.Result) {
	meta.log.Info("receive private message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
}
