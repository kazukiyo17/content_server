package redis_mq

import "C"
import (
	service "content_server/service/scene"
	"fmt"
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
		msg, err3 := redisMQClient.GetMsgBlock(TEST_STREAM_KEY)
		if err3 != nil {
			fmt.Println("GetMsgByGroupConsumer Failed. err:", err3)
			continue
		}
		for k, v := range msg {
			fmt.Println("---------------------------------------------------")
			fmt.Println("get msg:", k, v)
			service.GenerateScene(k, v)
		}
	}
}

//func hasGenerated(sceneId int64) bool {
//	// 1 查询redis
//	rKey := "cos:" + strconv.FormatInt(sceneId, 10)
//	if redis.Exists(rKey) {
//		return true
//	}
//	// 2 查询db
//	cosUrl, err := model.GetCosUrlBySceneId(sceneId)
//	if err != nil {
//		return false
//	}
//	if cosUrl == "" {
//		return false
//	}
//	// 2. 存入redis
//	redis.Set(rKey, cosUrl)
//	return true
//}
