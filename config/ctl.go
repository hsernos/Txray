package config

import (
	"os/exec"
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
		c.exeCmd = exec.Command("v2ray", "--config", tool.Join(c.getConfigPath() , "config.json"))
	}
	c.exeCmd.Start()
	log.Info("v2ray --config ", tool.Join(c.getConfigPath() , "config.json"))
	log.Info("开启v2ray服务, 监听端口：",c.Settings.Port,"，选定节点索引为：",c.Index)
	
}

// Stop 停止v2ray服务
func (c *Config) Stop() {
	if c.exeCmd != nil {
		c.exeCmd.Process.Kill()
		c.exeCmd = nil
	}
}
