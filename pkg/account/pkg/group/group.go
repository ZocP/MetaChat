package group

type GroupType string

type Group struct {
	GroupName  string   `json:"group_name"`
	GroupID    string   `json:"group_id"`
	GroupUsers []string `json:"group_users"`
}

//尽量避免直接使用fields
func (g *Group) GetName() string {
	return g.GroupName
}

func (g *Group) GetID() string {
	return g.GroupID
}

func (g *Group) GetUsers() []string {
	return g.GroupUsers
}
