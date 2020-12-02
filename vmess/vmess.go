package vmess

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// vmess链接信息
type Vmess struct {
	Add  string `json:"add"`
	Aid  int    `json:"aid,string"`
	Host string `json:"host"`
	ID   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Port uint   `json:"port,string"`
	PS   string `json:"ps"`
	TLS  string `json:"tls"`
	Type string `json:"type"`
	V    string `json:"v"`
}

// vmess对象转化成vmess链接
func (v *Vmess) ToLink() string {
	b, _ := json.Marshal(v)
	return "vmess://" + base64.StdEncoding.EncodeToString(b)
}

// 链接转化成vmess对象
func Link2vmessobj(link string) *Vmess {
	if len(link) > 8 && link[:8] == "vmess://" {
		decodeData, _ := base64.StdEncoding.DecodeString(link[8:])
		v := Vmess{}
		err := json.Unmarshal(decodeData, &v)
		if err != nil {
			return nil
		}
		return &v
	}
	return nil
}

// 解析订阅文本转化成vmess链接
func Sub2links(subtext string) []string {
	decodeData, _ := base64.URLEncoding.DecodeString(subtext)
	list := strings.Split(strings.TrimRight(string(decodeData), "\n"), "\n")
	return list
}

// 解析vmessList链接
func Links2vmessObjs(list []string) []*Vmess {
	result := make([]*Vmess, 0, len(list))
	for _, x := range list {
		v := Link2vmessobj(x)
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}
