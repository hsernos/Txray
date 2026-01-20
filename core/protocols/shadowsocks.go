// core/protocols/shadowsocks.go 负责 Shadowsocks 协议的定义与相关操作
package protocols

import (
	"bytes"
	"fmt"
	"net/url"
)

// ShadowSocks 结构体定义了 Shadowsocks 协议所需的基本信息
type ShadowSocks struct {
	url.Values
	Password string `json:"password"` // 密码
	Address  string `json:"address"`  // 地址
	Port     int    `json:"port"`     // 端口
	Remarks  string `json:"remarks"`  // 别名
	Method   string `json:"method"`   // 加密方法
}

// GetProtocolMode 返回协议模式
func (s *ShadowSocks) GetProtocolMode() Mode {
	return ModeShadowSocks
}

// GetName 返回 Shadowsocks 的别名
func (s *ShadowSocks) GetName() string {
	return s.Remarks
}

// GetAddr 返回 Shadowsocks 的地址
func (s *ShadowSocks) GetAddr() string {
	return s.Address
}

// GetPort 返回 Shadowsocks 的端口
func (s *ShadowSocks) GetPort() int {
	return s.Port
}

// GetInfo 返回 Shadowsocks 的详细信息
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

// GetLink 返回 Shadowsocks 的链接信息
func (s *ShadowSocks) GetLink() string {
	if len(s.Values) == 0 {
		// 未使用插件
		src := fmt.Sprintf("%s:%s@%s:%d", s.Method, s.Password, s.Address, s.GetPort())
		return fmt.Sprintf("ss://%s#%s", base64EncodeWithEq(src), url.QueryEscape(s.Remarks))
	} else {
		// 使用了插件
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

// Check 检查 Shadowsocks 的基本信息是否完整
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
