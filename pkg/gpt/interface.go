package gpt

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

var GPTInstance *GPT

type GPT struct {
	Model    string
	APIKey   string
	MaxToken int64
}

func SetConfig(viper *viper.Viper) {
	zap.L().Info("set gpt config")
	gpt := &GPT{}
	err := viper.UnmarshalKey("gpt", gpt)
	if err != nil {
		return
	}
	zap.L().Info("gpt key is set, it is : " + gpt.APIKey)
	GPTInstance = gpt
}

func SendReq(content string) (*Response, error) {
	url := "https://api.openai.com/v1/chat/completions"
	method := "POST"

	payload, err := json.Marshal(Request{
		Model: GPTInstance.Model,
		Messages: []Messages{{
			USER,
			content,
		},
		},
		MaxTokens: GPTInstance.MaxToken,
	})

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}
	if GPTInstance == nil {
		return nil, nil
	}
	req.Header.Add("Authorization", GPTInstance.APIKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := &Response{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil

}
