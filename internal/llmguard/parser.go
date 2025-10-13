package llmguard

import (
	"fmt"
	"time"
)

type LLMRequestInfo struct {
	Provider string // 当前仅支持ollama，后续扩展其他provider
	APIType  string // chat, completion
	Question string
	Images   []byte
	Stream   bool // 如果是true，说明是sse模式

}

func (l *LLMRequestInfo) GetQuestion() string {
	return l.Question
}

func (l *LLMRequestInfo) GetAPIType() string {
	return l.APIType
}

func (l *LLMRequestInfo) GetStream() bool {
	return l.Stream
}

func (l *LLMRequestInfo) GetProvider() string {
	return l.Provider
}

func ParseLLMRequest(provider string, uri, body string) (*LLMRequestInfo, error, bool) {
	switch provider {
	case "ollama":
		question, images, apiType, stream, err := ollamaParseReq(uri, body)
		if err != nil {
			return nil, err, false
		}
		return &LLMRequestInfo{
			Provider: "ollama",
			APIType:  apiType,
			Question: question,
			Images:   images,
			Stream:   stream,
		}, nil, true
	default:
		return nil, fmt.Errorf("provider %s is not support", provider), false
	}
}

func GetLLMResMessage(questionInfo *LLMRequestInfo) string {
	// 设置告警消息
	switch questionInfo.GetProvider() {
	case "ollama":
		if questionInfo.GetAPIType() == "chat" && questionInfo.Stream {
			res_template := GetOllamaChatAnsMessage()
			timeString := time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")
			return fmt.Sprintf(res_template, timeString, "您输入的内容包含了提示词注入攻击")

		}
		if questionInfo.GetAPIType() == "completion" && questionInfo.Stream {
			res_template := GetOllamaCompletionAnsMessage()
			timeString := time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")
			return fmt.Sprintf(res_template, timeString, "您输入的内容包含了提示词注入攻击")
		}
	default:
		return "您输入的内容包含了提示词注入攻击"
	}
	return ""
}

func GetLLMRespMessage(responseInfo *LLMResponseInfo) string {
	// 设置告警消息
	switch responseInfo.GetProvider() {
	case "ollama":
		if responseInfo.GetAPIType() == "chat" {
			res_template := GetOllamaChatAnsMessage()
			timeString := time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")
			return fmt.Sprintf(res_template, timeString, "您输入的内容包含了提示词注入攻击")
		}
		if responseInfo.GetAPIType() == "completion" {
			res_template := GetOllamaCompletionAnsMessage()
			timeString := time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")
			return fmt.Sprintf(res_template, timeString, "您输入的内容包含了提示词注入攻击")
		}
	default:
		return "您输入的内容包含了提示词注入攻击"
	}
	return ""
}

type LLMResponseInfo struct {
	Provider string // 当前仅支持ollama，后续扩展其他provider
	APIType  string // chat, completion
	Answer   string
	Images   []byte
	Stream   bool // 如果是true，说明是sse模式
	Done     bool // 是否完成
}

func (lr *LLMResponseInfo) GetAnswer() string {
	return lr.Answer
}

func (lr *LLMResponseInfo) GetAPIType() string {
	return lr.APIType
}

func (lr *LLMResponseInfo) GetStream() bool {
	return lr.Stream
}

func (lr *LLMResponseInfo) GetProvider() string {
	return lr.Provider
}

func (lr *LLMResponseInfo) IsDone() bool {
	return lr.Done
}

func ParseLLMResponse(provider string, uri, body string) (*LLMResponseInfo, error, bool) {
	switch provider {
	case "ollama":
		answer, images, apiType, stream, done, err := ollamaParseRes(uri, body)
		if err != nil {
			return nil, err, false
		}
		return &LLMResponseInfo{
			Provider: "ollama",
			APIType:  apiType,
			Answer:   answer,
			Images:   images,
			Stream:   stream,
			Done:     done,
		}, nil, true
	default:
		return nil, fmt.Errorf("provider %s is not support", provider), false
	}
}
