package xray

import (
	"Txray/log"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func KillProcessByName(name string) error {
	name = strings.ToLower(name)
	processes, _ := process.Processes()
	for _, p := range processes {
		pName, _ := p.Name()
		pName = strings.TrimSuffix(strings.ToLower(pName), ".exe")
		if pName == name {
			err := p.Kill()
			if err != nil {
				process, err := os.FindProcess(int(p.Pid))
				if err != nil {
					return err
				}
				return process.Kill()

			}
		}
	}
	return nil
}

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
