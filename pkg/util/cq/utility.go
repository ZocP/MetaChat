package cq

import "regexp"

func IsAtMe(msg string, qq string) bool {
	exp, err := regexp.Compile(`\[CQ:at,qq=` + qq + `\]`)
	if err != nil {
		return false
	}
	return exp.MatchString(msg)

}
