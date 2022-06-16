package saucenao

import (
	"bytes"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	APIKey = "4e393ac5d7382ee2f409909226bb8f64db77d233"
)

func Recognize(url string) gjson.Result {
	zap.S().Debug("识别图片 ", zap.String("url", url))
	apikey := APIKey // api key (get it here https://saucenao.com/user.php?page=search-api)

	client := &http.Client{}
	var data = strings.NewReader(``)
	var url1 = "https://saucenao.com/search.php?output_type=2&testmode=1&numres=5&url=" + url + "&api_key=" + apikey

	req, _ := http.NewRequest("GET", url1, data)
	res, _ := client.Do(req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	gjson.Parse(newStr)
	return gjson.Parse(newStr)
}
