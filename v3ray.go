package main

import (
	"github.com/abiosoft/ishell"
	"os"
	"v3ray/cmd"
	log "v3ray/logger"
	"v3ray/tool"
)

const (
	version = "v1.1.3"
	name    = "v3ray"
)

func init() {
	dir := os.Getenv("V3RAY_HOME")
	if !tool.PathExists(tool.Join(dir, "v2ray", "v2ray")) {
		log.Error(tool.Join(dir, "v2ray") + " 目录下不存在v2ray可执行文件")
		log.Error("请在 https://github.com/v2ray/v2ray-core/releases 下载对应版本")
		log.Error("并将解压后的目录下的所有文件移动到 ", tool.Join(dir, "v2ray"), " 文件夹下")
		os.Exit(1)
	}
}

func main() {
	shell := ishell.New()
	shell.Printf("%s - V2ray Shell Client - %s\n", name, version)
	cmd.InitShell(shell)
	cmd.InitConfig()
	// run shell
	shell.Run()
	defer cmd.Kill()
}
