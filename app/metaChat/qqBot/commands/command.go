package commands

import (
	"regexp"
	"strings"
)

//map[flag][]values

type Command struct {
	Name     string
	Param    map[string][]string
	Raw      string
	HasParam bool
}

func UnpackCommand(cmd string) Command {
	str := cmd[2:]
	params := strings.Split(str, " ")
	if len(params) == 1 {
		return Command{Name: params[0], Raw: cmd, HasParam: false}
	}

	param := make(map[string][]string)
	var flag string
	for _, v := range params[1:] {
		compiler := regexp.MustCompile("^-")
		if compiler.MatchString(v) {
			flag = v[1:]
			param[flag] = []string{}
			continue
		}
		param[flag] = append(param[flag], v)
	}

	return Command{
		Name:     params[0],
		Param:    param,
		Raw:      cmd,
		HasParam: true,
	}
}
