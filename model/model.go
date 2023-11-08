package model

import (
	"content_server/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Setup() {
	_dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User, setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host, setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name)
	d, err := gorm.Open(mysql.Open(_dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	DB = d
}
