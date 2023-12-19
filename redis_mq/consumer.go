package redis_mq

import "C"
import (
	"content_server/model/scene"
	"content_server/redis"
	"content_server/service/scene_service"
	"fmt"
	"strconv"
)

type SceneMsgConsumer struct {
	// 你的消息队列客户端
	redisCli *RedisStreamMQClient
}

type ISceneMsgConsumer interface {
	consume(groupName string, consumerName string, msgAmount int32)
}

func Consume() {
	for {
		// 取消息
		msg, err3 := redisMQClient.GetMsgByGroupConsumer(TEST_STREAM_KEY, TEST_GROUP_NAME, TEST_CONSUMER_NAME)
		if err3 != nil {
			fmt.Println("GetMsgByGroupConsumer Failed. err:", err3)
			continue
		}
		for k, v := range msg {
			kInt, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}
			if hasGenerated(kInt) {
				continue
			}
			scene_service.GenerateScene(k, v)
		}
	}
}

func hasGenerated(sceneId int64) bool {
	// 1 查询redis
	rKey := "cos:" + strconv.FormatInt(sceneId, 10)
	if redis.Exists(rKey) {
		return true
	}
	// 2 查询db
	cosUrl, err := scene.GetCosUrlBySceneId(sceneId)
	if err != nil {
		return false
	}
	if cosUrl == "" {
		return false
	}
	// 2. 存入redis
	err = redis.Set(rKey, cosUrl)
	return true
}
