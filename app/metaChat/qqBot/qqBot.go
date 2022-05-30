package qqBot

import (
	"MetaChat/app/metaChat/qqBot/account"
	"MetaChat/app/metaChat/qqBot/config"
	"MetaChat/app/metaChat/qqBot/io"
	"MetaChat/app/metaChat/qqBot/io/ws"
	"MetaChat/pkg/cq"
	"MetaChat/pkg/signal"
	"time"

	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type QQBot struct {
	log *zap.Logger

	*account.AccountInfo

	config *config.Config

	stop        *signal.StopHandler
	stopCh      chan bool
	connReadyCh <-chan bool

	throwCh chan gjson.Result
	getCh   chan cq.CQResp

	eventHandlers EventHandler
	//use for handling message with event needed
	echoHandlerCh map[string]chan gjson.Result
	//use for send and receive message
	IOHandler io.IOHandler

	isReady bool
}

type EventHandler map[string]map[string]func(msg gjson.Result)

func NewQQBot(log *zap.Logger, config *config.Config, handler io.IOHandler, stopHandler *signal.StopHandler) *QQBot {
	return &QQBot{
		log:           log,
		config:        config,
		stop:          stopHandler,
		stopCh:        make(chan bool),
		IOHandler:     handler,
		throwCh:       make(chan gjson.Result),
		getCh:         make(chan cq.CQResp),
		echoHandlerCh: make(map[string]chan gjson.Result),
		isReady:       false,
	}
}

func (qq *QQBot) OnStart() {
	qq.stop.Add(qq)
	qq.RegisterHandlers()

	qq.connReadyCh = qq.IOHandler.GetOnReadyCh()
	msgCh := qq.IOHandler.GetMessageCh()
	//loop to listen message from IOHandler
	go func() {
		for {
			select {
			case ready := <-qq.connReadyCh:
				if ready {
					qq.isReady = true
					go qq.initAccountInfo()
				} else {
					qq.isReady = false
				}
			case msg := <-msgCh:
				if qq.isReady {
					go qq.onMessage(msg)
					break
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

func (qq *QQBot) SendMessage(msg cq.CQResp) {
	qq.IOHandler.SendMessage(msg)
}

func (qq *QQBot) GetMessageCh() <-chan gjson.Result {
	return qq.throwCh
}

func (qq *QQBot) IsReady() bool {
	return qq.isReady
}

//注册一个channel，用于接受上下文发来的消息
func (qq *QQBot) RegisterEchoHandler(id string) {
	qq.echoHandlerCh[id] = make(chan gjson.Result)
}

//等待channel的消息
func (qq *QQBot) WaitForResult(id string) gjson.Result {
	data := <-qq.echoHandlerCh[id]
	delete(qq.echoHandlerCh, id)
	return data
}

//将需要上层MC处理的消息放入channel，等待上层MetaChat处理
func (qq *QQBot) throw(msg gjson.Result) {
	qq.throwCh <- msg
}

func (qq *QQBot) onMessage(msg gjson.Result) {
	//TODO: handle event
	//handle message that is registered and wait for processing result
	if ch, ok := qq.echoHandlerCh[msg.Get(cq.ECHO).String()]; ok {
		ch <- msg
		return
	}
	//检测注册的event并且处理
	for key, valueMap := range qq.eventHandlers {
		if gjson.Get(msg.String(), key).Exists() {
			for event, handler := range valueMap {
				if gjson.Get(msg.String(), event).Exists() {
					go handler(msg)
				}
			}
		}
	}
}

func (qq *QQBot) notifyStop() {
	//TODO: send stop message to super admin
	//qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetPrivateMessage()))
}

func (qq *QQBot) initAccountInfo() {
	qq.log.Info("intializing account info ...")
	time.Sleep(time.Second * 5)

	login, loginID := cq.GetCQRespEcho(cq.ACTION_GET_LOGIN_INFO, nil)
	qq.SendMessage(login)
	qq.RegisterEchoHandler(loginID)
	loginResult := qq.WaitForResult(loginID)

	group, groupID := cq.GetCQRespEcho(cq.ACTION_GET_GROUP_LIST, nil)
	qq.SendMessage(group)
	qq.RegisterEchoHandler(groupID)
	groupResult := qq.WaitForResult(groupID)

	groupList := make(map[int64]*account.Group)

	groupResult.Get(cq.DATA).ForEach(func(key, value gjson.Result) bool {
		groupList[value.Get(cq.GROUP_ID).Int()] = &account.Group{
			GroupID:   value.Get(cq.GROUP_ID).Int(),
			GroupName: value.Get(cq.GROUP_NAME).String(),
		}
		return true
	})

	friend, friendID := cq.GetCQRespEcho(cq.ACTION_GET_FRIEND_LIST, nil)
	qq.SendMessage(friend)
	qq.RegisterEchoHandler(friendID)
	friendResult := qq.WaitForResult(friendID)

	friendList := make(map[string]*account.User)

	friendResult.Get(cq.DATA).ForEach(func(key, value gjson.Result) bool {
		friendList[value.Get(cq.USER_ID).String()] = &account.User{
			UserID:   value.Get(cq.USER_ID).String(),
			Nickname: value.Get(cq.NICKNAME).String(),
		}
		return true
	})

	qq.AccountInfo = account.NewAccountInfo(
		loginResult.Get(cq.USER_ID).Int(),
		loginResult.Get(cq.NICKNAME).String(),
		friendList,
		groupList,
	)
	qq.log.Info("account info initialized! ", zap.Any("account", qq.AccountInfo))
	qq.log.Info("QQBot is ready!")
}

func Provide() fx.Option {
	return fx.Options(fx.Provide(
		NewQQBot,
		ws.NewWS,
		config.NewConfig,
	))
}
