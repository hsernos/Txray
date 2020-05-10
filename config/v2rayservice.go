package config

import (
	"bufio"
	"os/exec"
	"time"
	log "v3ray/logger"
	"v3ray/tool"
)

// Start 开启v2ray服务
func (c *Config) Start(i int) {
	if i >= len(c.Nodes) {
		log.Warn("没有该节点")
		return
	} else if i < 0 {
		if int(c.GetNodeIndex()) >= len(c.Nodes) {
			log.Warn("该节点已删除")
			return
		}
	} else {
		c.Index = uint(i)
	}
	c.SaveJSON()
	c.GenConfig()
	if c.exeCmd == nil {
		c.exeCmd = exec.Command(tool.Join(c.getConfigPath(), "v2ray", "v2ray"), "--config", tool.Join(c.getConfigPath(), "config.json"))
	}
	stdout, _ := c.exeCmd.StdoutPipe()
	c.exeCmd.Start()
	r := bufio.NewReader(stdout)
	lines := new([]string)
	go readInfo(r, lines)
	code := new(int)
	go checkProc(c.exeCmd, code)
	time.Sleep(time.Duration(500) * time.Millisecond)

	if *code > 0 {
		log.Error("开启v2ray服务失败,查看下面报错信息来检查出错问题")
		for _, x := range *lines {
			log.Error(x)
		}
	} else {
		log.Info("开启v2ray服务成功, 监听socks5/http端口：", c.Settings.Port, "/", c.Settings.Http, "，选定节点索引：", c.Index)
	}
}

// Stop 停止v2ray服务
func (c *Config) Stop() {
	if c.exeCmd != nil {
		c.exeCmd.Process.Kill()
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
