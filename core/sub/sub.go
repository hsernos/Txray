// core/sub/sub.go 负责订阅对象的定义与相关操作
package sub

import (
	"Txray/log"
	"crypto/md5"
	"fmt"
	"net/http"
)

// Subscirbe 结构体表示一个订阅对象
type Subscirbe struct {
	Name  string `json:"name"`  // 订阅名称
	Url   string `json:"url"`   // 订阅链接
	Using bool   `json:"using"` // 是否启用该订阅
}

// New方法用于创建一个新的订阅对象
func NewSubscirbe(url, name string) *Subscirbe {
	return &Subscirbe{Name: name, Url: url, Using: true}
}

// ID 方法返回订阅的唯一标识符，使用 MD5 哈希订阅链接
func (s *Subscirbe) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s.Url)))
}

// UpdataNode 方法根据提供的选项更新节点信息
// opt: 更新选项，包括代理模式、地址、端口和超时时间
// 返回值: 返回更新后的节点链接列表
func (s *Subscirbe) UpdataNode(opt UpdataOption) []string {
	var res *http.Response
	var err error
	// 根据代理模式选择不同的获取方式
	switch opt.proxyMode() {
		case MIXED:
		res, err = GetByMixedProxy(s.Url, opt.addr(), opt.port(), opt.timeout())
	case SOCKS:
		res, err = GetBySocks5Proxy(s.Url, opt.addr(), opt.port(), opt.timeout())
	case HTTP:
		res, err = GetByHTTPProxy(s.Url, opt.addr(), opt.port(), opt.timeout())
	default:
		res, err = GetNoProxy(s.Url, opt.timeout())
	}
	if err != nil {
		log.Error(err)
		return []string{}
	}
	log.Info("访问 [", s.Url, "] -- ", res.Status)
	text := ReadDate(res)
	res.Body.Close()
	return Sub2links(text)
}
