package config

import (
	log "Tv2ray/logger"
	"Tv2ray/tools"
	"bufio"
	"os/exec"
	"path/filepath"
	"time"
)

// 检查v2ray程序和资源文件 geoip.dat、geosite.dat是否存在且在同一目录
func CheckV2rayFile() bool {
	v2rayPath := GetV2rayPath()
	if v2rayPath == "" {
		log.Error("在", tools.GetRunPath(), " 下没有找到v2ray程序")
		log.Error("请在 https://github.com/v2fly/v2ray-core/releases 下载对应版本")
		log.Error("并将解压后的文件夹或所有文件移动到 ", tools.GetRunPath(), " 下")
		return false
	} else {
		path := filepath.Dir(GetV2rayPath())
		if tools.IsFile(tools.PathJoin(path, "geoip.dat")) && tools.IsFile(tools.PathJoin(path, "geosite.dat")) {
			return true
		} else {
			log.Error("在 ", path, " 下没有找到v2ray程序的资源文件 geoip.dat 和 geosite.dat")
			log.Error("请在 https://github.com/v2fly/v2ray-core/releases 下载对应版本")
			log.Error("并将缺失的文件移动到 ", path, " 下")
			return false
		}
	}
}

// 查找v2ray程序所在绝对路径
func GetV2rayPath() string {
	path := tools.GetRunPath()
	files, _ := tools.FindFileByName(path, "v2ray", ".exe")
	if len(files) == 0 {
		return ""
	}
	return files[0]
}

// 开启v2ray服务
func (c *Config) Start(i int) bool {
	if i >= len(c.Nodes) {
		log.Warn("没有该节点")
		return false
	} else if i < 0 {
		if int(c.GetNodeIndex()) >= len(c.Nodes) {
			log.Warn("该节点已删除")
			return false
		}
	} else {
		c.Index = uint(i)
	}
	c.SaveJSON()
	c.GenConfig()
	if c.exeCmd == nil {
		c.exeCmd = exec.Command(GetV2rayPath(), "--config", tools.PathJoin(tools.GetRunPath(), "config.json"))
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
		return false
	} else {
		log.Info("开启v2ray服务成功, 监听socks5/http端口：", c.Settings.Port, "/", c.Settings.Http, "，选定节点索引：", c.Index)
	}
	return true
}

// 停止v2ray服务
func (c *Config) Stop() {
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
