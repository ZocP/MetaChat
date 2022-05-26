package cq

import (
	"MetaChat/app/metaChat/cq/group"
	"MetaChat/app/metaChat/cq/user"
)

type RawInfo struct {
}

type QQBot struct {
	AccountId  int64
	Nickname   string
	FriendList map[int64]*user.User
	GroupList  map[int64]*group.Group
}

func (qq *QQBot) GetAccountId() int64 {
	return qq.AccountId
}

func (qq *QQBot) GetGroup(id int64) *group.Group {
	return qq.GroupList[id]
}
