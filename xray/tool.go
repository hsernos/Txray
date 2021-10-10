package xray

import (
	"Txray/log"
	"fmt"
	"net/http"
	"net/url"
	"time"
)


// 获取节点代理访问外网的延迟
func TestNode(url string, port int, timeout int) (int, string) {
	start := time.Now()
	res, e := GetBySocks5Proxy(url, "127.0.0.1", port, time.Duration(timeout)*time.Second)
	elapsed := time.Since(start)
	if e != nil {
		log.Warn(e)
		return -1, "Error"
	}
	result, status := int(float32(elapsed.Nanoseconds())/1e6), res.Status
	defer res.Body.Close()
	return result, status
}

// 通过Socks5代理访问网站
func GetBySocks5Proxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("socks5://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeOut,
	}
	return client.Get(objUrl)
}
