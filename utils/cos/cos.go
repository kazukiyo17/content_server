package cos

import (
	"content_server/setting"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sony/sonyflake"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var client *cos.Client
var flake *sonyflake.Sonyflake

func Setup() {
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
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

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

func UploadScene(content string, key string) (cosUrl string, err error) {
	//client := newClient()
	//key, err := generateKey()
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
	objectUrl := GetObjectUrl(name)
	return objectUrl.String(), nil
}

// UploadImage 将base64编码的图片上传到cos
func UploadImage(base64Str string) (url string, key string, err error) {
	//client := newClient()
	// 生成随机key
	key, err = generateKey()
	if err != nil {
		panic(err)
		return "", key, nil
	}
	// content 为 base64 编码的图片
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		panic(err)
	}
	name := "image/" + key + ".jpg"
	// 上传字节流
	f := strings.NewReader(string(data))
	_, err = client.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}
	// 获取图片URL
	imgUrl := GetObjectUrl(name)
	return imgUrl.String(), key, nil
}

// GetSceneContent 下载cos中的文件
func GetSceneContent(sceneId int64) (res string) {
	//client := newClient()
	key := strconv.FormatInt(sceneId, 10)
	path := "scene/" + key + ".txt"
	resp, err := client.Object.Get(context.Background(), path, nil)
	if err != nil {
		//panic(err)
		return ""
	}
	// buf := new(bytes.Buffer)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		res += string(buf[:n])
		if err != nil {
			break
		}
	}
	return res
}

// 获取cos中的文件URL
func GetObjectUrl(path string) (res *url.URL) {
	//client := newClient()
	ourl := client.Object.GetObjectURL(path)
	return ourl
}

// 读取cos中的文件
func getScene(key string) (res string) {
	//client := newClient()
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

func GenerateSceneCosPath() (string, string, error) {
	// 生成随机key
	key, err := generateKey()
	if err != nil {
		panic(err)
		return "", "", err
	}
	u := "https://fake-buddha-1300084664.cos.ap-shanghai.myqcloud.com/scene/" + key + ".txt"
	return key, u, nil
}

// 生成随机key
func generateKey() (string, error) {
	//flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	// 等待一毫秒
	time.Sleep(time.Millisecond)
	if err != nil {
		return "", err
	}
	// 控制台输出ID
	fmt.Println(id)
	return strconv.FormatInt(int64(id), 10), nil
}
