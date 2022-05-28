package account

type AccountInfo struct {
	AccountId   int64
	Nickname    string
	FriendList  map[int64]*User
	GroupList   map[int64]*Group
	addGroupCh  chan *Group
	delGroupCh  chan int64
	addFriendCh chan *User
	delFriendCh chan int64
}

func NewAccountInfo(id int64, name string, flist map[int64]*User, glist map[int64]*Group) *AccountInfo {

	account := &AccountInfo{
		AccountId:  id,
		Nickname:   name,
		FriendList: flist,
		GroupList:  glist,
		addGroupCh: make(chan *Group),
		delGroupCh: make(chan int64),
	}
	return account
}

func (account *AccountInfo) OnStart() {
	go func() {
		for {
			select {
			case newGroup := <-account.addGroupCh:
				account.GroupList[newGroup.GroupID] = newGroup
			case groupId := <-account.delGroupCh:
				delete(account.GroupList, groupId)
			case newFriend := <-account.addFriendCh:
				account.FriendList[newFriend.UserID] = newFriend
			case friendId := <-account.delFriendCh:
				delete(account.FriendList, friendId)
			}
		}
	}()
}

func (account *AccountInfo) GetGroupInfo(groupId int64) Group {
	return *account.GroupList[groupId]
}

func (account *AccountInfo) AddGroup(group *Group) {
	account.addGroupCh <- group
}

func (account *AccountInfo) DelGroup(groupId int64) {
	account.delGroupCh <- groupId
}

func (account *AccountInfo) AddFriend(friend *User) {
	account.addFriendCh <- friend
}

func (account *AccountInfo) DelFriend(friendId int64) {
	account.delFriendCh <- friendId
}
