package user

const (
	ADMIN  = "admin"
	NORMAL = "normal"
)

type User struct {
	UserId    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	Character string `json:"character"`
}
