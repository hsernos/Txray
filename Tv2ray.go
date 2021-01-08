package main

import (
	"Tv2ray/cmd"
	"Tv2ray/config"
	"github.com/abiosoft/ishell"
	"os"
)

const (
	version = "v1.2.0"
	name    = "Tv2ray"
)

func init() {
	// 检测v2ray核心文件
	if !config.CheckV2rayFile() {
		os.Exit(1)
	}
}

func main() {
	shell := ishell.New()
	shell.Printf("%s - V2ray Shell Client - %s\n", name, version)
	cmd.InitShell(shell)
	cmd.InitConfig()
	shell.Run()
	defer cmd.Kill()
}
