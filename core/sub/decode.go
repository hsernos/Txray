package sub

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func base64Decode(str string) string {
	i := len(str) % 4
	switch i {
	case 1:
		str = str[:len(str)-1]
	case 2:
		str += "=="
	case 3:
		str += "="
	}
	var data []byte
	var err error
	if strings.Contains(str, "-") || strings.Contains(str, "_") {
		data, err = base64.URLEncoding.DecodeString(str)
	} else {
		data, err = base64.StdEncoding.DecodeString(str)
	}
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

// 解析订阅文本
func Sub2links(subtext string) []string {
	data := base64Decode(subtext)
	s := strings.ReplaceAll(data, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	list := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return list
}
