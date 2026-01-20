// core/protocols/trojan.go 负责 Trojan 协议的定义与相关操作
package protocols

import (
	"bytes"
	"fmt"
	"net/url"
)

// Trojan 结构体表示一个 Trojan 协议的实例
type Trojan struct {
	url.Values
	Password string `json:"password"` // 密码
	Address  string `json:"address"`  // 地址
	Port     int    `json:"port"`     // 端口
	Remarks  string `json:"remarks"`  // 备注
}

// GetProtocolMode 返回协议模式
func (t *Trojan) GetProtocolMode() Mode {
	return ModeTrojan
}

// GetName 返回备注
func (t *Trojan) GetName() string {
	return t.Remarks
}

// GetAddr 返回地址
func (t *Trojan) GetAddr() string {
	return t.Address
}

// GetPort 返回端口
func (t *Trojan) GetPort() int {
	return t.Port
}

// GetInfo 返回 Trojan 实例的详细信息
func (t *Trojan) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "别名", t.Remarks))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "地址", t.Address))
	buf.WriteString(fmt.Sprintf("%3s: %d\n", "端口", t.Port))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "密码", t.Password))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "SNI", t.Sni()))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "ECH配置列表", t.EchConfigList()))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "ECH强制查询", t.EchForceQuery()))
	buf.WriteString(fmt.Sprintf("%3s: %s", "协议", t.GetProtocolMode()))
	return buf.String()
}

// GetLink 返回 Trojan 链接
func (t *Trojan) GetLink() string {
	u := &url.URL{
		Scheme:   "trojan",
		Host:     fmt.Sprintf("%s:%d", t.Address, t.Port),
		Fragment: t.Remarks,
		User:     url.User(t.Password),
		RawQuery: t.Values.Encode(),
	}
	return u.String()
}

// Sni 返回 SNI 信息
func (t *Trojan) Sni() string {
	if t.Has("sni") {
		return t.Get("sni")
	}
	return ""
}

// EchConfigList 返回 ECH 配置列表
func (t *Trojan) EchConfigList() string {
	if t.Has("echConfigList") {
		return t.Get("echConfigList")
	}
	return ""
}

// EchForceQuery 返回 ECH 强制查询信息
func (t *Trojan) EchForceQuery() string {
	if t.Has("echForceQuery") {
		return t.Get("echForceQuery")
	}
	return ""
}

// Check 检查 Trojan 实例的各项属性是否有效
func (t *Trojan) Check() *Trojan {
	if t.Password != "" && t.Address != "" && t.Port > 0 && t.Port < 65535 && t.Remarks != "" {
		return t
	}
	return nil
}
