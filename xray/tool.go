package xray

import (
	"Txray/log"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 获取节点代理访问外网的延迟
func TestNode(url string, port int, timeout int) (int, string) {
	result, status := -1, "Error"
	for i := 0; i < 2; i++ {
		start := time.Now()
		res, e := GetBySocks5Proxy(url, "127.0.0.1", port, time.Duration(timeout)*time.Second)
		elapsed := time.Since(start)
		if e != nil {
			if strings.HasSuffix(fmt.Sprint(e), "(Client.Timeout exceeded while awaiting headers)") {
				log.Warnf("请求超时, 网络环境波动 或 配置test.timeout=%d 过小", timeout)
			} else if strings.HasSuffix(fmt.Sprint(e), "connect: connection refused"){
				log.Warn("代理服务未启动或启动中, 配置test.before过小")
			} else {
				log.Warn(e)
			}
		} else {
			result, status := int(float32(elapsed.Nanoseconds())/1e6), res.Status
			defer res.Body.Close()
			return result, status
		}
	}
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
