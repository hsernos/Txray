package protocols

import (
	"Txray/core/protocols/field"
	"bytes"
	"fmt"
	"net/url"
)

type Hysteria2 struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Remarks  string `json:"remarks"`
	url.Values
}

// GetProtocolMode 获取协议模式
func (v *Hysteria2) GetProtocolMode() Mode {
	return ModeHysteria2
}

// GetName 获取别名
func (v *Hysteria2) GetName() string {
	return v.Remarks
}

// GetAddr 获取远程地址
func (v *Hysteria2) GetAddr() string {
	return v.Address
}

// GetPort 获取远程端口
func (v *Hysteria2) GetPort() int {
	return v.Port
}

// GetNetwork 获取远程传输方式
func (v *Hysteria2) GetNetwork() string {
	return "hysteria"
}

func (t *Hysteria2) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%10s: [%s]\n", "协议", t.GetProtocolMode()))
	buf.WriteString(fmt.Sprintf("%10s: %s\n", "别名", t.Remarks))
	buf.WriteString(fmt.Sprintf("%10s: %s\n", "地址", t.Address))
	buf.WriteString(fmt.Sprintf("%10s: %d\n", "端口", t.Port))
	buf.WriteString(fmt.Sprintf("%10s: %s\n", "密码", t.Password))
	buf.WriteString(fmt.Sprintf("%6s: %s\n", "跳跃端口范围", t.GetValue(field.Mport)))
	buf.WriteString(fmt.Sprintf("%6s: %s\n", "端口跳跃间隔", t.GetValue(field.Hopinterval)))
	buf.WriteString(fmt.Sprintf("%12s: %s", "FinalMask", t.GetValue(field.Finalmask)))
	return buf.String()
}

// GetLink 获取节点分享链接
func (v *Hysteria2) GetLink() string {
	u := url.URL{
		Scheme:   "hysteria2",
		User:     url.User(v.Password),
		Host:     fmt.Sprintf("%s:%d", v.GetAddr(), v.GetPort()),
		RawQuery: v.Values.Encode(),
		Fragment: v.Remarks,
	}
	return u.String()
}

func (v *Hysteria2) GetValue(field field.Field) string {
	if v.Has(field.Key) {
		return v.Get(field.Key)
	}
	return field.Value
}

func (v *Hysteria2) GetHostValue(field field.Field) string {
	if v.Has(field.Key) {
		return v.Get(field.Key)
	}
	return v.Address
}


func (t *Hysteria2) Check() *Hysteria2 {
	if t.Password != "" && t.Address != "" && t.Port > 0 && t.Port < 65535 && t.Remarks != "" {
		return t
	}
	return nil
}