package dashscope

import (
	"bytes"
	"content_server/setting"
	"encoding/json"
	"net/http"
)

// TextGenerationRequest 生成文本的请求参数
type TextGenerationRequest struct {
	Prompt string `json:"prompt"`
}

// TextGenerationResponse 生成文本的响应参数
type TextGenerationResponse struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

// TextGenerator 生成文本的接口
type TextGenerator interface {
	GenerateText(request TextGenerationRequest) (response TextGenerationResponse, err error)
}

// GenerateText 调用接口生成文本
func (t *TextGenerationRequest) GenerateText() (response TextGenerationResponse, err error) {
	url := setting.DashScopeSetting.Url
	apiKey := setting.DashScopeSetting.ApiKey
	// 构造body
	bytesData, _ := json.Marshal(TextGenerationRequest{
		Prompt: t.Prompt,
	})
	body := bytes.NewReader(bytesData)
	// 构造请求
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text")
	// 发送请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// 解析响应
	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	response.Text = respData["output"].(map[string]interface{})["text"].(string)
	response.FinishReason = respData["output"].(map[string]interface{})["finish_reason"].(string)
	if err != nil {
		return
	}
	return
}
