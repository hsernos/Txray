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
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/hpcloud/tail"
)

var Xray *exec.Cmd
var readInfoMu sync.Mutex

// TestNodes 依次启动各节点代理并测试真实连接延迟
func TestNodes() {
	testUrl := setting.TestUrl()
	testTimeout := setting.TestTimeout()
	manager := manage.Manager
	indexList := core.IndexList("all", manager.NodeLen())
	for _, index := range indexList {
		n := manager.GetNode(index)
		if n == nil {
			continue
		}
		if run(n.Protocol) {
			result, status := TestNode(testUrl, setting.Mixed(), testTimeout)
			log.Infof("%6s [ %s ] 节点: %d, 延迟: %dms", status, testUrl, index, result)
			if result == -1 {
				n.ConnDelay = 99999
			} else {
				n.ConnDelay = float64(result)
			}
		} else {
			n.ConnDelay = 99999
		}
	}
	Stop()
}

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

// RunConfig 使用指定配置文件启动 xray-core
func RunConfig(configPath string) bool {
	if !filepath.IsAbs(configPath) {
		if abs, err := filepath.Abs(configPath); err == nil {
			configPath = abs
		}
	}
	if !IsExistFile(configPath) {
		log.Errorf("配置文件不存在: %s", configPath)
		return false
	}
	if startXrayCmd([]string{"run", "-c", configPath}, false) {
		log.Infof("启动成功, 使用配置文件: %s", configPath)
		return true
	}
	return false
}

// RunConfDir 使用指定配置目录启动 xray-core
func RunConfDir(dir string) bool {
	if !filepath.IsAbs(dir) {
		if abs, err := filepath.Abs(dir); err == nil {
			dir = abs
		}
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		log.Errorf("配置目录不存在: %s", dir)
		return false
	}
	if startXrayCmd([]string{"run", "-confdir", dir}, false) {
		log.Infof("启动成功, 使用配置目录: %s", dir)
		return true
	}
	return false
}

// run 根据节点协议启动相应的 xray 子进程
// node : 节点的协议配置
func run(node protocols.Protocol) bool {
	switch node.GetProtocolMode() {
	case protocols.ModeShadowSocks, protocols.ModeTrojan, protocols.ModeVMess, protocols.ModeSocks, protocols.ModeVLESS, protocols.ModeVMessAEAD:
		file := GenConfig(node)
		if file == "" {
			return false
		}
		return startXrayCmd([]string{"run", "-c", file}, true)
	default:
		log.Infof("暂不支持%v协议", node.GetProtocolMode())
		return false
	}
}

func startXrayCmd(args []string, waitMixedPort bool) bool {
	Stop()
	if waitMixedPort {
		waitTCPClosed("127.0.0.1", setting.Mixed(), 2*time.Second)
	}
	detach := shouldDetachAfterStart(waitMixedPort)
	Xray = exec.Command(XrayPath, args...)
	applyDaemonAttrs(Xray)
	lines := new([]string)
	var stdout io.Closer
	var stderr io.Closer
	// 后台脱离时不要使用 StdoutPipe/StderrPipe。父进程退出后管道会断开，
	// xray 后续写日志可能收到 broken pipe 并提前退出。
	if detach {
		nullOut, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			log.Error(err)
			Xray = nil
			return false
		}
		stdout = nullOut
		stderr = nullOut
		Xray.Stdout = nullOut
		Xray.Stderr = nullOut
	} else {
		out, err := Xray.StdoutPipe()
		if err != nil {
			log.Error(err)
			Xray = nil
			return false
		}
		stdout = out
		errOut, err := Xray.StderrPipe()
		if err != nil {
			log.Error(err)
			Xray = nil
			return false
		}
		stderr = errOut
	}
	if err := Xray.Start(); err != nil {
		log.Error(err)
		Xray = nil
		return false
	}
	if waitMixedPort || !detach {
		if rc, ok := stdout.(io.ReadCloser); ok {
			go readInfo(bufio.NewReader(rc), lines)
		}
		if rc, ok := stderr.(io.ReadCloser); ok {
			go readInfo(bufio.NewReader(rc), lines)
		}
	}
	var status chan struct{}
	if !detach {
		status = make(chan struct{}, 1)
		go checkProc(Xray, status)
	}
	if waitMixedPort {
		return waitDetachedMixedStart(Xray.Process.Pid, stdout, 5*time.Second)
	}
	return waitConfigXrayStart(status, lines, stdout, detach)
}

// waitConfigXrayStart 等待 run -c/-d 启动：尽快返回，进程提前退出则立即失败。
func waitConfigXrayStart(status <-chan struct{}, lines *[]string, stdout io.Closer, detach bool) bool {
	const (
		minReady = 300 * time.Millisecond
		interval = 50 * time.Millisecond
	)
	readyAt := time.Now().Add(minReady)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-status:
			logStartupError(lines)
			cleanupXrayStart(stdout, detach)
			return false
		case <-ticker.C:
			if time.Now().Before(readyAt) {
				continue
			}
			if Xray == nil || Xray.Process == nil || !processAlive(Xray.Process.Pid) {
				logStartupError(lines)
				cleanupXrayStart(stdout, detach)
				return false
			}
			return finishXrayStart(Xray.Process.Pid, stdout, detach)
		}
	}
}

func waitDetachedMixedStart(pid int, stdout io.Closer, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	address := fmt.Sprintf("%s:%d", "127.0.0.1", setting.Mixed())
	for time.Now().Before(deadline) {
		if !processAlive(pid) {
			log.Error("xray 启动失败或已退出")
			cleanupXrayStart(stdout, true)
			return false
		}
		if socksHandshake(address) {
			return finishXrayStart(pid, stdout, true)
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Errorf("xray service start timeout, mixed port %d is not ready", setting.Mixed())
	Stop()
	if stdout != nil {
		_ = stdout.Close()
	}
	return false
}

// shouldDetachAfterStart 节点 run 始终脱离；run -c/-d 仅在命令行非交互模式下脱离。
func shouldDetachAfterStart(waitMixedPort bool) bool {
	if waitMixedPort {
		return true
	}
	return len(os.Args) > 1
}

// finishXrayStart 保存 pid；detach 为 true 时将子进程脱离父进程，避免 Txray 退出时连带终止。
func finishXrayStart(pid int, stdout io.Closer, detach bool) bool {
	setting.SetPid(pid)
	if detach {
		detachXrayChild(stdout)
	}
	return true
}

func cleanupXrayStart(stdout io.Closer, release bool) {
	if stdout != nil {
		_ = stdout.Close()
	}
	if release && Xray != nil && Xray.Process != nil {
		_ = Xray.Process.Release()
	}
	Xray = nil
}

func applyDaemonAttrs(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		return
	}
	attr := &syscall.SysProcAttr{}
	if f := reflect.ValueOf(attr).Elem().FieldByName("Setsid"); f.IsValid() && f.Kind() == reflect.Bool {
		f.SetBool(true)
	}
	cmd.SysProcAttr = attr
}

// detachXrayChild 关闭可安全关闭的输出句柄并 Release 子进程，由 setting.pid 负责后续 stop。
func detachXrayChild(stdout io.Closer) {
	if stdout != nil {
		_ = stdout.Close()
	}
	if Xray != nil && Xray.Process != nil {
		_ = Xray.Process.Release()
	}
	Xray = nil
}

func logStartupError(lines *[]string) {
	log.Error("开启xray服务失败, 查看下面报错信息来检查出错问题")
	for _, x := range *lines {
		log.Error(x)
	}
}

// Stop 停止服务
// waitTCPClosed waits for the old listener to release the mixed port.
func waitTCPClosed(host string, port int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	address := fmt.Sprintf("%s:%d", host, port)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err != nil {
			return
		}
		conn.Close()
		time.Sleep(100 * time.Millisecond)
	}
}

func socksHandshake(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
		return false
	}
	if _, err := conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
		return false
	}
	buf := make([]byte, 2)
	if _, err := conn.Read(buf); err != nil {
		return false
	}
	return buf[0] == 0x05 && buf[1] == 0x00
}

// Stop stops the running xray service.
func Stop() {
	if Xray != nil {
		if Xray.Process != nil {
			Xray.Process.Kill()
		}
		Xray = nil
	}
	if pid := setting.Pid(); pid != 0 {
		if pid == os.Getpid() {
			log.Warn("setting.pid 指向当前 Txray 进程，已清理并跳过停止")
			_ = setting.SetPid(0)
			return
		}
		if processAlive(pid) {
			process, err := os.FindProcess(pid)
			if err == nil {
				if err := process.Kill(); err != nil {
					log.Warn(err)
				}
			}
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
	const maxStartupLines = 20
	for {
		line, _, err := r.ReadLine()
		if len(line) != 0 {
			readInfoMu.Lock()
			if len(*lines) < maxStartupLines {
				*lines = append(*lines, string(line[:]))
			}
			readInfoMu.Unlock()
		}
		if err != nil {
			return
		}
	}
}

// 检查进程状态
func checkProc(c *exec.Cmd, status chan struct{}) {
	c.Wait()
	status <- struct{}{}
}
