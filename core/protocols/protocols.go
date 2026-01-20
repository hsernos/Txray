// core/protocols/protocols.go 负责定义所有协议的接口和通用结构体
package protocols

import (
	"encoding/json"
	"regexp"
)

// Protocol 是所有协议的接口，定义了协议的基本操作
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
// Serialize 将协议对象序列化为字符串
func Serialize(p Protocol) string {
	jsonData, _ := json.Marshal(p)
	return string(p.GetProtocolMode()) + ": " + string(jsonData)
}

// 反序列化
// Deserialize 将字符串反序列化为协议对象
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
