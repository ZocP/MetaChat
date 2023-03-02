package gpt

const (
	USER      = "user"
	SYSTEM    = "system"
	ASSISTANT = "assistant"
)

type Request struct {
	Model     string     `json:"model"`
	Messages  []Messages `json:"messages"`
	User      string     `json:"user"`
	MaxTokens int64      `json:"max_tokens"`
}
type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type ErrResp struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}
