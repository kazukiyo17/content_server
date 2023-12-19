package bg_img

import (
	"content_server/model"
	"strconv"
)

type BackgroundImg struct {
	model.Model
	ImgId   int64  `gorm:"bigint;index"`
	COSUrl  string `gorm:"varchar(255)"`
	SceneId int64  `gorm:"bigint;index"`
	Prompt  string `gorm:"varchar(255)"`
}

func NewBackgroundImg(key string, url string, sceneId string, prompt string) *BackgroundImg {
	imgIdInt, _ := strconv.ParseInt(key, 10, 64)
	sceneIdInt, _ := strconv.ParseInt(sceneId, 10, 64)
	return &BackgroundImg{
		ImgId:   imgIdInt,
		COSUrl:  url,
		SceneId: sceneIdInt,
		Prompt:  prompt,
	}
}

func SaveBackgroundImg(b *BackgroundImg) (err error) {
	err = model.DB.Model(&BackgroundImg{}).Create(&b).Error
	return err
}
