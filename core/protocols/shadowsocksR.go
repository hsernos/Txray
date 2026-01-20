// core/protocols/shadowsocksR.go 负责 ShadowsocksR 协议的定义与相关操作
package protocols

import (
	"bytes"
	"fmt"
)

// ShadowSocksR 结构体定义了 ShadowsocksR 协议所需的各项参数
type ShadowSocksR struct {
	Address    string `json:"address"`    // 服务器地址
	Port       int    `json:"port"`       // 服务器端口
	Protocol   string `json:"protocol"`   // 使用的协议
	Method     string `json:"method"`     // 加密方法
	Obfs       string `json:"obfs"`       // 混淆方式
	Password   string `json:"password"`   // 连接密码
	ObfsParam  string `json:"obfsParam"`  // 混淆参数
	ProtoParam string `json:"protoParam"` // 协议参数
	Remarks    string `json:"remarks"`    // 备注
	Group      string `json:"group"`      // 分组
}

// GetProtocolMode 返回协议模式
func (s *ShadowSocksR) GetProtocolMode() Mode {
	return ModeShadowSocksR
}

// GetName 返回备注信息
func (s *ShadowSocksR) GetName() string {
	return s.Remarks
}

// GetAddr 返回服务器地址
func (s *ShadowSocksR) GetAddr() string {
	return s.Address
}

// GetPort 返回服务器端口
func (s *ShadowSocksR) GetPort() int {
	return s.Port
}

// GetInfo 返回该连接的详细信息，包含各个参数的说明
func (s *ShadowSocksR) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "别名", s.Remarks))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "地址", s.Address))
	buf.WriteString(fmt.Sprintf("%7s: %d\n", "端口", s.Port))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "加密方法", s.Method))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "密码", s.Password))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "协议", s.Protocol))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "协议参数", s.ProtoParam))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "混淆", s.Obfs))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "混淆参数", s.ObfsParam))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "分组", s.Group))
	buf.WriteString(fmt.Sprintf("%7s: %s", "协议", s.GetProtocolMode()))
	return buf.String()
}

// GetLink 生成该连接的链接字符串，供客户端使用
func (s *ShadowSocksR) GetLink() string {
	formatParams := "/?obfsparam=%s&protoparam=%s&remarks=%s&group=%s"
	params := fmt.Sprintf(formatParams, base64Encode(s.ObfsParam), base64Encode(s.ProtoParam),
		base64Encode(s.Remarks), base64Encode(s.Group))
	result := fmt.Sprintf("%s:%d:%s:%s:%s:%s%s", s.Address, s.Port, s.Protocol, s.Method, s.Obfs,
		base64Encode(s.Password), params)
	return "ssr://" + base64EncodeWithEq(result)
}
