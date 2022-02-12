package protocols

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type VMess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port int    `json:"port"`
	Id   string `json:"id"`
	Scy  string `json:"scy"`
	Aid  int    `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
}

func (v *VMess) GetProtocolMode() Mode {
	return ModeVMess
}

func (v *VMess) GetName() string {
	return v.Ps
}

func (v *VMess) GetAddr() string {
	return v.Add
}

func (v *VMess) GetPort() int {
	return v.Port
}

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
	buf.WriteString(fmt.Sprintf("%7s: %s", "协议", v.GetProtocolMode()))
	return buf.String()
}

func (v *VMess) GetLink() string {
	data := map[string]string{
		"v":    v.V,
		"ps":   v.Ps,
		"add":  v.Add,
		"port": strconv.Itoa(v.Port),
		"id":   v.Id,
		"aid":  strconv.Itoa(v.Aid),
		"scy": v.Scy,
		"net":  v.Net,
		"type": v.Type,
		"host": v.Host,
		"path": v.Path,
		"tls":  v.Tls,
		"sni": v.Sni,
		"alpn": v.Alpn,
	}
	jsonData, _ := json.Marshal(data)
	return "vmess://" + base64EncodeWithEq(string(jsonData))
}

func (v *VMess) Check() *VMess {
	if v.Add != "" && v.Port > 0 && v.Port <= 65535 && v.Ps != "" && v.Id != "" && v.Net != "" {
		return v
	}
	return nil
}
