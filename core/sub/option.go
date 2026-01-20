// core/sub/option.go 负责订阅选项的定义与相关操作
package sub

import (
	"Txray/core/setting"
	"time"
)

// ProxyProtocol 定义代理协议类型
type ProxyProtocol int

const (
	NONE ProxyProtocol = iota // 无代理
	MIXED                     // 混合代理
	SOCKS                     // SOCKS代理
	HTTP                      // HTTP代理
)

// UpdataOption 结构体用于保存订阅的相关选项
type UpdataOption struct {
	Key       string        // 订阅的唯一标识
	ProxyMode ProxyProtocol // 代理模式
	Addr      string        // 代理地址
	Port      int           // 代理端口
	Timeout   time.Duration // 超时时间
}

// proxyMode 返回 UpdataOption 中设置的代理模式
func (o *UpdataOption) proxyMode() ProxyProtocol {
	return o.ProxyMode
}

// addr 返回 UpdataOption 中设置的代理地址，默认返回 "127.0.0.1"
func (o *UpdataOption) addr() string {
	if o.Addr == "" {
		return "127.0.0.1"
	}
	return o.Addr
}

// port 返回 UpdataOption 中设置的代理端口，
// 如果未设置端口，则根据代理模式返回默认端口
func (o *UpdataOption) port() int {
	if o.Port != 0 {
		return o.Port
	}
	switch o.ProxyMode {
	case MIXED:
		return setting.Mixed()
	case SOCKS:
		return setting.Socks()
	case HTTP:
		return setting.Http()
	}
	return o.Port
}

// timeout 返回 UpdataOption 中设置的超时时间，默认返回 5 秒
func (o *UpdataOption) timeout() time.Duration {
	if o.Timeout == 0 {
		return 5 * time.Second
	}
	return o.Timeout
}
