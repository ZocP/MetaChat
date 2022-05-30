package account

type User struct {
	UserID      string `json:"user_id"`
	Nickname    string `json:"nickname"`
	AccountType string `json:"account_type"`
}

func (user User) GetAccountType() string {
	return user.AccountType
}

func (user User) GetUserID() string {
	return user.UserID
}

func (user User) GetNickName() string {
	return user.Nickname
}
