package lolicon

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	NOTR18 = 0
	R18    = 1
	MIX    = 2
)

type Param struct {
	Adult      int      `json:"r18,omitempty"`
	Num        int      `json:"num,omitempty"`
	Uid        int      `json:"uid,omitempty"`
	KeyWord    string   `json:"keyword,omitempty"`
	Tag        []string `json:"tag,omitempty"`
	Size       []string `json:"size,omitempty"`
	Proxy      string   `json:"proxy,omitempty"`
	DateAfter  int      `json:"dateAfter,omitempty"`
	DateBefore int      `json:"dateBefore,omitempty"`
	DSC        bool     `json:"dsc,omitempty"`
}

func GetRandomPictureJSON(param Param) (gjson.Result, error) {
	url := "https://api.lolicon.app/setu/v2"
	method := "POST"
	raw, err := json.Marshal(param)
	if err != nil {
		return gjson.Result{}, err
	}
	payload := strings.NewReader(string(raw))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(body), nil
}

func ParseParam(m map[string][]string) (Param, error) {
	var (
		result Param
		err    error
		num    int
	)
	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		switch k {
		case "r18":
			//result.Adult, err = strconv.Atoi(v[0])
			//if err != nil{
			//	return result, err
			//}
		case "num":
			num, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
			if num > 3 {
				num = 3
			}
			result.Num = num
		case "uid":
			result.Uid, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
		case "keyword":
			result.KeyWord = v[0]
		case "tag":
			result.Tag = v
		case "size":
			result.Size = v
		case "proxy":
			result.Proxy = v[0]
		case "dateAfter":
			result.DateAfter, err = strconv.Atoi(v[0])

			if err != nil {
				return result, err
			}
		case "dateBefore":
			result.DateBefore, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
		case "dsc":
			result.DSC, err = strconv.ParseBool(v[0])
			if err != nil {
				return result, err
			}
		}
	}
	result.Adult = NOTR18
	return result, nil
}

func ParseParamAll(m map[string][]string) (Param, error) {
	var (
		result Param
		err    error
		num    int
	)
	for k, v := range m {
		switch k {
		case "r18":
			result.Adult, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
		case "num":
			num, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
			if num > 50 {
				num = 50
			}
			result.Num = num
		case "uid":
			result.Uid, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
		case "keyword":
			result.KeyWord = v[0]
		case "tag":
			result.Tag = v
		case "size":
			result.Size = v
		case "proxy":
			result.Proxy = v[0]
		case "dateAfter":
			result.DateAfter, err = strconv.Atoi(v[0])

			if err != nil {
				return result, err
			}
		case "dateBefore":
			result.DateBefore, err = strconv.Atoi(v[0])
			if err != nil {
				return result, err
			}
		case "dsc":
			result.DSC, err = strconv.ParseBool(v[0])
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}
