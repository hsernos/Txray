package protocols

import (
	"bytes"
	"fmt"
)

type ShadowSocksR struct {
	Address    string `json:"address"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	Method     string `json:"method"`
	Obfs       string `json:"obfs"`
	Password   string `json:"password"`
	ObfsParam  string `json:"obfsParam"`
	ProtoParam string `json:"protoParam"`
	Remarks    string `json:"remarks"`
	Group      string `json:"group"`
}

func (s *ShadowSocksR) GetProtocolMode() Mode {
	return ModeShadowSocksR
}

func (s *ShadowSocksR) GetName() string {
	return s.Remarks
}
func (s *ShadowSocksR) GetAddr() string {
	return s.Address
}

func (s *ShadowSocksR) GetPort() int {
	return s.Port
}

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

func (s *ShadowSocksR) GetLink() string {
	formatParams := "/?obfsparam=%s&protoparam=%s&remarks=%s&group=%s"
	params := fmt.Sprintf(formatParams, base64Encode(s.ObfsParam), base64Encode(s.ProtoParam),
		base64Encode(s.Remarks), base64Encode(s.Group))
	result := fmt.Sprintf("%s:%d:%s:%s:%s:%s%s", s.Address, s.Port, s.Protocol, s.Method, s.Obfs,
		base64Encode(s.Password), params)
	return "ssr://" + base64EncodeWithEq(result)
}
