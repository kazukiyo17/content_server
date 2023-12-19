package scene

import (
	"content_server/model"
)

type ModelScene struct {
	model.Model
	SceneId       int64  `json:"tag_id" gorm:"type:bigint;index"`
	ChooseContent string `json:"choose_content" gorm:"type:varchar(255)"`
	CreatorId     int64  `json:"creator_id" gorm:"type:bigint;index"`
	ParentSceneId int64  `json:"parent_scene_id" gorm:"type:bigint;index"`
	COSUrl        string `json:"cos_url" gorm:"type:varchar(255)"`
	ShortDesc     string `json:"desc" gorm:"type:varchar(600)"`
}

func GetChooseContentBySceneId(sceneId string) (chooseContent string, err error) {
	scene := &ModelScene{}
	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	chooseContent = scene.ChooseContent
	return chooseContent, nil
}

func SaveScene(scene *ModelScene) (err error) {
	err = model.DB.Model(&ModelScene{}).Create(&scene).Error
	return
}

func UpdateSceneBySceneId(sceneId string, scene *ModelScene) (err error) {
	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).Updates(&scene).Error
	return
}

func GetCosUrlBySceneId(sceneId int64) (cosUrl string, err error) {
	scene := &ModelScene{}
	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	cosUrl = scene.COSUrl
	return cosUrl, nil
}

func GetDescBySceneId(sceneId string) (desc string, err error) {
	scene := &ModelScene{}
	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	desc = scene.ShortDesc
	return desc, nil
}
