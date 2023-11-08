package consumer

import "C"
import (
	"content_server/redis_mq"
	"content_server/service/scene_service"
	"encoding/json"
	"fmt"
)

type SceneMsgConsumer struct {
	// 你的消息队列客户端
	redisCli *redis_mq.RedisStreamMQClient
}

type ISceneMsgConsumer interface {
	consume(groupName string, consumerName string, msgAmount int32)
}

func (mq *SceneMsgConsumer) consume(groupName string, consumerName string, msgAmount int32) {
	for {
		// 取消息
		msg, err3 := mq.redisCli.GetMsgByGroupConsumer("", groupName, consumerName, msgAmount)
		if err3 != nil {
			fmt.Println("GetMsgByGroupConsumer Failed. err:", err3)
			return
		}
		// 消费消息
		for sceneId, val := range msg {
			//var scene *model.Scene
			err := json.Unmarshal(val, &scene)
			if err != nil {
				return
			}
			// 读取Scene
			sceneService := scene_service.Scene{SceneId: sceneId}
			scene, err := sceneService.GetSceneBySceneId()
			// 调用NLP接口

			// 后处理

			// 写入MySQL

			// 写入Redis
		}
	}
}
