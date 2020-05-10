package tool

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	log "v3ray/logger"
)

// Get get访问url
func Get(url string) string {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("url访问错误 ", err.Error())
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}

// GetByProxy 本地socks5代理访问url
func GetByProxy(url string, port uint) string {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:"+UintToStr(port), nil, &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	},
	)
	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{Transport: transport}
	res, err := client.Get(url)
	if err != nil {
		log.Error("url访问错误 ", err.Error())
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
