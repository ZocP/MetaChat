package commands

import "strings"

type Command struct {
	Name     string
	Param    []string
	Raw      string
	HasParam bool
}

func UnpackCommand(cmd string) Command {
	str := cmd[2:]
	params := strings.Split(str, " ")
	if len(params) == 1 {
		return Command{Name: params[0], Raw: cmd, HasParam: false}
	}
	return Command{
		Name:     params[0],
		Param:    params[1:],
		Raw:      cmd,
		HasParam: true,
	}
}
