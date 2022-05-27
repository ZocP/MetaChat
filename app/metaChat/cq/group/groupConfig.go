package group

const (
	MODE_TRANSFER = "transfer"
	MODE_REPEAT   = "repeat"
)

type Group struct {
	GroupName string `json:"group_name"`
	GroupID   int64  `json:"group_id"`
	BotMode   string `json:"mode"`

	UserHandler map[int64]func()
}

func (g *Group) GetMode() string {
	return g.BotMode
}

func (g *Group) SetUserHandler(id int64, f func()) {
	g.UserHandler[id] = f
}

func (g *Group) GetUserHandler(id int64) func() {
	return g.UserHandler[id]
}

func (g *Group) ClearUserHandler(id int64) {
	delete(g.UserHandler, id)
}

func (g *Group) GetID() int64 {
	return g.GroupID
}
