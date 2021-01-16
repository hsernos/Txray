package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 通过http代理访问网站
func GetByHTTPProxy(objUrl, proxyAddress string, proxyPort, timeOut uint) (*http.Response, error) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("http://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeOut) * time.Second,
	}
	return client.Get(objUrl)
}

// 通过Socks5代理访问网站
func GetBySocks5Proxy(objUrl, proxyAddress string, proxyPort, timeOut uint) (*http.Response, error) {

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(fmt.Sprintf("socks5://%s:%d", proxyAddress, proxyPort))
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeOut) * time.Second,
	}
	return client.Get(objUrl)
}

// 不通过代理访问网站
func GetNoProxy(url string, timeOut uint) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
	}
	return client.Get(url)
}

// 读取http响应的内容
func ReadDate(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
