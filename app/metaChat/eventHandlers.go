package metaChat

import (
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (meta *MetaChat) handleCQMessage(msg gjson.Result) {
	meta.log.Debug("收到CQ消息", zap.Any("msg", msg))
	//meta.mc.SendMessage()
}

