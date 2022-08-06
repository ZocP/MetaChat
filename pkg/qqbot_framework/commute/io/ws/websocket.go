package ws

import (
	"MetaChat/pkg/qqbot_framework/commute/config"
	"MetaChat/pkg/qqbot_framework/commute/io"
	"MetaChat/pkg/util/cq"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type WS struct {
	*websocket.Conn
	connected bool
	config    *config.Config
	messageCh chan gjson.Result
	writeCh   chan cq.CQResponse
	rawCh     chan []byte
	ready     chan bool
	stopCh    chan bool
	log       *zap.Logger
}

func (ws *WS) GetMessageCh() <-chan gjson.Result {
	return ws.messageCh
}

func (ws *WS) SendMessage(msg cq.CQResponse) {
	ws.writeCh <- msg
}

func (ws *WS) GetOnReadyCh() <-chan bool {
	return ws.ready
}

func NewWS(config *config.Config, log *zap.Logger) io.IOHandler {
	return &WS{
		config:    config,
		log:       log,
		messageCh: make(chan gjson.Result),
		rawCh:     make(chan []byte),
		ready:     make(chan bool),
		writeCh:   make(chan cq.CQResponse),
		stopCh:    make(chan bool),
	}
}

func (ws *WS) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := Upgrade(c.Writer, c.Request)
		if err != nil {
			ws.log.Error("error while upgrading connection", zap.Error(err))
			return
		}
		ws.Conn = conn
		ws.connected = true
		ws.ready <- true
		ws.log.Info("websocket connected")
		go ws.listen()
		go ws.ReadMessage()
	}
}

func (ws *WS) OnDisconnect() {
	ws.connected = false
	ws.ready <- false
}

func (ws *WS) listen() {
	defer func() {
		err := ws.Close()
		if err != nil {
			ws.log.Error("error while closing connection", zap.Error(err))
		}
	}()
	for {
		select {
		case <-ws.stopCh:
			ws.OnDisconnect()
			ws.log.Info("websocket closed")
			return
		case raw := <-ws.rawCh:
			data := gjson.Parse(string(raw))
			//过滤心跳包，不作pong处理
			if data.Get(cq.META_EVENT_TYPE).String() != cq.META_EVENT_TYPE_HEARTBEAT {
				ws.messageCh <- data
			}
		case write := <-ws.writeCh:
			if err := ws.Conn.WriteJSON(write); err != nil {
				ws.log.Error("send message error", zap.Error(err))
				if websocket.IsCloseError(err) {
					ws.log.Info("websocket closed")
					ws.connected = false
				}
			}
		}
	}
}

func (ws *WS) ReadMessage() {
	go func() {
		for {
			_, raw, err := ws.Conn.ReadMessage()
			if err != nil {
				ws.log.Error("error while reading message", zap.Error(err))
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					ws.log.Info("connection closed")
					ws.connected = false
					ws.ready <- false
					break
				}
				ws.log.Info("connection abnormally closed", zap.Error(err))
				ws.connected = false
				ws.ready <- false
				break
			}
			ws.rawCh <- raw
		}
	}()
}
