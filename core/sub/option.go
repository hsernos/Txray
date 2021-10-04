package sub

import (
	"Txray/core/setting"
	"time"
)

type ProxyProtocol int

const (
	NONE ProxyProtocol = iota
	SOCKS
	HTTP
)

type UpdataOption struct {
	Key       string
	ProxyMode ProxyProtocol
	Addr      string
	Port      int
	Timeout   time.Duration
}

func (o *UpdataOption) proxyMode() ProxyProtocol {
	return o.ProxyMode
}

func (o *UpdataOption) addr() string {
	if o.Addr == "" {
		return "127.0.0.1"
	}
	return o.Addr
}

func (o *UpdataOption) port() int {
	if o.Port != 0 {
		return o.Port
	}
	switch o.ProxyMode {
	case SOCKS:
		return setting.Socks()
	case HTTP:
		return setting.Http()
	}
	return o.Port
}
func (o *UpdataOption) timeout() time.Duration {
	if o.Timeout == 0 {
		return 5 * time.Second
	}
	return o.Timeout
}
