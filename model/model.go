package model

import (
	"content_server/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedAt  int `json:"created_on"`
	UpdatedAt int `json:"modified_on"`
	DeletedAt  int `json:"deleted_on"`
}

func Setup() {
	_dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User, setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host, setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name)
	d, err := gorm.Open(mysql.Open(_dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := d.DB()
	if err != nil {
		panic("failed to get db")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = d
}
