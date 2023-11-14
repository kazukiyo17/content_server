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
	//Parameters Parameters `json:"parameters"`
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
func (t *TextGenerationRequest) GenerateText() (response TextGenerationResponse, err error) {
	url := setting.DashScopeSetting.Url
	apiKey := setting.DashScopeSetting.ApiKey
	// 构造body
	bytesData, _ := json.Marshal(TextGenerationRequest{
		Model: "qwen-turbo",
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
		return
	}
	response.Text = respData["output"].(map[string]interface{})["text"].(string)
	response.FinishReason = respData["output"].(map[string]interface{})["finish_reason"].(string)
	if err != nil {
		return
	}
	return
}

func NewTextGenerationRequest(choose string, content string) (req TextGenerationRequest) {
	prompt := fmt.Sprintf("我在以下剧情中选择了[%s]，续写角色对话，不要剧情描述，只要角色对话，changeBg为场景描述(如村间小路，"+
		"不要出现角色名)，最后给出三个分支choose，用|分隔，不少于500字。角色有赵卓群（诡异男子），周远山（反派富商），林秀芝（村民女孩），"+
		"陈大志（村长）。剧本如下：%s", choose, content)
	req = TextGenerationRequest{
		Input: Input{
			Prompt: prompt,
		},
	}
	return req
}

// Test 测试
func Test() {
	request := TextGenerationRequest{
		Input: Input{
			Prompt: "我在以下剧情中选择了[跟着林秀芝]，续写角色对话，不要剧情描述，只要角色对话，changeBg为风景描述(如村间小路，不要出现角色名)，最后给出三个分支choose，用|分隔，不少于500字。角色有赵卓群（诡异男子），周远山（反派富商），林秀芝（村民女孩），陈大志（村长）。剧本如下：changeBg: 一个静谧的小山村，空气清新，景色宜人。\\n 清晨的阳光照射进小镇的街道，我正在探寻这个地方。\\n此时一个小女孩蹦蹦跳跳地走在前方\\n我:小朋友，请问这里是不是有一座寺庙？\\n林秀芝:我叫林秀芝，大哥哥你想去哪里？\\n我:我想去拍那座寺庙的照片。\\n林秀芝:那座寺庙？\\n林秀芝不语，盯着我看了足足一分钟\\n随后女孩诡异地笑了\\n林秀芝:秀芝不知道哦\\nchoose:独自向前探索|跟着林秀芝|到村庄里寻找其他人",
		},
	}
	response, err := request.GenerateText()
	if err != nil {
		return
	}
	println(response.Text)
}
