package ws

import (
	"MetaChat/app/metaChat/cq/config"
	"MetaChat/pkg/cq"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type WS struct {
	*websocket.Conn
	config *config.Config
	log    *zap.Logger
}

func NewWS(conn *websocket.Conn, config *config.Config, log *zap.Logger) *WS {
	return &WS{
		Conn:   conn,
		config: config,
		log:    log,
	}
}

func (ws *WS) WriteJson(data interface{}) error {
	ws.log.Debug("writing json", zap.Any("data", data))
	return ws.WriteJSON(data)
}

func (ws *WS) ReadMessage() (gjson.Result, error) {
	_, raw, err := ws.Conn.ReadMessage()
	if err != nil {
		return gjson.Result{}, err
	}
	eventJson := gjson.Parse(string(raw))
	if err != nil {
		ws.log.Error("error while unmarshalling message", zap.Error(err))
		return gjson.Result{}, err
	}
	if eventJson.Get(cq.META_EVENT_TYPE).String() != cq.META_EVENT_TYPE_HEARTBEAT {
		ws.log.Debug("receive message", zap.Any("message", eventJson.String()))
	}
	return eventJson, nil
}
