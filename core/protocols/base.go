package protocols

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type Protocol interface {
	// 解析链接
	ParseLink(link string) bool
	// 获取协议模式
	GetProtocolMode() string
	// 获取别名
	GetName() string
	// 获取远程地址
	GetAddr() string
	// 获取远程端口
	GetPort() int
	// 获取节点数据
	GetInfo() string
	// 获取节点分享链接
	GetLink() string
}

// 解析订阅文本
func Sub2links(subtext string) []string {
	data := base64Decode(subtext)
	s := strings.ReplaceAll(data, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	list := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return list
}

func base64Encode(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

func base64EncodeWithEq(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

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

// 解析链接
func ParseLink(link string) Protocol {
	url := strings.Split(link, "://")
	if len(url) == 2 && len(url[1]) != 0 {
		proto := url[0]
		switch strings.ToLower(proto) {
		case "vmess":
			v := new(VMess)
			if v.ParseLink(link) {
				return v
			} else {
				return nil
			}
		case "vless":
			// TODO
		case "ss":
			s := new(ShadowSocks)
			if s.ParseLink(link) {
				return s
			} else {
				return nil
			}
		case "ssr":
			s := new(ShadowSocksR)
			if s.ParseLink(link) {
				return s
			} else {
				return nil
			}
		case "trojan":
			t := new(Trojan)
			if t.ParseLink(link) {
				return t
			} else {
				return nil
			}
		}
	}
	return nil
}

func IsSupportLinkFormat(link string) bool {
	if strings.Contains(link, "://") {
		class := strings.ToLower(strings.Split(link, "://")[0])
		return class == "ss" || class == "vmess" || class == "trojan" || class == "ssr"
	}
	return false
}
