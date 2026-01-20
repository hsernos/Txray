// core/protocols/vmessAEAD.go 负责 VMess AEAD 协议的定义与相关操作
package protocols

import (
	"Txray/core/protocols/field"
	"bytes"
	"fmt"
	"net/url"
)

// VMessAEAD 结构体定义了 VMess AEAD 协议所需的基本信息
type VMessAEAD struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	Remarks string `json:"remarks"`
	url.Values
}

// GetProtocolMode 获取协议模式
func (v *VMessAEAD) GetProtocolMode() Mode {
	return ModeVMessAEAD
}

// GetName 获取别名
func (v *VMessAEAD) GetName() string {
	return v.Remarks
}

// GetAddr 获取远程地址
func (v *VMessAEAD) GetAddr() string {
	return v.Address
}

// GetPort 获取远程端口
func (v *VMessAEAD) GetPort() int {
	return v.Port
}

// GetInfo 获取节点数据
func (v *VMessAEAD) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "别名", v.Remarks))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "地址", v.Address))
	buf.WriteString(fmt.Sprintf("%7s: %d\n", "端口", v.Port))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "用户ID", v.ID))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "加密方式", v.GetValue(field.VMessEncryption)))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "传输协议", v.GetValue(field.NetworkType)))
	switch v.GetValue(field.NetworkType) {
	case "tcp":
		buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装类型", v.GetValue(field.TCPHeaderType)))
	case "kcp":
		buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装类型", v.GetValue(field.MkcpHeaderType)))
		if v.GetValue(field.Seed) != "" {
			buf.WriteString(fmt.Sprintf("%7s: %s\n", "KCP种子", v.GetValue(field.Seed)))
		}
	case "ws":
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Path", v.GetValue(field.WsPath)))
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Host", v.GetValue(field.WsHost)))
	case "h2":
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Path", v.GetValue(field.WsPath)))
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Host", v.GetHostValue(field.WsHost)))
	case "quic":
		buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装类型", v.GetValue(field.QuicHeaderType)))
		buf.WriteString(fmt.Sprintf("%7s: %s\n", "QUIC加密", v.GetValue(field.QuicSecurity)))
		if v.GetValue(field.QuicSecurity) != "none" {
			buf.WriteString(fmt.Sprintf("%5s: %s\n", "加密密钥", v.GetValue(field.QuicKey)))
		}
	case "grpc":
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "ServiceName", v.GetValue(field.GrpcServiceName)))
		buf.WriteString(fmt.Sprintf("%5s: %s\n", "传输模式", v.GetValue(field.GrpcMode)))
	case "splithttp":
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Path", v.GetValue(field.SpPath)))
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Host", v.GetHostValue(field.SpHost)))
		buf.WriteString(fmt.Sprintf("%7s: %s\n", "Mode", v.GetValue(field.SpMode)))
		buf.WriteString(fmt.Sprintf("%7s: %v\n", "extra", v.GetExtraValue(field.SpExtra)))
	case "xhttp":
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Path", v.GetValue(field.XhPath)))
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "Host", v.GetHostValue(field.XhHost)))
		buf.WriteString(fmt.Sprintf("%7s: %s\n", "Mode", v.GetValue(field.XhMode)))
		buf.WriteString(fmt.Sprintf("%7s: %v\n", "extra", v.GetExtraValue(field.XhExtra)))
	}
	if v.GetValue(field.Security) == "reality" {
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "FingerPrint", v.GetValue(field.TLSFingerPrint)))
		buf.WriteString(fmt.Sprintf("%7s: %s\n", "reality公钥", v.GetValue(field.RealityPublicKey)))
		buf.WriteString(fmt.Sprintf("%9s: %s\n", "realityID", v.GetValue(field.RealityShortId)))
		buf.WriteString(fmt.Sprintf("%7s: %s\n", "REALITY爬虫", v.GetValue(field.RealitySpiderX)))
		buf.WriteString(fmt.Sprintf("%7s: %s\n", "MLDSA65验证", v.GetValue(field.RealityMldsa65Verify)))
	}
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "底层传输", v.GetValue(field.Security)))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "SNI", v.GetValue(field.SNI)))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "Alpn", v.GetValue(field.Alpn)))
	buf.WriteString(fmt.Sprintf("%7s: %s", "协议", v.GetProtocolMode()))

	return buf.String()
}

// GetLink 获取节点分享链接
func (v *VMessAEAD) GetLink() string {
	u := url.URL{
		Scheme:   "vmess",
		User:     url.User(v.ID),
		Host:     fmt.Sprintf("%s:%d", v.GetAddr(), v.GetPort()),
		RawQuery: v.Values.Encode(),
		Fragment: v.Remarks,
	}
	return u.String()
}

// GetValue 根据字段获取对应的值
func (v *VMessAEAD) GetValue(field field.Field) string {
	if v.Has(field.Key) {
		return v.Get(field.Key)
	}
	return field.Value
}

// GetExtraValue 根据ExtraField获取对应的值，当值为空时返回空对象
func (v *VMessAEAD) GetExtraValue(field field.ExtraField) interface{} {
	if v.Has(field.Key) {
		value := v.Get(field.Key)
		if value == "" {
			return map[string]interface{}{}
		}
		return value
	}
	return map[string]interface{}{}
}

// H2Host SNI
func (v *VMessAEAD) GetHostValue(field field.Field) string {
	if v.Has(field.Key) {
		return v.Get(field.Key)
	}
	return v.Address
}

// Check 检查 VMessAEAD 结构体中的必需字段是否有效
func (v *VMessAEAD) Check() *VMessAEAD {
	if v.ID != "" && v.Port > 0 && v.Port <= 65535 && v.Address != "" && v.Remarks != "" {
		return v
	}
	return nil
}
