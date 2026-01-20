// xray/service.go 负责 xray 核心服务的启动、停止、状态管理等功能
package xray

import (
	"Txray/core"
	"Txray/core/manage"
	"Txray/core/protocols"
	"Txray/core/setting"
	"Txray/log"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/hpcloud/tail"
)

var Xray *exec.Cmd

// Start 根据用户选择的节点信息，启动相应的 xray 服务
// key : 用户的节点选择的 key
func Start(key string) {
	testUrl := setting.TestUrl()
	testTimeout := setting.TestTimeout()
	manager := manage.Manager
	indexList := core.IndexList(key, manager.NodeLen())
	if len(indexList) == 0 {
		log.Warn("没有选取到节点")
	} else if len(indexList) == 1 {
		index := indexList[0]
		node := manager.GetNode(index)
		manager.SetSelectedIndex(index)
		manager.Save()
		exe := run(node.Protocol)
		if exe {
			if setting.Socks() == 0 {
				if setting.Http() == 0 {
					log.Infof("启动成功, 监听mixed端口: %d, 所选节点: %d", setting.Mixed(), manager.SelectedIndex())
				} else {
					log.Infof("启动成功, 监听mixed/http端口: %d/%d, 所选节点: %d", setting.Mixed(), setting.Http(), manager.SelectedIndex())
				}
			} else {
				if setting.Http() == 0 {
					log.Infof("启动成功, 监听mixed/socks端口: %d/%d, 所选节点: %d", setting.Mixed(), setting.Socks(), manager.SelectedIndex())
				} else {
					log.Infof("启动成功, 监听mixed/socks/http端口: %d/%d/%d, 所选节点: %d", setting.Mixed(), setting.Socks(), setting.Http(), manager.SelectedIndex())
				}
			}
			result, status := TestNode(testUrl, setting.Mixed(), testTimeout)
			log.Infof("%6s [ %s ] 延迟: %dms", status, testUrl, result)
		}
	} else {
		min := 100000
		i := -1
		for _, index := range indexList {
			node := manager.GetNode(index)
			exe := run(node.Protocol)
			if exe {
				result, status := TestNode(testUrl, setting.Mixed(), testTimeout)
				log.Infof("%6s [ %s ] 节点: %d, 延迟: %dms", status, testUrl, index, result)
				if result > 0 && result <= setting.TestMinTime() {
					i = index
					min = result
					break
				}
				if result != -1 && min > result {
					i = index
					min = result
				}
			} else {
				return
			}
		}
		if i != -1 {
			log.Info("延迟最小的节点为：", i, "，延迟：", min, "ms")
			manager.SetSelectedIndex(i)
			manager.Save()
			node := manager.GetNode(i)
			exe := run(node.Protocol)
			if exe {
				if setting.Socks() == 0 {
					if setting.Http() == 0 {
						log.Infof("启动成功, 监听mixed端口: %d, 所选节点: %d", setting.Mixed(), manager.SelectedIndex())
					} else {
						log.Infof("启动成功, 监听mixed/http端口: %d/%d, 所选节点: %d", setting.Mixed(), setting.Http(), manager.SelectedIndex())
					}
				} else {
					if setting.Http() == 0 {
						log.Infof("启动成功, 监听mixed/socks端口: %d/%d, 所选节点: %d", setting.Mixed(), setting.Socks(), manager.SelectedIndex())
					} else {
						log.Infof("启动成功, 监听mixed/socks/http端口: %d/%d/%d, 所选节点: %d", setting.Mixed(), setting.Socks(), setting.Http(), manager.SelectedIndex())
					}
				}
			} else {
				log.Error("启动失败")
			}
		} else {
			log.Info("所选节点全部不能访问外网")
		}

	}
}

// run 根据节点协议启动相应的 xray 子进程
// node : 节点的协议配置
func run(node protocols.Protocol) bool {
	Stop()
	switch node.GetProtocolMode() {
	case protocols.ModeShadowSocks, protocols.ModeTrojan, protocols.ModeVMess, protocols.ModeSocks, protocols.ModeVLESS, protocols.ModeVMessAEAD:
		file := GenConfig(node)
		Xray = exec.Command(XrayPath, "-c", file)
	default:
		log.Infof("暂不支持%v协议", node.GetProtocolMode())
		return false
	}
	stdout, _ := Xray.StdoutPipe()
	_ = Xray.Start()
	r := bufio.NewReader(stdout)
	lines := new([]string)
	go readInfo(r, lines)
	status := make(chan struct{})
	go checkProc(Xray, status)
	stopper := time.NewTimer(time.Millisecond * 300)
	select {
	case <-stopper.C:
		setting.SetPid(Xray.Process.Pid)
		return true
	case <-status:
		log.Error("开启xray服务失败, 查看下面报错信息来检查出错问题")
		for _, x := range *lines {
			log.Error(x)
		}
		return false
	}
}

// Stop 停止服务
func Stop() {
	if Xray != nil {
		Xray.Process.Kill()
		Xray = nil
	}
	if setting.Pid() != 0 {
		process, err := os.FindProcess(setting.Pid())
		if err == nil {
			process.Kill()
		}
		setting.SetPid(0)
	}
	// 日志文件过大清除
	file, _ := os.Stat(core.LogFile)
	if file != nil {
		fileSize := float64(file.Size()) / (1 << 20)
		if fileSize > 5 {
			os.Remove(core.LogFile)
		}
	}
}

// 查看xray日志
func ShowLog() {
	t, _ := tail.TailFile(core.LogFile, tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	})
	for line := range t.Lines {
		fmt.Println(line.Text)
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
