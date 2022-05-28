package ws

import (
	"MetaChat/app/metaChat/qqBot/config"
	"MetaChat/app/metaChat/qqBot/io"
	"MetaChat/pkg/cq"
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
	rawCh     chan []byte
	ready     chan bool
	log       *zap.Logger
}

func (ws *WS) GetMessageCh() <-chan gjson.Result {
	return ws.messageCh
}

func (ws *WS) SendMessage(msg cq.CQResp) {
	if err := ws.Conn.WriteJSON(msg); err != nil {
		ws.log.Error("send message error", zap.Error(err))
		if websocket.IsCloseError(err) {
			ws.log.Info("websocket closed")
			ws.connected = false
		}
	}
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
	}
}

func (ws *WS) OnConnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := Upgrade(c.Writer, c.Request)
		if err != nil {
			ws.log.Error("error while upgrading connection", zap.Error(err))
			return
		}
		ws.Conn = conn
		ws.connected = true
		ws.ready <- true
		go ws.listen()
	}
}

func (ws *WS) OnDisconnect() {
	ws.connected = false
	ws.ready <- false
}

func (ws *WS) listen() {
	defer func(ws *WS) {
		err := ws.Close()
		if err != nil {
			ws.log.Error("error while closing connection", zap.Error(err))
		}
	}(ws)
	for {
		//TODO: 更改成channel
		message, err := ws.ReadMessage()
		if err != nil {
			ws.log.Error("error while reading message", zap.Error(err))
			if websocket.IsCloseError(err) {
				ws.log.Info("connection closed")
				break
			}
		}
		ws.messageCh <- message
	}
}

func (ws *WS) ReadMessage() (gjson.Result, error) {

	_, raw, err := ws.Conn.ReadMessage()
	if err != nil {
		return gjson.Result{}, err
	}

	eventJson := gjson.Parse(string(raw))
	if err != nil {
		return gjson.Result{}, err
	}
	if eventJson.Get(cq.META_EVENT_TYPE).String() != cq.META_EVENT_TYPE_HEARTBEAT {
		//ws.log.Debug("receive message", zap.Any("message", eventJson.String()))
	}
	return eventJson, nil
}
