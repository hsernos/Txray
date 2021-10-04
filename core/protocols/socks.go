package protocols

import (
	"bytes"
	"fmt"
	"net/url"
)

type Socks struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Remarks  string `json:"remarks"`
}

// GetProtocolMode 获取协议模式
func (s *Socks) GetProtocolMode() Mode {
	return ModeSocks
}

// GetName 获取别名
func (s *Socks) GetName() string {
	return s.Remarks
}

// GetAddr 获取远程地址
func (s *Socks) GetAddr() string {
	return s.Address
}

// GetPort 获取远程端口
func (s *Socks) GetPort() int {
	return s.Port
}

// GetInfo 获取节点数据
func (s *Socks) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "别名", s.Remarks))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "地址", s.Address))
	buf.WriteString(fmt.Sprintf("%7s: %d\n", "端口", s.Port))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "用户", s.Username))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "密码", s.Password))
	buf.WriteString(fmt.Sprintf("%7s: %s", "协议", s.GetProtocolMode()))
	return buf.String()
}

// GetLink 获取节点分享链接
func (s *Socks) GetLink() string {
	u := url.URL{
		Scheme:   "socks",
		Host:     fmt.Sprintf("%s:%d", s.Address, s.Port),
		Fragment: s.Remarks,
	}
	if s.Username != "" && s.Password != "" {
		u.User = url.UserPassword(s.Username, s.Password)
	}
	return u.String()
}

func (s *Socks) Check() *Socks {
	if s.Address != "" && s.Port > 0 && s.Port < 65535 && s.Remarks != "" {
		if s.Username == "" && s.Password == "" {
			return s
		}
		if s.Username != "" && s.Password != "" {
			return s
		}
	}
	return nil
}
