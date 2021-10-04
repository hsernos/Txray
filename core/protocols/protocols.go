package protocols

import (
	"encoding/json"
	"regexp"
)

type Protocol interface {
	// GetProtocolMode 获取协议模式
	GetProtocolMode() Mode
	// GetName 获取别名
	GetName() string
	// GetAddr 获取远程地址
	GetAddr() string
	// GetPort 获取远程端口
	GetPort() int
	// GetInfo 获取节点数据
	GetInfo() string
	// GetLink 获取节点分享链接
	GetLink() string
}

// 序列化
func Serialize(p Protocol) string {
	jsonData, _ := json.Marshal(p)
	return string(p.GetProtocolMode()) + ": " + string(jsonData)
}

// 反序列化
func Deserialize(text string) Protocol {
	expr := "(^[a-zA-Z]*?): ({.*?}$)"
	r, _ := regexp.Compile(expr)
	result := r.FindStringSubmatch(text)
	if len(result) == 3 {
		jsonText := result[2]
		var data Protocol
		switch result[1] {
		case string(ModeVMess):
			data = new(VMess)
		case string(ModeShadowSocks):
			data = new(ShadowSocks)
		case string(ModeShadowSocksR):
			data = new(ShadowSocksR)
		case string(ModeTrojan):
			data = new(Trojan)
		case string(ModeSocks):
			data = new(Socks)
		case string(ModeVLESS):
			data = new(VLess)
		case string(ModeVMessAEAD):
			data = new(VMessAEAD)
		}
		err := json.Unmarshal([]byte(jsonText), &data)
		if err != nil {
			return nil
		}
		return data
	}
	return nil
}
