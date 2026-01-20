// core/protocols/base64.go 负责 base64 编解码相关工具函数
package protocols

import (
	"encoding/base64"
	"strings"
)

// base64Encode 对给定字符串进行 base64 编码，不使用填充字符 '='
// 参数:
//   str: 需要编码的字符串
// 返回值:
//   编码后的字符串
func base64Encode(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

// base64EncodeWithEq 对给定字符串进行 base64 编码，使用标准填充字符 '='
// 参数:
//   str: 需要编码的字符串
// 返回值:
//   编码后的字符串
func base64EncodeWithEq(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64Decode 对给定的 base64 编码字符串进行解码
// 参数:
//   str: 需要解码的字符串
// 返回值:
//   解码后的字符串，如果解码失败则返回空字符串
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
	// 判断是使用 URL 安全的 base64 编码还是标准的 base64 编码
	if strings.Contains(str, "-") || strings.Contains(str, "_") {
		data, err = base64.URLEncoding.DecodeString(str)
	} else {
		data, err = base64.StdEncoding.DecodeString(str)
	}
	if err != nil {
		return ""
	}
	return string(data)
}
