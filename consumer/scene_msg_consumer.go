package consumer

import "C"
import (
	"content_server/redis_mq"
	"fmt"
)

type SceneMsgConsumer struct {
	// 你的消息队列客户端
	redisCli *redis_mq.RedisStreamMQClient
}

func (mq *SceneMsgConsumer) consume(groupName string, consumerName string, msgAmount int32) {
	for {
		// 取消息
		msg, err3 := mq.redisCli.GetMsgByGroupConsumer("", groupName, consumerName, msgAmount)
		if err3 != nil {
			fmt.Println("GetMsgByGroupConsumer Failed. err:", err3)
			return
		}
		// 消费消息: 1. 从Cache/db中读取Scene的Prompt 2. 调用NLP接口 3. 将结果写入Cache/db
		//	msg 是map[string]map[string][]string
		for _, v := range msg {
			for _, v2 := range v {
				for _, v3 := range v2 {
					fmt.Println("msg:", v3)
				}
			}
		}
	}
}
