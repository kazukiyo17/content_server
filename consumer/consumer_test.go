package consumer

import (
	"content_server/redis_mq"
	"fmt"
	"time"
)

const (
	STREAM_MQ_MAX_LEN  = 500000 //消息队列最大长度
	READ_MSG_AMOUNT    = 1000   //每次读取消息的条数
	READ_MSG_BLOCK_SEC = 30     //阻塞读取消息时间
	TEST_STREAM_KEY    = "TestStreamKey1"
)

func PrintMsgMap(msgMap map[string]map[string][]string) (key2msgIds map[string][]string, msgCount int32) {
	key2msgIds = make(map[string][]string, 0)
	msgCount = 0
	for streamKey, val := range msgMap {
		fmt.Println("StreamKey:", streamKey)
		vecMsgId := make([]string, 0)
		for msgId, msgList := range val {
			fmt.Println("MsgId:", msgId)
			vecMsgId = append(vecMsgId, msgId)
			for msgIndex := 0; msgIndex < len(msgList); msgIndex = msgIndex + 2 {
				var msgKey = msgList[msgIndex]
				var msgVal = msgList[msgIndex+1]
				msgCount++
				fmt.Println("MsgKey:", msgKey, "MsgVal:", msgVal)
			}
		}
		key2msgIds[streamKey] = vecMsgId
	}
	return
}

// 非阻塞消费
func testGroupConsumer(redisCli *redis_mq.RedisStreamMQClient, groupName string, consumerName string, msgAmount int32) {

	/*
		fmt.Println("Start Test Function CreateConsumerGroup")
		err3 := redisCli.CreateConsumerGroup(common.TEST_STREAM_KEY, groupName, "0")
		if err3 != nil {
			fmt.Println("CreateConsumerGroup Failed. err:", err3)
			return
		}
		fmt.Println("End Start Test Function CreateConsumerGroup")
	*/
	startTime := time.Now()
	fmt.Println("Start Test Function GetMsgByGroupConsumer")
	msgMap3, err3 := redisCli.GetMsgByGroupConsumer(TEST_STREAM_KEY, groupName, consumerName, msgAmount)

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
