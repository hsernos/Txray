package protocols

import (
	"bytes"
	"fmt"
	"net/url"
)

type ShadowSocks struct {
	url.Values
	Password string `json:"password"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Remarks  string `json:"remarks"`
	Method   string `json:"method"`
}

func (s *ShadowSocks) GetProtocolMode() Mode {
	return ModeShadowSocks
}

func (s *ShadowSocks) GetName() string {
	return s.Remarks
}
func (s *ShadowSocks) GetAddr() string {
	return s.Address
}

func (s *ShadowSocks) GetPort() int {
	return s.Port
}

func (s *ShadowSocks) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "别名", s.Remarks))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "地址", s.Address))
	buf.WriteString(fmt.Sprintf("%3s: %d\n", "端口", s.Port))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "密码", s.Password))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "加密", s.Method))
	if s.Has("plugin") {
		buf.WriteString(fmt.Sprintf("%3s: %s\n", "其他", "plugin="+s.Get("plugin")))
	}
	buf.WriteString(fmt.Sprintf("%3s: %s", "协议", s.GetProtocolMode()))
	return buf.String()
}

func (s *ShadowSocks) GetLink() string {
	if len(s.Values) == 0 {
		src := fmt.Sprintf("%s:%s@%s:%d", s.Method, s.Password, s.Address, s.GetPort())
		return fmt.Sprintf("ss://%s#%s", base64EncodeWithEq(src), url.QueryEscape(s.Remarks))
	} else {
		src := fmt.Sprintf("%s:%s", s.Method, s.Password)
		u := &url.URL{
			Scheme:   "ss",
			Host:     fmt.Sprintf("%s:%d", s.Address, s.Port),
			Fragment: s.Remarks,
			User:     url.User(base64EncodeWithEq(src)),
			RawQuery: s.Values.Encode(),
		}
		return u.String()
	}
}

func (s *ShadowSocks) Check() *ShadowSocks {
	if s.Address != "" && s.Port > 0 && s.Port < 65535 && s.Remarks != "" {
		if s.Method == "none" {
			return s
		} else if s.Password != "" && s.Method != "" {
			return s
		}
	}
	return nil
}
