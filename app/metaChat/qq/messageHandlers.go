package qq

import (
	"MetaChat/pkg/cq"
	"MetaChat/pkg/qqBot"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type QQ struct {
	qqBot.Context
	log *zap.Logger

	throwCh chan gjson.Result
}

func (qq *QQ) MessageHandler(ctx qqBot.Context, msg gjson.Result) {
	qq.log.Debug("收到CQ消息", zap.Any("msg", msg))
	if msg.Get(cq.POST_TYPE).String() == cq.POST_TYPE_MESSAGE {
		ctx.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetGroupMessage(msg.Get(cq.GROUP_ID).Int(), msg.Get(cq.RAW_MESSAGE).String())))
	}
}

func (qq *QQ) throw(result gjson.Result) {
	qq.throwCh <- result
}

func (qq *QQ) GetThrow() <-chan gjson.Result {
	return qq.throwCh
}

func NewQQ(log *zap.Logger, bot *qqBot.QQBot) *QQ {
	return &QQ{
		log:     log,
		Context: bot,
	}
}

func (qq *QQ) onStart() {
	for {
		select {}
	}
}

func Provide() fx.Option {
	return fx.Provide(NewQQ)
}
