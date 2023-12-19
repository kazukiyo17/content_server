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

func Generate(text string) (base64Img string, err error) {
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
		Resolution: common.StringPtr("1280:720"),
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
		return "", err
	}
	return *response.Response.ResultImage, nil
}

//
//// 测试
//func Test() {
//	resp, err := Generate("(character design), a young man with short black hair and determined eyes," +
//		"casual wear,backpack,hands behind his back, cartoonish lithographs, clean background, " +
//		"super details,  --ar 3:4 --niji 5")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(resp)
//}
