package main

import (
	"Txray/cmd"
	"Txray/core/config"
	"Txray/core/service"
	"Txray/log"
	"Txray/tools"
	"github.com/abiosoft/ishell"
	"os"
)

const (
	version = "v2.1.0"
	name    = "Txray"
)

func init() {
	// 初始化日志
	absPath := tools.PathJoin(config.GetConfigDir(), "info.log")
	log.Init(
		log.GetConsoleZapcore(log.INFO),
		log.GetFileZapcore(absPath, log.INFO, 5),
	)
}

func main() {
	shell := ishell.New()
	cmd.InitShell(shell)
	shell.Set("version", version)
	shell.Set("name", name)
	if len(os.Args) > 1 {
		_ = shell.Process(os.Args[1:]...)
	} else {
		shell.Printf("%s - Xray Shell Client - %s\n", name, version)
		shell.Run()
	}
	defer service.Stop()
}
