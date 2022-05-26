package cq

import (
	"MetaChat/app/metaChat/cq/config"
	"MetaChat/app/metaChat/cq/group"
	"MetaChat/app/metaChat/cq/user"
	"MetaChat/app/metaChat/cq/ws"
	"MetaChat/app/metaChat/eventBridge/request"
	"MetaChat/app/metaChat/eventBridge/response"
	"MetaChat/pkg/signal"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CQEventHandler struct {
	*QQBot
	config       *config.Config
	conn         *ws.WS
	eventChannel chan gjson.Result
	replyChannel chan response.CQResp
	log          *zap.Logger
	stopHandler  *signal.StopHandler
	stopCh       chan chan bool
}

func NewCQEventHandler(viper *viper.Viper, logger *zap.Logger, handler *signal.StopHandler) *CQEventHandler {
	return &CQEventHandler{
		log:          logger,
		config:       config.Unmarshal(viper),
		eventChannel: make(chan gjson.Result),
		replyChannel: make(chan response.CQResp),
		stopHandler:  handler,
	}
}

func (cq *CQEventHandler) OnStart() {
	cq.stopHandler.Add(cq)
	for {
		select {
		case reply := <-cq.replyChannel:
			if err := cq.conn.WriteJSON(reply); err != nil {
				cq.log.Error("error while writing message", zap.Error(err))
			}

			//case stop := <-cq.stopCh:
			//	if err := cq.conn.WriteJSON(); err != nil{
			//		cq.log.Error("error while writing message", zap.Error(err))
			//	}
			//	stop <- true

		}
	}
}

func (cq *CQEventHandler) OnStop() error {
	cq.stopCh <- make(chan bool)
	<-cq.stopCh
	err := cq.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (cq *CQEventHandler) OnConnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := ws.Upgrade(c.Writer, c.Request)
		if err != nil {
			cq.log.Error("error while upgrading connection", zap.Error(err))
			return
		}
		cq.conn = ws.NewWS(conn, cq.config, cq.log)
		cq.log.Info("connection with CQHttp established")
		cq.QQBot = cq.getNewQQBot()
		go cq.listen()
	}
}

func (cq *CQEventHandler) getNewQQBot() *QQBot {
	if err := cq.conn.WriteJSON(response.GetCQResp("get_login_info", nil)); err != nil {
		cq.log.Error("error while writing message", zap.Error(err))
	}
	accountInfo, err := cq.conn.ReadMessage()
	if err != nil {
		cq.log.Error("error while reading message", zap.Error(err))
	}

	if err := cq.conn.WriteJSON(response.GetCQResp("get_group_list", nil)); err != nil {
		cq.log.Error("error while writing message", zap.Error(err))
	}

	groupInfo, err := cq.conn.ReadMessage()
	groupList := make(map[int64]*group.Group)
	groupInfo.Get(request.DATA).ForEach(func(key, value gjson.Result) bool {
		groupList[value.Get("group_id").Int()] = &group.Group{

			GroupID:   value.Get("group_id").Int(),
			GroupName: value.Get("group_name").String(),
			BotMode:   group.MODE_REPEAT,
		}
		return true
	})

	if err := cq.conn.WriteJSON(response.GetCQResp("get_friend_list", nil)); err != nil {
		cq.log.Error("error while writing message", zap.Error(err))
	}
	friendInfo, err := cq.conn.ReadMessage()
	friendList := make(map[int64]*user.User)
	friendInfo.Get(request.DATA).ForEach(func(key, value gjson.Result) bool {
		character := user.NORMAL
		//if user is admin
		if cq.isAdmin(value.Get("user_id").Int()) {
			character = user.ADMIN
		}
		friendList[value.Get("user_id").Int()] = &user.User{
			Nickname:  value.Get("nickname").String(),
			Character: character,
		}
		return true
	})

	cq.log.Debug("initialized qq info")
	return &QQBot{
		AccountId:  accountInfo.Get("user_id").Int(),
		Nickname:   accountInfo.Get("nickname").String(),
		FriendList: friendList,
		GroupList:  groupList,
	}
}

func (cq *CQEventHandler) isAdmin(id int64) bool {
	admin := []int64{
		1395437934,
	}
	for _, v := range admin {
		if v == id {
			return true
		}
	}
	return false
}

func (cq *CQEventHandler) GetBot() *QQBot {
	return cq.QQBot
}

func (cq *CQEventHandler) listen() {
	for {
		eventJson, err := cq.conn.ReadMessage()
		if err != nil {
			cq.log.Error("error while reading message", zap.Error(err))
		}
		cq.eventChannel <- eventJson
	}
}

func Provide() fx.Option {
	return fx.Provide(NewCQEventHandler)
}

func (cq *CQEventHandler) GetEventCh() chan gjson.Result {
	return cq.eventChannel
}

func (cq *CQEventHandler) GetReplyCh() chan response.CQResp {
	return cq.replyChannel
}
