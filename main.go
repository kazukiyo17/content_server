package main

import (
	"content_server/model"
	"content_server/redis_mq"
	"content_server/setting"
)

func init() {
	setting.Setup()
	model.Setup()
	redis_mq.Setup()
}

func main() {

}
