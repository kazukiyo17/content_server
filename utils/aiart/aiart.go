package aiart

import (
	"content_server/setting"
	"fmt"
	aiart "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/aiart/v20221229"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

func Generate(text string) (resp string, err error) {
	credential := common.NewCredential(
		setting.TencentCloudSetting.SecretID,
		setting.TencentCloudSetting.SecretKey,
	)
	// 实例化一个客户端配置对象，可以指定超时时间等配置
	cpf := profile.NewClientProfile()
	// 实例化要请求产品的client对象
	client, _ := aiart.NewClient(credential, regions.Shanghai, cpf)
	// 实例化一个请求对象
	request := aiart.NewTextToImageRequest()
	request.Prompt = &text
	request.Styles = []*string{common.StringPtr("301"), common.StringPtr("201")}
	request.ResultConfig = &aiart.ResultConfig{
		Resolution: common.StringPtr("768:1024"),
	}
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.TextToImage(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		panic(err)
	}
	return response.ToJsonString(), err
}

// 测试
func Test() {
	resp, err := generate("高大的中年男子，修剪整齐的短发，眼神深沉而狡猾，黑色西装，银框眼镜，优雅而又神秘。仅保留上半身，白色背景")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
