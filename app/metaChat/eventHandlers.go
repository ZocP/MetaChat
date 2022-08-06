package metaChat

import (
	"MetaChat/pkg/util/cq"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (meta *MetaChat) handleCQMessage(msg gjson.Result) {
	meta.log.Debug("收到CQ消息", zap.Any("msg", msg))
	//meta.mc.SendMessage()
}

func (meta *MetaChat) handleMCMessage(msg gjson.Result) {
	meta.log.Debug("收到MC消息", zap.Any("msg", msg))
	meta.qq.SendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.GroupMessage(962403891, msg.String())))

}
