package mc

import "time"

const (
	MESSAGE_TYPE_CHAT    = "chat"
	MESSAGE_TYPE_COMMAND = "command"
)

type MCResponse struct {
	MessageType string      `json:"msg_type"`
	TimeStamp   int64       `json:"timestamp"`
	Data        interface{} `json:"data"`
}

type CommandData struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func NewMCResponse(mt string, data interface{}) MCResponse {
	return MCResponse{
		MessageType: mt,
		TimeStamp:   time.Now().Unix(),
		Data:        data,
	}
}
