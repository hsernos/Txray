// core/protocols/socks.go 负责 Socks 协议的定义与相关操作
package protocols

import (
	"bytes"
	"fmt"
	"net/url"
)

// Socks 结构体定义了 Socks 代理服务器的相关属性
type Socks struct {
	Address  string `json:"address"`  // 代理服务器地址
	Port     int    `json:"port"`     // 代理服务器端口
	Username string `json:"username"` // 连接代理的用户名
	Password string `json:"password"` // 连接代理的密码
	Remarks  string `json:"remarks"`  // 代理的备注信息
}

// GetProtocolMode 获取协议模式
// 返回值：当前结构体所对应的协议模式，这里是 Socks
func (s *Socks) GetProtocolMode() Mode {
	return ModeSocks
}

// GetName 获取别名
// 返回值：代理的备注信息
func (s *Socks) GetName() string {
	return s.Remarks
}

// GetAddr 获取远程地址
// 返回值：代理服务器地址
func (s *Socks) GetAddr() string {
	return s.Address
}

// GetPort 获取远程端口
// 返回值：代理服务器端口
func (s *Socks) GetPort() int {
	return s.Port
}

// GetInfo 获取节点数据
// 返回值：格式化的字符串，包含代理的详细信息
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
// 返回值：格式化的 URL 字符串，可用于分享节点信息
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

// Check 检查 Socks 结构体中的字段值是否有效
// 返回值：如果有效返回指向当前 Socks 结构体的指针，否则返回 nil
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
