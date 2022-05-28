package account

type Group struct {
	GroupName string   `json:"group_name"`
	GroupID   int64    `json:"group_id"`
	Admin     int64    `json:"admin"`
	Users     []*int64 `json:"users"`
	BotMode   string   `json:"mode"`
}
