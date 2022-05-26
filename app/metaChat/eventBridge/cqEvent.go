package eventBridge

const (
	META_EVENT_TYPE_HEARTBEAT = "heartbeat"
	META_EVENT_TYPE_LIFECYCLE = "lifecycle"
)

type CQEvent struct {
	Interval      int    `json:"interval"`
	MetaEventType string `json:"meta_event_type"`
	PostType      string `json:"post_type"`
	SelfId        int64  `json:"self_id"`
	Status        struct {
		AppEnabled     bool        `json:"app_enabled"`
		AppGood        bool        `json:"app_good"`
		AppInitialized bool        `json:"app_initialized"`
		Good           bool        `json:"good"`
		Online         bool        `json:"online"`
		PluginsGood    interface{} `json:"plugins_good"`
		Stat           struct {
			PacketReceived  int `json:"PacketReceived"`
			PacketSent      int `json:"PacketSent"`
			PacketLost      int `json:"PacketLost"`
			MessageReceived int `json:"MessageReceived"`
			MessageSent     int `json:"MessageSent"`
			LastMessageTime int `json:"LastMessageTime"`
			DisconnectTimes int `json:"DisconnectTimes"`
			LostTimes       int `json:"LostTimes"`
		} `json:"stat"`
	} `json:"status"`
	Time int `json:"time"`
}
