package account

type ImplementedUser struct {
	UserID      string `json:"user_id"`
	Nickname    string `json:"nickname"`
	AccountType string `json:"account_type"`
}

func (user ImplementedUser) GetAccountType() string {
	return user.AccountType
}

func (user ImplementedUser) GetUserID() string {
	return user.UserID
}

func (user ImplementedUser) GetNickName() string {
	return user.Nickname
}
