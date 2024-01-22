package scene

import (
	"content_server/model/bg_img"
	model "content_server/model/scene"
	"content_server/redis"
	"content_server/setting"
	"content_server/utils/aiart"
	"content_server/utils/cos"
	"content_server/utils/dashscope"
	"content_server/utils/flake"
	"content_server/utils/processor"
	"encoding/json"
	"log"
	"strconv"
)

type Scene struct {
	Url      string `json:"url"`
	SceneId  string `json:"sceneId"`
	Username string `json:"username"`
}

//func GetChooseContentBySceneId(sceneId string) (chooseContent string, err error) {
//	key := "choose:" + sceneId
//	chooseContent = ""
//	if redis.Exists(key) {
//		chooseContent, err = redis.Get(key)
//		if err != nil {
//			return "", err
//		}
//	} else {
//		chooseContent, err = scene.GetChooseContentBySceneId(sceneId)
//		if err != nil {
//			return "", err
//		}
//		err = redis.Set(key, chooseContent)
//	}
//	return chooseContent, nil
//}
//
//func GetSceneContentBySceneId(sceneId int64) (sceneContent string, err error) {
//	key := "scene:" + strconv.FormatInt(sceneId, 10)
//	if redis.Exists(key) {
//		sceneContentBytes, err := redis.Get(key)
//		if err != nil {
//			return "", err
//		}
//		sceneContent = string(sceneContentBytes)
//		return sceneContent, nil
//	}
//	sceneContent = cos.GetSceneContent(sceneId)
//	return sceneContent, nil
//}

func GetDescBySceneId(sceneId string) (desc string, err error) {
	key := "desc:" + sceneId
	if redis.Exists(key) {
		desc, err = redis.Get(key)
		if err != nil {
			return "", err
		}
	}
	desc, err = model.GetDescBySceneId(sceneId)
	if err != nil {
		return "", err
	}
	redis.Set(key, desc, setting.ServerSetting.SceneExpire)
	return desc, nil
}

func updateSceneAfterGenerate(scene *model.Scene) error {
	sceneId := strconv.FormatInt(scene.SceneId, 10)
	log.Printf("update scene: %v", scene)
	err := model.UpdateSceneBySceneId(sceneId, scene)
	if err != nil {
		return err
	}
	sceneInfo := &Scene{
		SceneId:  sceneId,
		Username: scene.Creator,
		Url:      scene.COSUrl,
	}
	redis.Set("desc:"+sceneId, scene.ShortDesc, setting.ServerSetting.SceneExpire)
	jsonStr, err := json.Marshal(sceneInfo)
	if err == nil {
		log.Printf("jsonStr: %v", sceneInfo)
		redis.Set("scene:"+sceneId, string(jsonStr), setting.ServerSetting.SceneExpire)
		if scene.IsInit != 0 {
			isInitStr := strconv.Itoa(scene.IsInit)
			redis.Set("init:"+scene.Creator+isInitStr, string(jsonStr), setting.ServerSetting.SceneExpire)
		}
	}
	return nil
}

func GenerateScene(sceneId, sceneInfo string) {
	var scene = &model.Scene{}
	err := json.Unmarshal([]byte(sceneInfo), &scene)
	log.Printf("scene info: %v", scene)
	if err != nil {
		log.Printf("unmarshal scene info failed. err: %v", err)
		return
	}
	parentId := scene.ParentSceneId
	chooseContent := scene.ChooseContent
	sceneIdInt, err := strconv.ParseInt(sceneId, 10, 64)

	// 简述
	sceneDesc, err := GetDescBySceneId(strconv.FormatInt(parentId, 10))
	if err != nil {
		return
	}

	// 3. 开始生成
	// 3.1 生成新剧本简述
	descReq := dashscope.NewDescGenerationRequest(chooseContent, sceneDesc)
	rawDesc, err := descReq.GenerateText()
	if err != nil {
		log.Printf("generate desc failed. err: %v", err)
		return
	}
	// 3.2 生成背景
	bgDescReq := dashscope.NewBgPromptGenerationRequest(rawDesc)
	rawBgPrompt, err := bgDescReq.GenerateText()
	if err != nil {
		log.Printf("generate bg prompt failed. err: %v", err)
		return
	}
	bgUrl, _, err := generateBgImg(rawBgPrompt, sceneIdInt)
	// 3.3 生成新剧本
	req := dashscope.NewTextGenerationRequest(rawDesc)
	rawText, err := req.GenerateText()
	if err != nil {
		log.Printf("generate text failed. err: %v", err)
		return
	}
	// 3.4 生成分支选项
	chooseReq := dashscope.NewChooseGenerationRequest(rawDesc)
	rawChoose, err := chooseReq.GenerateText()
	if err != nil {
		log.Printf("generate choose failed. err: %v", err)
		return
	}
	// 3.4 解析
	newSceneContent, chooses := processor.GenerateNewScene(rawText, bgUrl, rawChoose)
	cosUrl, err := cos.UploadScene(newSceneContent, sceneId)
	if err != nil {
		log.Printf("upload scene failed. err: %v", err)
		return
	}
	// 5. 更新数据库
	scene.COSUrl = cosUrl
	scene.ShortDesc = rawDesc
	err = updateSceneAfterGenerate(scene)
	if err != nil {
		log.Printf("update scene failed. err: %v\n", err)
		return
	}
	// 6. 写入子剧本
	childScenes := make([]*model.Scene, 0)
	for k, v := range chooses {
		err, newScene := model.SaveUngeneratedScene(k, sceneIdInt, v, scene.Creator)
		if err != nil {
			log.Printf("save scene failed. err: %v", err)
			continue
		}
		childScenes = append(childScenes, newScene)
	}
	jsonStr, err := json.Marshal(childScenes)
	redis.Set("childs:"+sceneId, string(jsonStr), setting.ServerSetting.SceneExpire)
	redis.Delete("childs:" + strconv.FormatInt(scene.ParentSceneId, 10))
	//if scene.IsInit != 0 {
	//	redis.Set("init:" + scene.Creator + string(rune(scene.IsInit)), string(jsonStr), setting.ServerSetting.SceneExpire)
	//}
	if err != nil {
		return
	}

}

func generateBgImg(rawBgPrompt string, sceneId int64) (string, int64, error) {
	url := "https://fake-buddha-1300084664.cos.ap-shanghai.myqcloud.com/image%2F496806037978889767.jpg"
	bgId, err := flake.Generate()
	if err != nil {
		return url, bgId, err
	}
	imgBase54, err := aiart.Generate(rawBgPrompt)
	if err != nil {
		return url, bgId, err
	}
	// 上传COS
	url, err = cos.UploadImage(imgBase54, bgId)
	if err != nil {
		return url, bgId, err
	}
	// 保存至MySQL
	//bgImg := bg_img.NewBackgroundImg(key, url, sceneId, bgDesc)
	bgImg := &bg_img.BackgroundImg{
		ImgId:   bgId,
		COSUrl:  url,
		SceneId: sceneId,
		Prompt:  rawBgPrompt,
	}
	err = bg_img.SaveBackgroundImg(bgImg)
	return url, bgId, nil
}
