package protocols

import (
	"Txray/core/protocols/mode"
	"Txray/tools"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type ShadowSocksR struct {
	Address    string
	Port       string
	Protocol   string
	Method     string
	Obfs       string
	Password   string
	ObfsParam  string
	ProtoParam string
	Remarks    string
	Group      string
}

func (s *ShadowSocksR) ParseLink(link string) bool {
	if strings.ToLower(link[:6]) == "ssr://" {
		link = link[6:]
	}
	link = base64Decode(link)
	if link == "" {
		return false
	}
	expr := `^([a-zA-Z0-9-\.]*):([0-9]{1,5}):([a-z0-9_]*):([a-z0-9-]*):([a-z0-9_\.]*):([a-zA-Z0-9-_=]*)(/.(.*$))?`
	r, _ := regexp.Compile(expr)
	result := r.FindStringSubmatch(link)
	if len(result) != 9 {
		return false
	}
	s.Address = result[1]
	s.Port = result[2]
	if !tools.IsNetPort(s.Port) {
		return false
	}
	s.Protocol = result[3]
	s.Method = result[4]
	s.Obfs = result[5]
	s.Password = base64Decode(result[6])
	if s.Password == "" {
		return false
	}
	for _, str := range strings.Split(result[8], "&") {
		if strings.HasPrefix(str, "obfsparam=") {
			s.ObfsParam = base64Decode(str[10:])
		} else if strings.HasPrefix(str, "protoparam=") {
			s.ProtoParam = base64Decode(str[11:])
		} else if strings.HasPrefix(str, "remarks=") {
			s.Remarks = base64Decode(str[8:])
		} else if strings.HasPrefix(str, "group=") {
			s.Group = base64Decode(str[6:])
		}
	}
	if s.Remarks == "" {
		s.Remarks = s.Address + ":" + s.Port
	}
	return true
}

func (s *ShadowSocksR) GetProtocolMode() string {
	return mode.ShadowSocksR
}

func (s *ShadowSocksR) GetName() string {
	return s.Remarks
}
func (s *ShadowSocksR) GetAddr() string {
	return s.Address
}

func (s *ShadowSocksR) GetPort() int {
	if tools.IsNetPort(s.Port) {
		return tools.StrToInt(s.Port)
	}
	return -1
}

func (s *ShadowSocksR) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "别名", s.Remarks))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "地址", s.Address))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "端口", s.Port))
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
	result := fmt.Sprintf("%s:%s:%s:%s:%s:%s%s", s.Address, s.Port, s.Protocol, s.Method, s.Obfs,
		base64Encode(s.Password), params)
	return "ssr://" + base64EncodeWithEq(result)
}
