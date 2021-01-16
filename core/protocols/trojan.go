package protocols

import (
	"Txray/core/protocols/mode"
	"Txray/tools"
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Trojan struct {
	Password string
	Address  string
	Port     string
	Remarks  string
}

func (t *Trojan) ParseLink(link string) bool {
	if strings.ToLower(link[:9]) == "trojan://" {
		link = link[9:]
	}
	expr := `(^[a-zA-Z0-9-]*?)@([a-zA-Z0-9-_\.]*?):([0-9]*)(.*)`
	r, _ := regexp.Compile(expr)
	result := r.FindStringSubmatch(link)
	if len(result) != 5 {
		return false
	}
	t.Password = result[1]
	t.Address = result[2]
	t.Port = result[3]
	other := result[4]
	index := strings.IndexByte(other, '#')
	if index >= 0 {
		name := other[index+1:]
		if name == "" {
			t.Remarks = t.Address + ":" + t.Port
		} else {
			t.Remarks, _ = url.QueryUnescape(name)
			t.Remarks = strings.Trim(t.Remarks, "\r")
		}
		other = other[:index]
	} else {
		t.Remarks = t.Address + ":" + t.Port
	}
	return true
}

func (t *Trojan) GetProtocolMode() string {
	return mode.Trojan
}

func (t *Trojan) GetName() string {
	return t.Remarks
}
func (t *Trojan) GetAddr() string {
	return t.Address
}

func (t *Trojan) GetPort() int {
	if tools.IsNetPort(t.Port) {
		return tools.StrToInt(t.Port)
	}
	return -1
}

func (t *Trojan) GetInfo() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "别名", t.Remarks))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "地址", t.Address))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "端口", t.Port))
	buf.WriteString(fmt.Sprintf("%3s: %s\n", "密码", t.Password))
	buf.WriteString(fmt.Sprintf("%3s: %s", "协议", t.GetProtocolMode()))
	return buf.String()
}

func (t *Trojan) GetLink() string {
	return fmt.Sprintf("trojan://%s@%s:%s#%s", t.Password, t.Address, t.Port, url.QueryEscape(t.Remarks))
}
