package commute

import (
	"MetaChat/pkg/qqbot_framework/commute/config"
	"MetaChat/pkg/qqbot_framework/commute/io"
	"MetaChat/pkg/qqbot_framework/commute/io/ws"
	"MetaChat/pkg/signal"
	"MetaChat/pkg/util/cq"
	"MetaChat/pkg/util/cq/condition"
	"context"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type QQBot struct {
	Context
	log         *zap.Logger
	config      *config.Config
	stopCh      chan bool
	connReadyCh <-chan bool
	IOHandler   io.IOHandler
	isReady     bool
}

func NewQQBot(log *zap.Logger, config *config.Config, handler io.IOHandler, stopHandler *signal.StopHandler) Context {
	bot := &QQBot{
		log:       log,
		config:    config,
		stopCh:    make(chan bool),
		IOHandler: handler,
		isReady:   false,
	}
	stopHandler.Add(bot)
	return bot
}

func (qq *QQBot) OnStart() {
	qq.connReadyCh = qq.IOHandler.GetOnReadyCh()
	qq.log.Info("framework started, waiting for connection ready")
	msgCh := qq.IOHandler.GetMessageCh()
	qq.initContext()
	go func() {
		for {
			select {
			case ready := <-qq.connReadyCh:
				if ready {
					qq.isReady = true
					qq.log.Info("connection ready!")
				} else {
					qq.isReady = false
				}
			case msg := <-msgCh:
				if qq.isReady {
					go qq.onMessage(msg)
					continue
				}
				qq.log.Info("receive message but bot not ready yet", zap.Any("msg", msg))
			case <-qq.stopCh:
				qq.notifyStop()
				break
			}
		}
	}()

}

func (qq *QQBot) OnStop() error {
	qq.stopCh <- true
	return nil
}

func (qq *QQBot) sendMessage(msg cq.CQResponse) {
	qq.IOHandler.SendMessage(msg)
}

func (qq *QQBot) onMessage(msg gjson.Result) {
	for _, eh := range EventHandlers {
		if eh != nil {
			go eh(qq, msg)
		}
	}
	ConditionHandlers.Range(func(key, value interface{}) bool {
		if key.(*condition.Condition).Fit(msg) {
			for _, eh := range value.([]ConditionHandler) {
				if eh != nil {
					go eh(qq, msg)
				}
			}
		}
		return true
	})
}

func (qq *QQBot) notifyStop() {
	//TODO: Do something when bot is stopped
	//qq.sendMessage(cq.CQResp(cq.ACTION_SEND_MESSAGE, cq.NewCQPrivateMessage()))
}

func (qq *QQBot) initContext() {
	qq.Context = NewCtx(qq)
}

func Provide() fx.Option {
	return fx.Options(fx.Provide(
		NewQQBot,
		ws.NewWS,
		config.NewConfig,
	),
		fx.Invoke(lc),
	)
}

func lc(lifecycle fx.Lifecycle, bot Context) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			bot.OnStart()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
