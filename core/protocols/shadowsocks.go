package protocols

import (
	"Txray/core/protocols/mode"
	"Txray/tools"
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ShadowSocks struct {
	Password string
	Address  string
	Port     string
	Remarks  string
	Method   string
}

func (s *ShadowSocks) ParseLink(link string) bool {
	if strings.ToLower(link[:5]) == "ss://" {
		link = link[5:]
	}
	expr := "(^[a-zA-Z0-9-=_]*)(#.*$)?"
	r, _ := regexp.Compile(expr)
	base := r.FindStringSubmatch(link)
	if len(base) != 3 {
		return false
	}
	data := base64Decode(base[1])
	if data == "" {
		return false
	}
	expr = `^([a-z0-9-]*):([a-zA-Z0-9]*)@([a-zA-Z0-9-_\.]*):([0-9]{1,5})$`
	r, _ = regexp.Compile(expr)
	result := r.FindStringSubmatch(data)
	if len(result) != 5 {
		return false
	}
	s.Method = result[1]
	s.Password = result[2]
	s.Address = result[3]
	s.Port = result[4]
	if !tools.IsNetPort(s.Port) {
		return false
	}
	if len(base[2]) > 1 {
		s.Remarks, _ = url.QueryUnescape(base[2][1:])
	} else {
		s.Remarks = s.Address + ":" + s.Port
	}
	return true
}

func (s *ShadowSocks) GetProtocolMode() string {
	return mode.ShadowSocks
}

func (s *ShadowSocks) GetName() string {
	return s.Remarks
}
func (s *ShadowSocks) GetAddr() string {
	return s.Address
}

func (s *ShadowSocks) GetPort() int {
	if tools.IsNetPort(s.Port) {
		return tools.StrToInt(s.Port)
	}
	return -1
}

func (s *ShadowSocks) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "别名", s.Remarks))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "地址", s.Address))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "端口", s.Port))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "密码", s.Password))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "加密", s.Method))
	buf.WriteString(fmt.Sprintf("%3s: %s", "协议", s.GetProtocolMode()))
	return buf.String()
}

func (s *ShadowSocks) GetLink() string {
	src := fmt.Sprintf("%s:%s@%s:%d", s.Method, s.Password, s.Address, s.GetPort())
	return fmt.Sprintf("ss://%s#%s", base64EncodeWithEq(src), url.QueryEscape(s.Remarks))
}
