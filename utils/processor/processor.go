package processor

import (
	"content_server/conf"
	"content_server/utils/flake"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//func GetDefaultEndScene() string {
//	return conf.EndScene
//}
//
//func FormatChooseContent(chooseContent string) (string, map[string][]string) {
//	var chooseMap = make(map[string][]string)
//	chooseContent = strings.ReplaceAll(chooseContent, "\n\n", "\n")
//	match := regexp.MustCompile(`\d+\.(.*)`)
//	sentences := strings.Split(chooseContent, "\n")
//	for _, sentence := range sentences {
//		// 提取内容
//		sentence = strings.ReplaceAll(sentence, " ", "")
//		choose := match.FindStringSubmatch(sentence)
//		key, url, _ := cos.GenerateSceneCosPath() // 生成随机key
//		chooseMap[key] = []string{url, choose[1]}
//	}
//	return parseChoose(chooseMap), chooseMap
//}

//func GenerateSceneContent(content string, sceneId string) (newContent string, chooses map[string][]string, err error) {
//	var newContentList []string
//	var chooseMap = make(map[string][]string)
//	currentRole := conf.ME // 标记当前的角色
//	sentences := strings.Split(content, "\n")
//	// 正则匹配：任意字符+数字+.或:+任意字符
//	r := regexp.MustCompile(`^.*\d+[.:](.*)$`)
//	for idx := 0; idx < len(sentences); idx++ {
//		sentence := sentences[idx]
//		// 将sentence 中空格去掉，中文冒号替换为英文冒号
//		sentence = strings.ReplaceAll(sentence, " ", "")
//		sentence = strings.ReplaceAll(sentence, "：", ":")
//		role := conf.GetRoleName(sentence)
//		matches := r.FindStringSubmatch(sentence)
//		if len(sentence) == 0 {
//			continue
//		}
//		// 如果是冒号分隔的句子：1.changeBg 2.对话 3. choose
//		if strings.Contains(sentence, ":") {
//			prefix := strings.Split(sentence, ":")[0]
//			suffix := strings.Split(sentence, ":")[1]
//			prefix = strings.ToLower(prefix)
//			prefix = strings.ReplaceAll(prefix, " ", "")
//			if prefix == "changebg" || strings.HasPrefix(sentence, "场景") {
//				// 如果suffix为空，则为下一行
//				if suffix == "" {
//					idx++
//					suffix = sentences[idx]
//				}
//				// 生成并上传Bg图片
//				//bgUrl := "https://fake-buddha-1300084664.cos.ap-shanghai.myqcloud.com/image%2F490302589137635912.jpg"
//				bgUrl, _, err := generateBg(suffix, sceneId)
//				if err != nil {
//					return "", nil, err
//				}
//				// 生成changeBg
//				newContentList = append(newContentList, parseChangeBg(bgUrl))
//			} else if len(matches) > 1 {
//				for _, s := range sentences[idx:] {
//					key, url, _ := cos.GenerateSceneCosPath()
//					sMatches := r.FindStringSubmatch(s)
//					if len(sMatches) <= 1 {
//						break
//					}
//					// [url, sMatches[1]]
//					chooseMap[key] = []string{url, sMatches[1]}
//				}
//				break
//			} else if role != "" {
//				// 对话
//				if role != "" && role != currentRole && role != conf.ME {
//					currentRole = role
//					newContentList = append(newContentList, parseChangeFigure(role, currentRole))
//				}
//				chatContent := fmt.Sprintf("%s:%s", role, suffix)
//				newContentList = append(newContentList, chatContent)
//				currentRole = role
//				continue
//			}
//		} else {
//			// 1. 旁白
//			if len(matches) > 1 {
//				for _, s := range sentences[idx:] {
//					key, url, _ := cos.GenerateSceneCosPath()
//					sMatches := r.FindStringSubmatch(s)
//					if len(sMatches) <= 1 {
//						break
//					}
//					// [url, sMatches[1]]
//					chooseMap[key] = []string{url, sMatches[1]}
//				}
//				break
//			} else {
//				if role != "" && role != currentRole && role != conf.ME {
//					currentRole = role
//					newContentList = append(newContentList, parseChangeFigure(role, currentRole))
//				}
//				newContentList = append(newContentList, ":"+sentence)
//			}
//		}
//	}
//	// newContentList 每个元素以;结尾
//	//if len(chooseMap) > 0 {
//	//	newContentList = append(newContentList, parseChoose(chooseMap))
//	//} else {
//	//	newContentList = append(newContentList, conf.EndScene)
//	//}
//
//	newContent = strings.Join(newContentList, ";\n")
//	println(newContent)
//	return newContent, chooseMap, nil
//}

func parseChangeFigure(newRole conf.RoleName, oldRole conf.RoleName) (figure string) {
	figure = fmt.Sprintf("changeFigure:%s -left -enter=enter-from-left -next;", conf.FigureFileMap[newRole])
	return figure
}

func parseChat(role conf.RoleName, sentenceContent string) (chat string) {
	chat = fmt.Sprintf("%s:%s;", role, sentenceContent)
	return chat
}

//func parseChoose(chooseMap map[string][]string) (newContent string) {
//	values := make([]string, 0, len(chooseMap))
//	for _, v := range chooseMap {
//		values = append(values, v[1]+":"+v[0])
//	}
//	choose := fmt.Sprintf("choose:%s;", strings.Join(values, "|"))
//	return choose
//}

func parseChangeBg(imgUrl string) (content string) {
	content = fmt.Sprintf("changeBg:%s -next;", imgUrl)
	return content
}

//func parseChat(sentenceKey string, sentenceContent string, currRole conf.RoleName) (newContent []string, role conf.RoleName) {
//	// 如果角色换了，则需要插入changeFigure
//	roleName := conf.GetRoleName(sentenceKey)
//	//newContent = parseChangeFigure(roleName, currRole)
//	newContent = append(newContent, parseChangeFigure(roleName, currRole))
//	// 按照冒号分隔（中文或英文）
//	//chatList := strings.Split(content, "：:")
//	chatContent := fmt.Sprintf("%s:%s", roleName, sentenceContent)
//	newContent = append(newContent, chatContent)
//	return newContent, roleName
//}

// 生成并上传Bg图片
//func generateBg(bgDesc string, sceneId string) (string, string, error) {
//	// 生成图片
//	bgDesc = bgDesc + "图片中不要出现任何人物"
//	imgBase54, err := aiart.Generate(bgDesc)
//	if err != nil {
//		return "", "", err
//	}
//	// 上传COS
//	url, key, err := cos.UploadImage(imgBase54)
//	if err != nil {
//		return "", key, err
//	}
//	// 保存至MySQL
//	//url := cos.GetObjectUrl("image/489014096096657676.jpg")
//	bgImg := bg_img.NewBackgroundImg(key, url, sceneId, bgDesc)
//	err = bg_img.SaveBackgroundImg(bgImg)
//	if err != nil {
//		return "", key, err
//	}
//	return url, key, nil
//}

func GenerateNewScene(rawScene string, bgUrl string, rawChooses string) (string, map[int64]string) {
	sceneList := make([]string, 0)
	chooseStr, chooses := formatChoose(rawChooses)
	sceneList = append(sceneList, parseChangeBg(bgUrl))
	sceneList = append(sceneList, formatChat(rawScene)...)
	sceneList = append(sceneList, chooseStr)
	newScene := strings.Join(sceneList, "\n")
	fmt.Println(newScene)
	return newScene, chooses
}

func getSentenceRole(sentence string) (role conf.RoleName) {
	if !strings.Contains(sentence, ":") {
		return ""
	}
	// 按开头
	if strings.HasPrefix(sentence, string(conf.ZHOU)) {
		return conf.ZHOU
	} else if strings.HasPrefix(sentence, string(conf.ZHAO)) {
		return conf.ZHAO
	} else if strings.HasPrefix(sentence, string(conf.LIN)) {
		return conf.LIN
	} else if strings.HasPrefix(sentence, string(conf.CHEN)) {
		return conf.CHEN
	} else if strings.HasPrefix(sentence, string(conf.ME)) {
		return conf.ME
	}
	return conf.ME
}

func formatChat(rawScene string) []string {
	// 去掉空行
	rawScene = strings.ReplaceAll(rawScene, "\n\n", "\n")
	// 按行分割
	sentences := strings.Split(rawScene, "\n")
	//// 正则匹配：任意字符+数字+.或:+任意字符
	//r := regexp.MustCompile(`^.*\d+[.:](.*)$`)

	// 当前角色
	currentRole := conf.ME
	chatList := make([]string, 0)
	for _, sentence := range sentences {
		sentence = strings.ReplaceAll(sentence, " ", "")
		sentence = strings.ReplaceAll(sentence, "：", ":")
		chatContent := sentence
		role := getSentenceRole(sentence)
		if strings.Contains(sentence, ":") {
			chatContent = strings.Split(sentence, ":")[1]
			if role != "" && role != currentRole && role != conf.ME {
				currentRole = role
				chatList = append(chatList, parseChangeFigure(role, currentRole))
			}
		}
		chatList = append(chatList, parseChat(role, chatContent))
	}
	return chatList
}


func formatChoose(rawChooses string) (string, map[int64]string) {
	var chooseMap = make(map[int64]string)
	var chooseList = make([]string, 0)
	chooseContent := strings.ReplaceAll(rawChooses, "\n\n", "\n")
	chooseContent = strings.ReplaceAll(chooseContent, " ", "")
	chooseContent = strings.ReplaceAll(chooseContent, "：", ":")
	//match := regexp.MustCompile(`\d+\.(.*)`)
	// 任意字符或数字+.或:+任意字符
	match := regexp.MustCompile( `^.*[.:](.*)$`)
	sentences := strings.Split(chooseContent, "\n")
	for _, sentence := range sentences {
		choose := match.FindStringSubmatch(sentence)
		chooseId, err := flake.Generate()
		if err != nil {
			continue
		}
		str := choose[1]
		chooseContent = strings.ReplaceAll(chooseContent, ":", "")
		chooseStr := fmt.Sprintf("%s:%s", str, "/api/v1/game/scene?sceneId=" +strconv.FormatInt(chooseId, 10))
		chooseList = append(chooseList, chooseStr)
		chooseMap[chooseId] = choose[1]
	}
	chooses := fmt.Sprintf("choose:%s;", strings.Join(chooseList, "|"))
	return chooses, chooseMap
}


//func Test() {
//	// 读test.txt
//	f, err := os.Open("test.txt")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer f.Close()
//	// 读取文件内容
//	var content string
//	buf := make([]byte, 1024)
//	for {
//		n, _ := f.Read(buf)
//		if 0 == n {
//			break
//		}
//		content += string(buf[:n])
//	}
//	// 生成场景
//	newContent, chooses, _ := GenerateSceneContent(content, "491213694139760899")
//	fmt.Println(newContent)
//	fmt.Println(chooses)
//	// 保存至res.txt
//	f2, err := os.OpenFile("res.txt", os.O_CREATE|os.O_WRONLY, 0666)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer f2.Close()
//	f2.WriteString(newContent)
//}
