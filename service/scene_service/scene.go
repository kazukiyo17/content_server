package scene_service

import (
	"content_server/model/scene"
	"content_server/redis"
	"content_server/utils/cos"
	"content_server/utils/dashscope"
	"content_server/utils/processor"
	"strconv"
	"strings"
)

func GetChooseContentBySceneId(sceneId string) (chooseContent string, err error) {
	key := "choose:" + sceneId
	chooseContent = ""
	if redis.Exists(key) {
		chooseContent, err = redis.Get(key)
		if err != nil {
			return "", err
		}
	} else {
		chooseContent, err = scene.GetChooseContentBySceneId(sceneId)
		if err != nil {
			return "", err
		}
		err = redis.Set(key, chooseContent)
	}
	return chooseContent, nil
}

func GetSceneContentBySceneId(sceneId int64) (sceneContent string, err error) {
	key := "scene:" + strconv.FormatInt(sceneId, 10)
	if redis.Exists(key) {
		sceneContentBytes, err := redis.Get(key)
		if err != nil {
			return "", err
		}
		sceneContent = string(sceneContentBytes)
		return sceneContent, nil
	}
	sceneContent = cos.GetSceneContent(sceneId)
	return sceneContent, nil
}

func GetDescBySceneId(sceneId string) (desc string, err error) {
	key := "desc:" + sceneId
	if redis.Exists(key) {
		desc, err = redis.Get(key)
		if err != nil {
			return "", err
		}
	} else {
		desc, err = scene.GetDescBySceneId(sceneId)
		if err != nil {
			return "", err
		}
		err = redis.Set(key, desc)
	}
	return desc, nil
}

func GenerateScene(chooseId string, sceneId string) {
	chooseIdInt, err := strconv.ParseInt(chooseId, 10, 64)
	if err != nil {
		return
	}
	_, err = strconv.ParseInt(sceneId, 10, 64)
	if err != nil {
		return
	}
	// 1. 根据k取Choose文本内容
	chooseContent, err := GetChooseContentBySceneId(chooseId)
	if err != nil {
		return
	}
	// 2. 根据v取剧本简述
	sceneDesc, err := GetDescBySceneId(sceneId)
	if err != nil {
		return
	}
	// 3.1 生成新剧本简述
	descReq := dashscope.NewDescGenerationRequest(chooseContent, sceneDesc)
	rawDesc, err := descReq.GenerateText()
	if err != nil {
		return
	}
	// 3.2 生成新剧本
	// rawDesc 如果分行了，去掉为空的行，取第一行
	//rawDesc = strings.Split(rawDesc, "\n")[0]
	req := dashscope.NewTextGenerationRequest(rawDesc)
	rawText, err := req.GenerateText()
	if err != nil {
		return
	}
	// 4. 解析
	newSceneContent, chooses, err := processor.GenerateSceneContent(rawText, chooseId)
	if err != nil {
		return
	}
	// 4. 上传cos
	cosUrl, err := cos.UploadScene(newSceneContent, chooseId)
	if err != nil {
		return
	}
	// 5. 更新数据库: CreatorId CreateTime COSKey COSUrl
	newScene := &scene.ModelScene{
		//SceneId:       chooseIdInt,
		//CreatorId:     0,
		//ParentSceneId: sceneIdInt,
		COSUrl:    cosUrl,
		ShortDesc: rawDesc,
	}
	err = scene.UpdateSceneBySceneId(chooseId, newScene)
	// 6. 写入redis
	childSceneIds := make([]string, 0)
	for k, v := range chooses {
		// k 转 int64
		kInt, _ := strconv.ParseInt(k, 10, 64)
		newChildScene := &scene.ModelScene{
			SceneId:       kInt,
			CreatorId:     0,
			ParentSceneId: chooseIdInt,
			ChooseContent: v[1],
		}
		err = scene.SaveScene(newChildScene)
		if err != nil {
			continue
		}
		childSceneIds = append(childSceneIds, k)
		err := redis.Set("choose:"+k, v[1])
		if err != nil {
			return
		}
	}
	err = redis.Set("child:"+chooseId, strings.Join(childSceneIds, ","))
	err = redis.Set("scene:"+chooseId, newSceneContent)
	err = redis.Set("cos:"+chooseId, cosUrl)
	err = redis.Set("desc:"+chooseId, rawDesc)
	if err != nil {
		return
	}

}
