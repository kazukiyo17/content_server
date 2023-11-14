package cos

import (
	"content_server/setting"
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/sony/sonyflake"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func newClient() (client *cos.Client) {
	u, _ := url.Parse("https://fake-buddha-1300084664.cos.ap-shanghai.myqcloud.com")
	// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com
	su, _ := url.Parse("https://cos.ap-shanghai.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  setting.TencentCloudSetting.SecretID,
			SecretKey: setting.TencentCloudSetting.SecretKey,
		},
	})
	return client
}

func UploadScene(content string) (resp string, err error) {
	client := newClient()
	key, err := generateKey()
	if err != nil {
		panic(err)
		return
	}
	name := "scene/" + key + ".txt"
	f := strings.NewReader(content)
	_, err = client.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}
	return
}

// UploadImage 将base64编码的图片上传到cos
func UploadImage(base64Str string) (url string, err error) {
	client := newClient()
	// content 为 base64 编码的图片
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		panic(err)
	}
	// 生成随机key
	key, err := generateKey()
	if err != nil {
		panic(err)
		return
	}
	name := "image/" + key + ".jpg"
	// 上传字节流
	f := strings.NewReader(string(data))
	_, err = client.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}
	// 获取图片URL
	imgUrl := getObjectUrl(name)
	return imgUrl.String(), nil
}

// 获取cos中的文件URL
func getObjectUrl(path string) (res *url.URL) {
	client := newClient()
	ourl := client.Object.GetObjectURL(path)
	return ourl
}

// 读取cos中的文件
func getScene(key string) (res string) {
	client := newClient()
	path := "scene/" + key + ".txt"
	resp, err := client.Object.Get(context.Background(), path, nil)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			break
		}
		res += string(buf[:n])
	}
	return res
}

// 生成随机key
func generateKey() (string, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString([]byte(strconv.FormatUint(id, 10))), nil
}
