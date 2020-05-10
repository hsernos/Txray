package main

import (
	"github.com/abiosoft/ishell"
	"os"
	"v3ray/cmd"
	log "v3ray/logger"
	"v3ray/tool"
)

const (
	version = "v1.1.2"
	name    = "v3ray"
)

func init() {
	dir := os.Getenv("V3RAY_HOME")
	if !tool.PathExists(tool.Join(dir, "v2ray", "v2ray")) {
		log.Error(tool.Join(dir, "v2ray") + "目录下不存在v2ray可执行文件")
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
