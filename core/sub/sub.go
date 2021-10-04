package sub

import (
	"Txray/log"
	"crypto/md5"
	"fmt"
	"net/http"
)

type Subscirbe struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Using bool   `json:"using"`
}

// New方法
func NewSubscirbe(url, name string) *Subscirbe {
	return &Subscirbe{Name: name, Url: url, Using: true}
}

func (s *Subscirbe) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s.Url)))
}

func (s *Subscirbe) UpdataNode(opt UpdataOption) []string {
	var res *http.Response
	var err error
	switch opt.proxyMode() {
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
