package cq

import (
	"regexp"
	"strings"
)

type CQCode struct {
	Type  string
	Param map[string]string
}

func (code CQCode) GetParam(key string) string {
	if val, ok := code.Param[key]; ok {
		return val
	}
	return ""
}

func GetImageCQCode(imageUrl string) string {
	return "[CQ:image,file=" + imageUrl + "]"
}

func ParseCQCode(code string) CQCode {
	var result CQCode
	compiler := regexp.MustCompile("^\\[(CQ:[a-zA-Z]+),(([a-zA-Z]+=[\\w.\\s:\\/-\\/-?=]+),?)+]$")
	if !compiler.MatchString(code) {
		return CQCode{}
	}
	str := regexp.MustCompile("CQ:([a-zA-Z]+)").FindString(code)
	result.Type = str
	params := regexp.MustCompile("([a-zA-Z]+=[\\w.\\s:\\/-\\/-?=]+)").FindAllString(code, -1)
	result.Param = make(map[string]string)
	for _, v := range params {
		parserKey := regexp.MustCompile("^[a-zA-Z]{0,10}=")
		key := parserKey.FindString(v)
		if len(key) >= 2 {
			key = key[:len(key)-1]
		}
		value := v[len(key)+1:]
		result.Param[key] = value
	}
	return result
}

func GetCQCode(tp string, param map[string]string) string {
	var sb strings.Builder
	sb.WriteString("[CQ:")
	sb.WriteString(tp)
	sb.WriteString(",")
	for k, v := range param {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString(",")
	}
	sb.WriteString("]")
	return sb.String()
}

func IsImageCQCode(code string) bool {
	compiler := regexp.MustCompile("^\\[(CQ:[a-zA-Z]+),(([a-zA-Z]+=[\\w.\\s:\\/-\\/-?=]+),?)+]$")
	return compiler.MatchString(code)
}
