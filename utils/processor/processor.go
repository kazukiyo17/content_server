package processor

import (
	"regexp"
	"strings"
)

const (
	CHARACTERS = ['1']
)

func Process(text string) (resp string, err error) {
	// 按\n分隔
	sentences := strings.Split(text, "\n")
	for _, sentence := range sentences {
		// 查找第一个冒号（中英文都可以）
		index := strings.IndexAny(sentence, ":：")
		if index != -1 {
			character := sentence[:index]
			content := sentence[index+1:]
			// 如果character中存在
		}else{

		}
	}
}
