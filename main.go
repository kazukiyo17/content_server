package main

import (
	"content_server/setting"
	"content_server/utils/dashscope"
)

func init() {
	setting.Setup()
	//model.Setup()
	//redis_mq.Setup()
}

func main() {
	// 测试aiart
	//aiart.Test()
	// 测试dashscope
	dashscope.Test()
}
