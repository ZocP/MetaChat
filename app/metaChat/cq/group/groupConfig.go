package group

const (
	MODE_TRANSFER = "transfer"
	MODE_REPEAT   = "repeat"
)

type Group struct {
	GroupName string `json:"group_name"`
	GroupID   int64  `json:"group_id"`
	BotMode   string `json:"mode"`
}

func (g *Group) GetMode() string {
	return g.BotMode
}

func (g *Group) GetID() int64 {
	return g.GroupID
}
