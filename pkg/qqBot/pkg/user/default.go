package user

type AccountType string

const (
	SUPERADMINUserType AccountType = "superadmin"
	ADMINUserType      AccountType = "admin"
	DefaultUserType    AccountType = "default"
)

type User struct {
	UserID      string      `json:"user_id"`
	NickName    string      `json:"nickname"`
	AccountType AccountType `json:"account_type"`
}

func (d *User) GetUserID() string {
	return d.UserID
}

func (d *User) GetNickName() string {
	if d.NickName == "" {
		return d.UserID
	}
	return d.NickName
}

func (d *User) GetAccountType() AccountType {
	return d.AccountType
}

func NewUser(userID string, accountType AccountType, nickname string) *User {
	return &User{
		UserID:      userID,
		NickName:    nickname,
		AccountType: accountType,
	}
}

func GetAccountType(t string) AccountType {
	switch t {
	case "default":
		return DefaultUserType
	case "admin":
		return ADMINUserType
	case "superadmin":
		return SUPERADMINUserType
	default:
		return DefaultUserType
	}
}

func DefaultUser() User {
	return User{
		UserID:      "default",
		NickName:    "default",
		AccountType: DefaultUserType,
	}
}
