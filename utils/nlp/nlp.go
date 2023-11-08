package nlp

import (
	"content_server/setting"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
)

func nlpRequest(text string) (resp string, err error) {
	credential := common.NewCredential(
		setting.NlpSetting.SecretID,
		setting.NlpSetting.SecretKey,
	)
	// 实例化一个客户端配置对象，可以指定超时时间等配置
	cpf := profile.NewClientProfile()
	// 实例化要请求产品的client对象
	client, _ := nlp.NewClient(credential, regions.Shanghai, cpf)
	// 实例化一个请求对象
	request := nlp.NewTextWritingRequest()
	request.Text = &text
	request.SourceLang = common.StringPtr("zh")
	request.Number = common.Int64Ptr(1)
	request.Style = common.StringPtr("urban_officialdom")
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.TextWriting(request)
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
