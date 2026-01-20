// core/sub/net.go 负责订阅相关的网络请求与数据获取
package sub

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// 通过Mixed代理访问网站
// objUrl: 目标网址
// proxyAddress: 代理地址
// proxyPort: 代理端口
// timeOut: 超时时间
// 返回http响应和错误信息
func GetByMixedProxy(objUrl, proxyAddress string, proxyPort int, timeOut time.Duration) (*http.Response, error) {
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

// 通过http代理访问网站
// objUrl: 目标网址
// proxyAddress: 代理地址
// proxyPort: 代理端口
// timeOut: 超时时间
// 返回http响应和错误信息
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
// objUrl: 目标网址
// proxyAddress: 代理地址
// proxyPort: 代理端口
// timeOut: 超时时间
// 返回http响应和错误信息
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
// url: 目标网址
// timeOut: 超时时间
// 返回http响应和错误信息
func GetNoProxy(url string, timeOut time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeOut,
	}
	return client.Get(url)
}

// 读取http响应的内容
// resp: http响应
// 返回响应体内容
func ReadDate(resp *http.Response) string {
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
