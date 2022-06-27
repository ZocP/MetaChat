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

const (
	RESULTS    = "results"
	DATA       = "data"
	URL        = "ext_urls"
	HEADER     = "header"
	SIMILARITY = "similarity"
	TITLE      = "title"
	THUMBNAIL  = "thumbnail"
)

type QQSauceNaoResult struct {
	MatchRate string
	URL       []string
	Title     string
	Thumbnail string
}

func GetQQSauceNaoResult(msg gjson.Result) string {
	resultRaw := ""
	result := make([]QQSauceNaoResult, 0, 10)
	msg.ForEach(func(key, value gjson.Result) bool {
		header := value.Get(HEADER)
		data := value.Get(DATA)
		urls := make([]string, 0, 10)
		data.Get(URL).ForEach(func(key, value gjson.Result) bool {
			urls = append(urls, value.String())
			return true
		})
		result = append(result, QQSauceNaoResult{
			MatchRate: header.Get(SIMILARITY).String(),
			URL:       urls,
			Title:     data.Get(TITLE).String(),
			Thumbnail: header.Get(THUMBNAIL).String(),
		})
		return true
	})
	for _, v := range result {
		resultRaw += "匹配率: " + v.MatchRate + "\n"
		resultRaw += "标题: " + v.Title + "\n"
		resultRaw += "链接: " + strings.Join(v.URL, "\n") + "\n"
	}
	return resultRaw
}

func Recognize(url string) gjson.Result {

	zap.S().Info("识别图片 ", zap.String("url", url))
	apikey := APIKey // api key (get it here https://saucenao.com/user.php?page=search-api)

	client := &http.Client{}
	var data = strings.NewReader(``)
	var url1 = "https://saucenao.com/search.php?output_type=2&testmode=1&numres=5&url=" + url + "&api_key=" + apikey

	req, _ := http.NewRequest("GET", url1, data)
	res, _ := client.Do(req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()
	result := gjson.Parse(newStr).Get(RESULTS)

	return result
}
