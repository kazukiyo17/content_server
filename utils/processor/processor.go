package processor

import (
	"fmt"
	"strings"
)

// 写死的角色名字["赵卓群", "周远山", "林秀芝", "陈大志"]
var roleNames = []string{"赵卓群", "周远山", "林秀芝", "陈大志"}

func FormatSceneContent(content string, imgUrls []string) (newContent string, err error) {
	// 逐行读取
	bgIndex := 0
	newContentList := []string{}
	currentRole := "我" // 标记当前的角色
	sentences := strings.Split(content, "\n")
	for _, sentence := range sentences {
		// 如果存在冒号
		if strings.Contains(sentence, "：:") {
			// 第一个冒号前的内容为角色名字,冒号可能有多个，不能用split
			key := sentence[:strings.Index(sentence, "：:")]
			chatContent := sentence[strings.Index(sentence, "：:")+1:]
			// 如果是choose
			if strings.HasPrefix(key, "choose") {
				chooseContent := parseChoose(chatContent)
				newContentList = append(newContentList, chooseContent)
				continue
			} else if strings.HasPrefix(key, "changeBg") {
				bgContent := parseChangeBg(imgUrls[bgIndex])
				newContentList = append(newContentList, bgContent)
				bgIndex++
				continue
			}
			// 不为choose则为角色对话
			// 切换立绘
			figure := parseChangeFigure(currentRole, key)
			if figure != "" {
				newContentList = append(newContentList, figure)
				currentRole = key
			}
			// 生成对话
			chat := fmt.Sprintf("%s:%s;", key, chatContent)
			newContentList = append(newContentList, chat)
		} else {
			chat := fmt.Sprintf("%s;", sentence)
			newContentList = append(newContentList, chat)
		}
	}
	newContent = strings.Join(newContentList, "\n")
	return newContent, nil
}

func parseChangeFigure(currentRole string, roleName string) (figure string) {
	// roleName以"我"开头或在roleNames中
	if strings.HasPrefix(roleName, "我") || currentRole == roleName {
		return
	}
	for _, name := range roleNames {
		if strings.HasPrefix(roleName, name) {
			figureFile := name + ".jpg"
			figure = fmt.Sprintf("changeFigure:%s -left -enter=enter-from-left -next;", figureFile)
			return figure
		}
	}
	return
}

func parseChoose(content string) (newContent string) {
	// 按|或｜分隔
	chooseList := strings.Split(content, "|｜")
	// 如果大于3个选项，选前三个
	if len(chooseList) > 3 {
		chooseList = chooseList[:3]
	}
	// 生成choose
	choose := fmt.Sprintf("choose:%s;", strings.Join(chooseList, "|"))
	return choose
}

func parseChangeBg(imgUrl string) (content string) {
	content = fmt.Sprintf("changeBg:%s -next;", imgUrl)
	return content
}
func ExtractBgDesc(content string) (bgDesc []string) {
	// 逐行读取
	sentences := strings.Split(content, "\n")
	for _, sentence := range sentences {
		// 如果存在冒号
		if strings.Contains(sentence, "：:") {
			// 第一个冒号前的内容为角色名字,冒号可能有多个，不能用split
			key := sentence[:strings.Index(sentence, "：:")]
			chatContent := sentence[strings.Index(sentence, "：:")+1:]
			// 如果是changeBg
			if strings.HasPrefix(key, "changeBg") {
				bgDesc = append(bgDesc, chatContent)
			}
		}
	}
	return bgDesc
}
