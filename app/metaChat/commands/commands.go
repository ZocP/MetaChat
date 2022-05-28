package commands

var (
	DefaultMCNormalUserCommands = []string{
		"/ls",
	}
)

type Command struct {
	Root  string
	Param []string
}
