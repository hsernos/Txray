package tool

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Vmess vmess链接信息
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


// GetVmesslink 获取vmess链接
func (v *Vmess) GetVmesslink() string {
	b, _ := json.Marshal(v)
	return "vmess://" + base64.StdEncoding.EncodeToString(b)
}

// SubToVmessList 解析订阅转化成vmess链接
func SubToVmessList(url string, port uint) []string {
	data := ""
	if port > 65535{
		data = Get(url)
	} else {
		data = GetByProxy(url, port)
	}
	decodeData, _ := base64.URLEncoding.DecodeString(data)
	list := strings.Split(string(decodeData), "\n")
	return list
}

// VmesslinkToObj 解析vmess链接
func VmesslinkToObj(vmesslink string) (*Vmess, error) {
	if len(vmesslink) > 8 && vmesslink[:8] == "vmess://" {
		decodeData, _ := base64.StdEncoding.DecodeString(vmesslink[8:])
		v := Vmess{}
		err := json.Unmarshal(decodeData, &v)
		if err != nil {
			return nil, err
		}
		return &v, err
	}
	return nil, nil
}

// VmessListToObj 解析vmessList链接
func VmessListToObj(list []string) []*Vmess {
	result := make([]*Vmess, 0, len(list))
	for _, x := range list {
		v, _ := VmesslinkToObj(x)
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}
