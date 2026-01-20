// core/protocols/vmess.go 负责 VMess 协议的定义与相关操作
package protocols

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// VMess 结构体定义了 VMess 协议所需的各项属性
type VMess struct {
	V             string `json:"v"`
	Ps            string `json:"ps"`
	Add           string `json:"add"`
	Port          int    `json:"port"`
	Id            string `json:"id"`
	Scy           string `json:"scy"`
	Aid           int    `json:"aid"`
	Net           string `json:"net"`
	Type          string `json:"type"`
	Host          string `json:"host"`
	Path          string `json:"path"`
	Tls           string `json:"tls"`
	Sni           string `json:"sni"`
	Alpn          string `json:"alpn"`
	EchConfigList string `json:"echConfigList"`
	EchForceQuery string `json:"echForceQuery"`
	PCS           string `json:"pinnedPeerCertSha256"`
}

// GetProtocolMode 返回协议模式
func (v *VMess) GetProtocolMode() Mode {
	return ModeVMess
}

// GetName 返回用户自定义的别名
func (v *VMess) GetName() string {
	return v.Ps
}

// GetAddr 返回服务器地址
func (v *VMess) GetAddr() string {
	return v.Add
}

// GetPort 返回端口号
func (v *VMess) GetPort() int {
	return v.Port
}

// GetInfo 返回 VMess 配置的详细信息，包含别名、地址、端口、用户ID、额外ID、加密方式等
func (v *VMess) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "别名", v.Ps))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "地址", v.Add))
	buf.WriteString(fmt.Sprintf("%7s: %d\n", "端口", v.Port))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "用户ID", v.Id))
	buf.WriteString(fmt.Sprintf("%7s: %d\n", "额外ID", v.Aid))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "加密方式", v.Scy))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装类型", v.Type))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装域名", v.Host))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "传输协议", v.Net))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "path", v.Path))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "传输安全", v.Tls))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "配置版本", v.V))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "SNI", v.Sni))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "Alpn", v.Alpn))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "Ech配置列表", v.EchConfigList))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "ECH强制查询", v.EchForceQuery))
	buf.WriteString(fmt.Sprintf("%16s: %s\n", "PinnedPeerCertSha256", v.PCS))
	buf.WriteString(fmt.Sprintf("%7s: %s", "协议", v.GetProtocolMode()))
	return buf.String()
}

// GetLink 生成 VMess 链接
func (v *VMess) GetLink() string {
	data := map[string]string{
		"v":             v.V,
		"ps":            v.Ps,
		"add":           v.Add,
		"port":          strconv.Itoa(v.Port),
		"id":            v.Id,
		"aid":           strconv.Itoa(v.Aid),
		"scy":           v.Scy,
		"net":           v.Net,
		"type":          v.Type,
		"host":          v.Host,
		"path":          v.Path,
		"tls":           v.Tls,
		"sni":           v.Sni,
		"alpn":          v.Alpn,
		"echConfigList": v.EchConfigList,
		"echForceQuery": v.EchForceQuery,
		"pinnedPeerCertSha256": v.PCS,
	}
	jsonData, _ := json.Marshal(data)
	return "vmess://" + base64EncodeWithEq(string(jsonData))
}

// GetValue 根据字段获取对应的值（VMess 不支持新传输协议的字段，返回空字符串）
func (v *VMess) GetValue(field interface{}) string {
	return ""
}

// GetExtraValue 根据ExtraField获取对应的值（VMess 不支持新传输协议的字段，返回空对象）
func (v *VMess) GetExtraValue(field interface{}) interface{} {
	return map[string]interface{}{}
}

// Check 检查 VMess 配置是否有效
func (v *VMess) Check() *VMess {
	if v.Add != "" && v.Port > 0 && v.Port <= 65535 && v.Ps != "" && v.Id != "" && v.Net != "" {
		return v
	}
	return nil
}
