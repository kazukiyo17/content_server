package conf

import "strings"

type RoleName string

const (
	ZHOU RoleName = "周远山"
	ZHAO RoleName = "赵卓群"
	LIN  RoleName = "林秀芝"
	CHEN RoleName = "陈大志"
	ME   RoleName = "我"
)

// 文件名Map
var FigureFileMap = map[RoleName]string{
	ZHOU: "zhou.jpg",
	ZHAO: "zhao.jpg",
	LIN:  "lin.jpg",
	CHEN: "chen.jpg",
	ME:   "me.jpg",
}

func GetRoleName(content string) (roleName RoleName) {
	// 按开头
	if strings.HasPrefix(content, string(ZHOU)) {
		return ZHOU
	} else if strings.HasPrefix(content, string(ZHAO)) {
		return ZHAO
	} else if strings.HasPrefix(content, string(LIN)) {
		return LIN
	} else if strings.HasPrefix(content, string(CHEN)) {
		return CHEN
	} else if strings.HasPrefix(content, string(ME)) {
		return ME
	}
	return ""
}
