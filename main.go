package main

import (
	"content_server/model"
	"content_server/redis"
	"content_server/redis_mq"
	"content_server/setting"
	"content_server/utils/cos"
)

func init() {
	setting.Setup()
	//model.Setup()
	redis_mq.Setup()
	redis.Setup()
	model.Setup()
	cos.Setup()
}

func main() {
	// 测试aiart
	//aiart.Test()
	// 测试dashscope
	//dashscope.Test()
	redis_mq.Consume()
	//test()
	//processor.Test()
}

// Test 测试
//func test() {
//	rawDesc := "我是村里的一名普通青年，生活平静。一日，诡异男子赵卓群突然出现，他似乎知道村中隐藏的秘密。在调查过程中，我结识了村民女孩林秀芝。"
//	key, _, _ := cos.GenerateSceneCosPath()
//	// 3.2 生成新剧本
//	//req := dashscope.NewTextGenerationRequest(rawDesc)
//	//rawText, err := req.GenerateText()
//	//if err != nil {
//	//	return
//	//}
//	// 读取test.txt
//	file, err := os.Open("test.txt")
//	if err != nil {
//		return
//	}
//	defer file.Close()
//	buf := make([]byte, 1024)
//	rawText := ""
//	for {
//		n, err := file.Read(buf)
//		rawText += string(buf[:n])
//		if err != nil {
//			break
//		}
//	}
//
//	// 4. 解析
//	newSceneContent, chooses, err := processor.GenerateSceneContent(rawText)
//	if err != nil {
//		return
//	}
//	// 4. 上传cos
//	// chooseIdInt转string
//	cosUrl, err := cos.UploadScene(newSceneContent, key)
//	if err != nil {
//		return
//	}
//	// 5. 更新数据库: CreatorId CreateTime COSKey COSUrl
//	sceneIdInt, _ := strconv.ParseInt(key, 10, 64)
//	newScene := &scene.ModelScene{
//		SceneId:       sceneIdInt,
//		ChooseContent: "",
//		CreatorId:     0,
//		ParentSceneId: 0,
//		COSUrl:        cosUrl,
//		ShortDesc:     rawDesc,
//	}
//	err = scene.SaveScene(newScene)
//	// map[string][]string
//	for k, v := range chooses {
//		chooseIdInt, _ := strconv.ParseInt(k, 10, 64)
//		newScene := &scene.ModelScene{
//			SceneId:       chooseIdInt,
//			ChooseContent: v[1],
//			CreatorId:     0,
//			ParentSceneId: sceneIdInt,
//		}
//		err = scene.SaveScene(newScene)
//	}
//	// 6. 写入redis
//	childSceneIds := make([]string, len(chooses))
//	for k, v := range chooses {
//		childSceneIds = append(childSceneIds, k)
//		err := redis.Set("choose:"+k, v[1])
//		if err != nil {
//			return
//		}
//	}
//	err = redis.Set("child:"+key, strings.Join(childSceneIds, ","))
//	err = redis.Set("scene:"+key, newSceneContent)
//	err = redis.Set("cos:"+key, cosUrl)
//	err = redis.Set("desc:"+key, rawDesc)
//	if err != nil {
//		return
//	}
//}
