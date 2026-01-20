// xray/tool.go 负责 xray 相关的工具函数实现，如文件检测、路径查找等
package xray

import (
	"Txray/log"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// 获取节点代理访问外网的延迟
// 参数:
//   - url: 目标网址
//   - port: 代理端口
//   - timeout: 超时时间（秒）
// 返回值:
//   - int: 延迟时间（毫秒），出错时返回-1
//   - string: HTTP响应状态，出错时返回"Error"
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
// 参数:
//   - objUrl: 目标网址
//   - proxyAddress: 代理地址
//   - proxyPort: 代理端口
//   - timeOut: 超时时间
// 返回值:
//   - *http.Response: HTTP响应
//   - error: 错误信息
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
