package cq

type QQBot struct {
	AccountId  int64
	Nickname   string
	FriendList map[int64]string
	GroupList  map[int64]string
	AdminList  map[int64]string
}

func (bot *QQBot) IsAdmin(id int64) bool {
	_, ok := bot.AdminList[id]
	return ok
}
