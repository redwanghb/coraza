package llmguard

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var (
	COMPLETION = "/api/generate"
	CHAT       = "/api/chat"
)

// 参考https://github.com/ollama/ollama/blob/main/docs/api.md

func ollamaParseReq(uri, body string) (question string, images []byte, apiType string, stream bool, err error) {
	switch uri {
	case COMPLETION:
		apiType = "completion"
		// TODO 后续扩展images，提取多模态图片内容列表
		images = nil
		question_result := gjson.Get(body, "prompt")
		if question_result.Exists() {
			question = question_result.String()
		} else {
			question = ""
		}
		stream_result := gjson.Get(body, "stream")
		if stream_result.Exists() {
			stream = stream_result.Bool()
		} else {
			stream = false
		}
		err = nil
		return
	case CHAT:
		apiType = "chat"

		// 提取stream内容，判断答案是以SSE返回还是一次性返回
		stream_result := gjson.Get(body, "stream")
		if stream_result.Exists() {
			stream = stream_result.Bool()
		} else {
			stream = false
		}
		// TODO 后续扩展images，提取多模态图片内容列表
		images = nil

		// 提取content内容
		result := gjson.Get(body, "messages")
		if !result.Exists() {
			question = ""
			err = fmt.Errorf("there is no messages for uri %s", uri)
			return
		}
		if !result.IsArray() {
			question = ""
			err = fmt.Errorf("messages of uri %s is not a array type", uri)
			return
		}
		messages := result.Array()
		if len(messages) == 0 {
			question = ""
			err = fmt.Errorf("messages of uri %s is a null array", uri)
			return
		}

		content := messages[len(messages)-1].Get("content")
		if content.Exists() {
			question = content.String()
			err = nil
		} else {
			question = ""
			err = fmt.Errorf("there is no content of uri %s", uri)
		}
		return
	default:
		err = fmt.Errorf("not support uri %s", uri)
		return
	}
}

func ollamaParseRes(uri, body string) (answer string, images []byte, apiType string, stream bool, done bool, err error) {
	switch uri {
	case COMPLETION:
		apiType = "completion"
		stream_result := gjson.Get(body, "stream")
		if stream_result.Exists() {
			stream = stream_result.Bool()
		} else {
			stream = false
		}
		answer_result := gjson.Get(body, "response")
		if answer_result.Exists() {
			answer = answer_result.String()
			err = nil
		} else {
			answer = ""
			err = fmt.Errorf("there is no response of uri %s", uri)
		}
		isDone := gjson.Get(body, "done")
		if isDone.Exists() {
			done = isDone.Bool()
		} else {
			done = false // 如果没有done字段，默认认为是未结束的回答
		}
		return
	case CHAT:
		apiType = "chat"
		stream_result := gjson.Get(body, "stream")
		if stream_result.Exists() {
			stream = stream_result.Bool()
		} else {
			stream = false
		}
		answer_result := gjson.Get(body, "message.content")
		if answer_result.Exists() {
			answer = answer_result.String()
			err = nil
		} else {
			answer = ""
			err = fmt.Errorf("there is no message.content of uri %s", uri)
		}
		isDone := gjson.Get(body, "done")
		if isDone.Exists() {
			done = isDone.Bool()
		} else {
			done = false // 如果没有done字段，默认认为是未结束的回答
		}
		// TODO 后续扩展images，提取多模态图片内容列表
		images = nil
		return
	default:
		err = fmt.Errorf("not support uri %s", uri)
		return
	}
}

func GetOllamaChatAnsMessage() string {
	return `{"model": "", "create_at": "%s", "message": {"role": "assistant", "content": "%s"}, "done": true}
	`
}

func GetOllamaCompletionAnsMessage() string {
	return `{"model": "", "create_at": "%s", "response": "%s", "done": true}
	`
}
