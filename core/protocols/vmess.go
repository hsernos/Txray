package protocols

import (
	"Txray/core/protocols/mode"
	"Txray/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type VMess struct {
	V    string
	Ps   string
	Add  string
	Port string
	Id   string
	Aid  string
	Net  string
	Type string
	Host string
	Path string
	Tls  string
}

func (v *VMess) ParseLink(link string) bool {
	if strings.ToLower(link[:8]) == "vmess://" {
		link = link[8:]
	}
	if len(link) == 0 {
		return false
	}
	jsonStr := base64Decode(link)
	if jsonStr == "" {
		return false
	}
	var mapResult map[string]string
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		return false
	}
	var ok bool
	if v.V, ok = mapResult["v"]; !ok {
		return false
	}
	if v.Ps, ok = mapResult["ps"]; !ok {
		return false
	}
	if v.Add, ok = mapResult["add"]; !ok {
		return false
	}
	if v.Port, ok = mapResult["port"]; !ok {
		return false
	}
	if !tools.IsNetPort(v.Port) {
		return false
	}
	if v.Id, ok = mapResult["id"]; !ok {
		return false
	}
	if v.Aid, ok = mapResult["aid"]; !ok {
		return false
	}
	if v.Net, ok = mapResult["net"]; !ok {
		return false
	}
	if v.Type, ok = mapResult["type"]; !ok {
		return false
	}
	if v.Host, ok = mapResult["host"]; !ok {
		return false
	}
	if v.Path, ok = mapResult["path"]; !ok {
		return false
	}
	if v.Tls, ok = mapResult["tls"]; !ok {
		return false
	}
	return true
}

func (v *VMess) GetProtocolMode() string {
	return mode.VMess
}

func (v *VMess) GetName() string {
	return v.Ps
}

func (v *VMess) GetAddr() string {
	return v.Add
}

func (v *VMess) GetPort() int {
	if tools.IsNetPort(v.Port) {
		return tools.StrToInt(v.Port)
	}
	return -1
}

func (v *VMess) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "别名", v.Ps))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "地址", v.Add))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "端口", v.Port))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "用户ID", v.Id))
	buf.WriteString(fmt.Sprintf("%7s: %s\n", "额外ID", v.Aid))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "加密方式", "auto"))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装类型", v.Type))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "伪装域名", v.Host))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "传输协议", v.Net))
	buf.WriteString(fmt.Sprintf("%9s: %s\n", "path", v.Path))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "传输安全", v.Tls))
	buf.WriteString(fmt.Sprintf("%5s: %s\n", "配置版本", v.V))
	buf.WriteString(fmt.Sprintf("%7s: %s", "协议", v.GetProtocolMode()))
	return buf.String()
}

func (v *VMess) GetLink() string {
	data := map[string]string{
		"v":    v.V,
		"ps":   v.Ps,
		"add":  v.Add,
		"port": v.Port,
		"id":   v.Id,
		"aid":  v.Aid,
		"net":  v.Net,
		"type": v.Type,
		"host": v.Host,
		"path": v.Path,
		"tls":  v.Tls,
	}
	jsonData, _ := json.Marshal(data)
	return "vmess://" + base64Encode(string(jsonData))
}
