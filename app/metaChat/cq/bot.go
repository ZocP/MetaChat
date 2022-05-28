package cq

import (
	"MetaChat/app/metaChat/cq/group"
	"MetaChat/app/metaChat/cq/user"
)

type QQBot struct {
	AccountId  int64
	Nickname   string
	FriendList map[int64]*user.User
	GroupList  map[int64]*group.Group
	AdminList  map[int64]*user.User
	addGroupCh chan *group.Group
	delGroupCh chan int64
}

func (qq *QQBot) OnStart() {
	go func() {
		for {
			select {
			case g := <-qq.addGroupCh:
				qq.GroupList[g.GetID()] = g

			case id := <-qq.delGroupCh:
				delete(qq.GroupList, id)

			}
		}
	}()
}

func (qq *QQBot) GetAccountId() int64 {
	return qq.AccountId
}

func (qq *QQBot) GetGroup(id int64) *group.Group {
	return qq.GroupList[id]
}

func (qq *QQBot) GetAdminList() map[int64]*user.User {
	return qq.AdminList
}

func (qq *QQBot) AddGroup(g *group.Group) {
	qq.addGroupCh <- g
}
