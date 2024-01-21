package scene

import (
	"content_server/model"
)

type Scene struct {
	model.Model
	SceneId       int64  `json:"scene_id" gorm:"type:bigint;index"`
	ChooseContent string `json:"choose_content" gorm:"type:varchar(255)"`
	Creator       string `json:"creator" gorm:"type:varchar(255)"`
	ParentSceneId int64  `json:"parent_scene_id" gorm:"type:bigint;index"`
	COSUrl        string `json:"cos_url" gorm:"type:varchar(255)"`
	ShortDesc     string `json:"desc" gorm:"type:varchar(600)"`
	IsInit        int    `json:"is_init" gorm:"type:int(11)"`
}

func GetChooseContentBySceneId(sceneId string) (chooseContent string, err error) {
	scene := &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	chooseContent = scene.ChooseContent
	return chooseContent, nil
}

func SaveScene(scene *Scene) (err error) {
	err = model.DB.Model(&Scene{}).Create(&scene).Error
	return
}

func SaveUngeneratedScene(sceneId, parentSceneId int64, choose, username string) (error, *Scene) {

	scene := &Scene{
		SceneId:       sceneId,
		ChooseContent: choose,
		Creator:       username,
		ParentSceneId: parentSceneId,
	}
	err := model.DB.Model(&Scene{}).Create(&scene).Error
	return err, scene
}

func UpdateSceneBySceneId(sceneId string, scene *Scene) (err error) {
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).Updates(&scene).Error
	return
}

func GetCosUrlBySceneId(sceneId int64) (cosUrl string, err error) {
	scene := &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	cosUrl = scene.COSUrl
	return cosUrl, nil
}

func GetDescBySceneId(sceneId string) (desc string, err error) {
	scene := &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	desc = scene.ShortDesc
	return desc, nil
}

//func GetParentSceneIdBySceneId(sceneId string) (parentSceneId int64, err error) {
//	scene := &ModelScene{}
//	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).First(&scene).Error
//	if err != nil {
//		return 0, err
//	}
//	parentSceneId = scene.ParentSceneId
//	return parentSceneId, nil
//}
