package consumer

import (
	"content_server/redis_mq"
	"fmt"
	"time"
)

func consume(redisCli *redis_mq.RedisStreamMQClient, groupName string, consumerName string, msgAmount int32) {
	startTime := time.Now()
	fmt.Println("Start Test Function GetMsgByGroupConsumer")
	msgMap3, err3 := redisCli.GetMsgByGroupConsumer(common.TEST_STREAM_KEY, groupName, consumerName, msgAmount)

	if err3 != nil {
		fmt.Println("GetMsgByGroupConsumer Failed. err:", err3)
		return
	}
	fmt.Println("End Test Function GetMsgByGroupConsumer")

	fmt.Println("Start Test Function ReplyAck")
	key2msgIds3, _ := PrintMsgMap(msgMap3)
	for streamKey, vecMsgId := range key2msgIds3 {
		//fmt.Println("streamKey:", streamKey, "groupName:", groupName, "consumerName:", consumerName, "Ack Msg Count:", msgCount)
		err3 = redisCli.ReplyAck(streamKey, groupName, vecMsgId)
		if err3 != nil {
			fmt.Println("ReplyAck Failed. err:", err3)
		}
	}

	fmt.Println("End Test Function ReplyAck")
	costTime := time.Since(startTime)
	fmt.Println("=========:START TIME:", startTime)
	fmt.Println("=========:COST TIME:", costTime)
}

func PrintMsgMap(msgMap map[string]map[string][]string) (key2msgIds map[string][]string, msgCount int32) {
	key2msgIds = make(map[string][]string, 0)
	msgCount = 0
	for streamKey, val := range msgMap {
		//fmt.Println("StreamKey:", streamKey)
		vecMsgId := make([]string, 0)
		for msgId, msgList := range val {
			//fmt.Println("MsgId:", msgId)
			vecMsgId = append(vecMsgId, msgId)
			for msgIndex := 0; msgIndex < len(msgList); msgIndex = msgIndex + 2 {
				//var msgKey = msgList[msgIndex]
				//var msgVal = msgList[msgIndex+1]
				msgCount++
				//fmt.Println("MsgKey:", msgKey, "MsgVal:", msgVal)
			}
		}
		key2msgIds[streamKey] = vecMsgId
	}
	return
}
