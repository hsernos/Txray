package service

import (
	"Txray/core/node"
	"Txray/core/protocols"
	"Txray/core/setting"
	"Txray/core/xray"
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var exeCmd *exec.Cmd

func Start(key string, isTcpSort bool) {
	indexList := make([]int, 0)
	testUrl := setting.TestSetting().Url
	testTimeout := setting.TestSetting().TimeOut
	if isTcpSort {
		temp := node.GetTcpSortIndex(false)
		for _, index := range format.IndexDeal(key, len(temp)) {
			indexList = append(indexList, temp[index])
		}
	} else {
		indexList = format.IndexDeal(key, node.NodesSize())
	}
	if len(indexList) == 0 {
		log.Warn("没有选取到节点")
	} else if len(indexList) == 1 {
		node.SetIndex(indexList[0] + 1)
		exe := run(node.GetSelectedIndex())
		if exe {
			if setting.HttpPort() == 0 {
				log.Infof("启动成功, 监听socks端口: %d, 所选节点: %d", setting.SocksPort(), node.GetSelectedIndex())
			} else {
				log.Infof("启动成功, 监听socks/http端口: %d/%d, 所选节点: %d", setting.SocksPort(), setting.HttpPort(), node.GetSelectedIndex())
			}
			result, status := TestNode(testUrl, setting.SocksPort(), testTimeout)
			log.Infof("%6s [ %s ] 延迟: %dms", status, testUrl, result)

		}
	} else {
		min := 100000
		i := -1
		for _, index := range indexList {
			exe := run(index + 1)
			if exe {
				result, status := TestNode(testUrl, setting.SocksPort(), testTimeout)
				log.Infof("%6s [ %s ] 节点: %d, 延迟: %dms", status, testUrl, index+1, result)
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
			node.SetIndex(i + 1)
			exe := run(i + 1)
			if exe {
				if setting.HttpPort() == 0 {
					log.Infof("启动成功, 监听socks端口: %d, 所选节点: %d", setting.SocksPort(), node.GetSelectedIndex())
				} else {
					log.Infof("启动成功, 监听socks/http端口: %d/%d, 所选节点: %d", setting.SocksPort(), setting.HttpPort(), node.GetSelectedIndex())
				}
			} else {
				log.Error("启动失败")
			}
		} else {
			log.Info("所选节点全部不能访问外网")
		}

	}
}

func run(index int) bool {
	Stop()
	switch node.GetNode(index).GetProtocolMode() {
	case protocols.ModeShadowSocks, protocols.ModeTrojan, protocols.ModeVMess:
		if xray.CheckFile() {
			file := xray.GenConfig(index)
			exeCmd = exec.Command(xray.GetPath(), "-c", file)
		} else {
			return false
		}
	default:
		return false
	}
	stdout, _ := exeCmd.StdoutPipe()
	_ = exeCmd.Start()
	r := bufio.NewReader(stdout)
	lines := new([]string)
	go readInfo(r, lines)
	status := make(chan struct{})
	go checkProc(exeCmd, status)
	stopper := time.NewTimer(time.Millisecond * 300)
	select {
	case <-stopper.C:
		return true
	case <-status:
		log.Error("开启xray服务失败, 查看下面报错信息来检查出错问题")
		for _, x := range *lines {
			log.Error(x)
		}
		return false
	}
}

// 获取节点代理访问外网的延迟
func TestNode(url string, port uint, timeout uint) (int, string) {
	start := time.Now()
	res, e := tools.GetBySocks5Proxy(url, "127.0.0.1", port, timeout)
	elapsed := time.Since(start)
	if e != nil {
		log.Warn(e)
		return -1, "Error"
	}
	result, status := strings.Trim(fmt.Sprintf("%4.0f", float32(elapsed.Nanoseconds())/1e6), " "), res.Status
	defer res.Body.Close()
	return tools.StrToInt(result), status
}

// 停止服务
func Stop() {
	if exeCmd != nil {
		err := exeCmd.Process.Kill()
		if err != nil {
			log.Error(err)
		}
		exeCmd = nil
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
func checkProc(c *exec.Cmd, status chan struct{}) {
	c.Wait()
	status <- struct{}{}
}
