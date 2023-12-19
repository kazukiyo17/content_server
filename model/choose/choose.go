package choose

type Choose struct {
	ChooseId      int    `gorm:"primary_key" json:"id"`
	SceneId       int64  `gorm:"index" json:"sceneId"`
	ParentSceneId int64  `gorm:"index" json:"parentSceneId"`
	Content       string `gorm:"type:varchar(255)" json:"content"`
	COSUrl        string `gorm:"type:varchar(255)" json:"cosUrl"`
}
