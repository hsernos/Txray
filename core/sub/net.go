package sub

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// 通过http代理访问网站
func GetByHTTPProxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("http://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeOut,
	}
	return client.Get(objUrl)
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

// 不通过代理访问网站
func GetNoProxy(url string, timeOut time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeOut,
	}
	return client.Get(url)
}

// 读取http响应的内容
func ReadDate(resp *http.Response) string {
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
