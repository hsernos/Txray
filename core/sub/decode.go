// core/sub/decode.go 负责订阅内容的解码与解析相关功能
package sub

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// base64Decode 用于解码 base64 编码的字符串
// str：待解码的字符串
// 返回值：解码后的字符串
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

// Sub2links 解析订阅文本
// subtext：待解析的订阅文本
// 返回值：解析出的链接列表
// 该函数首先调用 base64Decode 进行 base64 解码，
// 然后将解码后的文本中的 CRLF (\r\n) 和 CR (\r) 替换为 LF (\n)，
// 最后以 LF 为分隔符拆分字符串，返回拆分出的链接列表。
func Sub2links(subtext string) []string {
	data := base64Decode(subtext)
	s := strings.ReplaceAll(data, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	list := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return list
}
