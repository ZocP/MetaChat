package status

type AllowStatus int

const (
	Command_Not_Allowed     AllowStatus = 0
	Command_OK              AllowStatus = 1
	Command_Not_Allow_Param AllowStatus = 2
)
