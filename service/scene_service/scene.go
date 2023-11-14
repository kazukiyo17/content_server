package scene_service

import (
	"content_server/model"
	"content_server/redis"
	"content_server/utils/aiart"
	"content_server/utils/cos"
	"content_server/utils/dashscope"
	"content_server/utils/processor"
	"encoding/json"
)

type Scene struct {
	SceneId       int
	Prompt        string
	CreatorId     int
	CreateTime    int
	ParentSceneId int
	Chooses       []*Choose
	Content       string
}

type Choose struct {
	SceneId int
	Content string
	Key     string
}

// interface
type SceneServiceInterface interface {
	GetSceneBySceneId() (scene *model.Scene, err error)
	GetSceneByParentSceneId() (scenes []*model.Scene, err error)
	GetRedisKeyBySceneId() (redisKey string)
	// 调用NLP接口
	GnerateSceneNLP() (err error)
}

func (s *Scene) GetSceneBySceneId() (scene *model.Scene, err error) {
	var cacheScene *model.Scene
	// 1. Redis中读取
	redisKey := s.GetRedisKeyBySceneId()
	if redis.Exists(redisKey) {
		// 1.1 Redis中存在，直接返回
		data, err := redis.Get(redisKey)
		if err != nil {
			return
		} else {
			err := json.Unmarshal(data, &cacheScene)
			if err != nil {
				return nil, err
			}
			return cacheScene, nil
		}
	}
	// 2. Redis中不存在，从MySQL中读取
	scene = &model.Scene{SceneId: s.SceneId}
	scene, err = scene.GetSceneBySceneId()
	if err != nil {
		return nil, err
	}
	// 3. 将MySQL中读取的数据写入Redis
	data, err := json.Marshal(scene)
	if err != nil {
		return nil, err
	}
	err = redis.Set(redisKey, data, 3600)
	if err != nil {
		return nil, err
	}
	return scene, nil
}

func (s *Scene) GetRedisKeyBySceneId() (redisKey string) {
	redisKey = "scene:" + string(rune(s.SceneId))
	return redisKey
}

func (s *Scene) GetSceneByParentSceneId() (scenes []*model.Scene, err error) {
	scene := &model.Scene{ParentSceneId: s.ParentSceneId}
	scenes, err = scene.GetSceneByParentSceneId()
	if err != nil {
		return nil, err
	}
	return scenes, nil
}

// GenerateScene 生成分支的剧本
func (s *Scene) GenerateScene() (err error) {
	// 分支
	chooses := s.Chooses
	for _, choose := range chooses {
		// 生成分支的剧本
		req := dashscope.NewTextGenerationRequest(choose.Content, s.Content)
		response, err := req.GenerateText()
		if err != nil {
			return err
		}
		content := response.Text // 剧本元数据
		// 后处理：1. 提取背景描述 2. 格式化剧本
		bgDescList := processor.ExtractBgDesc(content)
		// 生成并上传
		imgUrls := []string{}
		for _, bgDesc := range bgDescList {
			if bgDesc == "" {
				continue
			}
			imgBase54, err := aiart.Generate(bgDesc)
			if err != nil {
				return err
			}
			// 上传COS, 随机一个key
			url, err := cos.UploadImage(imgBase54)
			if err != nil {
				return err
			}
			imgUrls = append(imgUrls, url)
		}
		newContent, err := processor.FormatSceneContent(content, imgUrls)
		if err != nil {
			return err
		}
		// 剧本上传COS

	}
}
