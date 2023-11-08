package scene_service

import (
	"content_server/model"
	"content_server/redis"
	"content_server/utils/nlp"
	"encoding/json"
)

type Scene struct {
	SceneId       int
	Prompt        string
	CreatorId     int
	CreateTime    int
	ParentSceneId int
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

func (s *Scene) GnerateSceneNLP() (err error) {
	text := ""
	res, err := nlp.nlpRequest(text)
	if err != nil {
		return err
	}

}
