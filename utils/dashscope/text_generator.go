package dashscope

import (
	"bytes"
	"content_server/setting"
	"encoding/json"
	"fmt"
	"net/http"
)

// TextGenerationRequest 生成文本的请求参数
type TextGenerationRequest struct {
	Model string `json:"model"`
	Input Input  `json:"input"`
}

// TextGenerationResponse 生成文本的响应参数
type TextGenerationResponse struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

type Input struct {
	Prompt   string    `json:"prompt"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Parameters struct {
	TopK int `json:"top_k"`
}

// TextGenerator 生成文本的接口
type TextGenerator interface {
	GenerateText(request TextGenerationRequest) (response TextGenerationResponse, err error)
}

// GenerateText 调用接口生成文本
func (t *TextGenerationRequest) GenerateText() (text string, err error) {
	url := setting.DashScopeSetting.Url
	apiKey := setting.DashScopeSetting.ApiKey
	// 构造body
	bytesData, _ := json.Marshal(TextGenerationRequest{
		Model: "qwen-max",
		Input: Input{
			Prompt: t.Input.Prompt,
		},
		//Parameters: Parameters{
		//	TopK: 50,
		//},
	})
	body := bytes.NewReader(bytesData)
	// 构造请求
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept", "text")
	// 发送请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	// 解析响应
	//var respData map[string]interface{}
	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	// 如果出错
	if respData["code"] != nil {
		err = fmt.Errorf("code: %s, message: %s", respData["code"], respData["message"])
		fmt.Printf("An API error has returned: %s", err)
		return "", err
	}
	// 如果没出错
	text = respData["output"].(map[string]interface{})["text"].(string)
	return text, err
}

func NewTextGenerationRequest(desc string) (req TextGenerationRequest) {
	// desc 如果分行了，则
	prompt := fmt.Sprintf("将以下剧情改写为对话形式，不少于1000字，剧情：%s", desc)
	req = TextGenerationRequest{
		Input: Input{
			Prompt: prompt,
		},
	}
	return req
}

func NewDescGenerationRequest(choose string, content string) (req TextGenerationRequest) {
	prompt := fmt.Sprintf("在以下剧情中我选择了“%s”，续写一段剧情简述，约200字，主角是“我”，配角有有赵卓群（诡异男子），"+
		"周远山（反派富商），林秀芝（村民女孩），陈大志（村长）。剧情：%s", choose, content)
	req = TextGenerationRequest{
		Input: Input{
			Prompt: prompt,
		},
	}
	return req
}

func NewChooseGenerationRequest(desc string) (req TextGenerationRequest) {
	prompt := fmt.Sprintf("为以下剧情生成2个分支选项，每个选项不超过20字：%s", desc)
	req = TextGenerationRequest{
		Input: Input{
			Prompt: prompt,
		},
	}
	return req
}

func NewBgPromptGenerationRequest(desc string) (req TextGenerationRequest) {
	prompt := fmt.Sprintf("20个字描述以下剧情所处的风景，不要出现人物：%s", desc)
	req = TextGenerationRequest{
		Input: Input{
			Prompt: prompt,
		},
	}
	return req
}