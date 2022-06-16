package qqBot

import (
	"MetaChat/pkg/cq"
	"MetaChat/pkg/qqBot/config"
	"MetaChat/pkg/qqBot/io"
	"MetaChat/pkg/qqBot/io/ws"
	"MetaChat/pkg/signal"
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

func NewQQBot(log *zap.Logger, config *config.Config, handler io.IOHandler, stopHandler *signal.StopHandler) *QQBot {
	bot := &QQBot{
		log:       log,
		config:    config,
		stopCh:    make(chan bool),
		IOHandler: handler,
		isReady:   false,
	}
	bot.OnStart()
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
					qq.onMessage(msg)
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

func (qq *QQBot) sendMessage(msg cq.CQResp) {
	qq.IOHandler.SendMessage(msg)
}

func (qq *QQBot) onMessage(msg gjson.Result) {
	for _, eh := range EventHandlers {
		if eh != nil {
			go eh(qq, msg)
		}
	}
}

func (qq *QQBot) notifyStop() {
	//TODO: Do something when bot is stopped
	//qq.sendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetPrivateMessage()))
}

func (qq *QQBot) initContext() {
	qq.Context = NewCtx(qq)
}

func Provide() fx.Option {
	return fx.Options(fx.Provide(
		NewQQBot,
		ws.NewWS,
		config.NewConfig,
	))
}
