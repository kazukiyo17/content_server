package model

type Scene struct {
	SceneId       int    `json:"tag_id" gorm:"index"`
	Prompt        string `json:"prompt" gorm:"type:varchar(255)"`
	CreatorId     int    `json:"creator_id" gorm:"index"`
	CreateTime    int    `json:"create_time" gorm:"index"`
	ParentSceneId int    `json:"parent_scene_id" gorm:"index"`
}

// interface
type SceneModelInterface interface {
	GetSceneBySceneId(sceneId int) (scene Scene, err error)
	GetSceneByParentSceneId(parentSceneId int) (scenes []Scene, err error)
}

func (s *Scene) GetSceneBySceneId(sceneId int) (scene Scene, err error) {
	err = DB.Where("scene_id = ?", sceneId).First(&scene).Error
	return
}

func (s *Scene) GetSceneByParentSceneId(parentSceneId int) (scenes []Scene, err error) {
	err = DB.Where("parent_scene_id = ?", parentSceneId).Find(&scenes).Error
	return
}
