package core

import (
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// 开启v2ray服务
func (c *Core) Start(key string, isTcpSort bool) {
	indexList := make([]int, 0)
	if isTcpSort {
		temp := c.GetTcpSortIndex(false)
		for _, index := range format.IndexDeal(key, len(temp)) {
			indexList = append(indexList, temp[index])
		}
	} else {
		indexList = format.IndexDeal(key, len(c.Nodes))
	}
	if len(indexList) == 0 {
		log.Warn("没有选取到节点")
	} else if len(indexList) == 1 {
		c.Index = uint(indexList[0])
		c.Save()
		c.exeCmd = c.run()
		if c.exeCmd != nil {
			if c.Settings.Http == 0 {
				log.Infof("启动成功, 监听socks端口: %d, 所选节点: %d", c.Settings.Socks, c.GetNodeIndex())
			} else {
				log.Infof("启动成功, 监听socks/http端口: %d/%d, 所选节点: %d", c.Settings.Socks, c.Settings.Http, c.GetNodeIndex())
			}
			result, status := c.TestNode("https://www.youtube.com")
			log.Infof("%6s [ https://www.youtube.com ] 延迟: %dms", status, result)

		}
	} else {
		min := 100000
		i := -1
		for _, index := range indexList {
			c.Index = uint(index)
			c.Save()
			c.exeCmd = c.run()
			if c.exeCmd != nil {
				result, status := c.TestNode("https://www.youtube.com")
				log.Infof("%6s [ https://www.youtube.com ] 节点: %d, 延迟: %dms", status, index+1, result)
				if result != -1 && min > result {
					i = index
					min = result
				}
			} else {
				return
			}
		}
		if i != -1 {
			log.Info("延迟最小的节点为：", i+1, "，延迟：", min, "ms")
			c.Index = uint(i)
			c.Save()
			c.exeCmd = c.run()
			if c.exeCmd != nil {
				if c.Settings.Http == 0 {
					log.Infof("启动成功, 监听socks端口: %d, 所选节点: %d", c.Settings.Socks, c.GetNodeIndex())
				} else {
					log.Infof("启动成功, 监听socks/http端口: %d/%d, 所选节点: %d", c.Settings.Socks, c.Settings.Http, c.GetNodeIndex())
				}
			} else {
				log.Error("启动失败")
			}
		} else {
			log.Info("所选节点全部不能访问外网")
		}

	}
}

// 获取节点代理访问外网的延迟
func (c *Core) TestNode(url string) (int, string) {
	start := time.Now()
	res, e := tools.GetBySocks5Proxy(url, "127.0.0.1", c.Settings.Socks, 10)
	elapsed := time.Since(start)
	if e != nil {
		log.Warn(e)
		return -1, "Error"
	}
	result, status := strings.Trim(fmt.Sprintf("%4.0f", float32(elapsed.Nanoseconds())/1e6), " "), res.Status
	defer res.Body.Close()
	return tools.StrToInt(result), status
}

func (c Core) run() *exec.Cmd {
	c.Stop()
	name, args := c.GenConfig()
	if args == nil {
		return nil
	}
	exe := exec.Command(name, args...)
	stdout, _ := exe.StdoutPipe()
	_ = exe.Start()
	r := bufio.NewReader(stdout)
	lines := new([]string)
	go readInfo(r, lines)
	code := new(int)
	go checkProc(exe, code)
	time.Sleep(time.Duration(500) * time.Millisecond)
	if *code > 0 {
		log.Error("开启v2ray服务失败, 查看下面报错信息来检查出错问题")
		for _, x := range *lines {
			log.Error(x)
		}
		return nil
	}
	return exe
}

// 停止v2ray服务
func (c *Core) Stop() {
	if c.exeCmd != nil {
		err := c.exeCmd.Process.Kill()
		if err != nil {
			log.Error(err)
		}
		c.exeCmd = nil
	}
}

func readInfo(r *bufio.Reader, lines *[]string) {
	for i := 0; i < 20; i++ {
		line, _, _ := r.ReadLine()
		if len(string(line[:])) != 0 {
			*lines = append(*lines, string(line[:]))
		}
	}
}

// 检查进程状态
func checkProc(c *exec.Cmd, code *int) {
	c.Wait()
	*code = c.ProcessState.ExitCode()
}
