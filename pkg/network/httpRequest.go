package network

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func GetFromUrlJSON(url string, query map[string]string) (gjson.Result, error) {
	var rawUrl string
	rawUrl += url + "?"
	for k, v := range query {
		rawUrl += fmt.Sprintf("&%s=%s", k, v)
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, rawUrl, nil)

	if err != nil {
		fmt.Println(err)
		return gjson.Result{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
