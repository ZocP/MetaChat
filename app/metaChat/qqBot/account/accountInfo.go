package account

type AccountInfo struct {
	AccountId   int64
	Nickname    string
	FriendList  map[string]*ImplementedUser
	GroupList   map[int64]*Group
	AdminList   map[string]*ImplementedUser
	addGroupCh  chan *Group
	delGroupCh  chan int64
	addFriendCh chan *ImplementedUser
	delFriendCh chan string
}

func NewAccountInfo(id int64, name string, flist map[string]*ImplementedUser, glist map[int64]*Group) *AccountInfo {

	account := &AccountInfo{
		AccountId:  id,
		Nickname:   name,
		FriendList: flist,
		GroupList:  glist,
		addGroupCh: make(chan *Group),
		delGroupCh: make(chan int64),
		AdminList: map[string]*ImplementedUser{
			"1395437934": {
				UserID:   "1395437934",
				Nickname: "ZOCP",
			},
		},
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

func (account *AccountInfo) AddFriend(friend *ImplementedUser) {
	account.addFriendCh <- friend
}

func (account *AccountInfo) DelFriend(friendId string) {
	account.delFriendCh <- friendId
}

func (account *AccountInfo) IsAdmin(id string) bool {
	_, ok := account.AdminList[id]
	return ok
}
