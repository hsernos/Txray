package main

import (
	"Txray/cmd"
	"Txray/core"
	"Txray/core/setting"
	"Txray/core/setting/key"
	"Txray/log"
	"github.com/abiosoft/ishell"
	"github.com/spf13/viper"
	"os"
)

const (
	version = "v3.0.0"
	name    = "Txray"
)

func init() {
	// 初始化日志
	absPath := core.PathJoin(core.GetConfigDir(), "info.log")
	log.Init(
		log.GetConsoleZapcore(log.INFO),
		log.GetFileZapcore(absPath, log.INFO, 5),
	)
}

func beforeOfRun(shell *ishell.Shell) {
	cmd := viper.GetString(key.RunBefore)
	if cmd != "" {
		for _, line := range setting.NewAlias("", cmd).GetCmd() {
			shell.Process(line...)
		}
		shell.Print("\n>>> ")
	}
}

func main() {
	shell := ishell.New()
	cmd.InitShell(shell)
	shell.Set("version", version)
	shell.Set("name", name)
	if len(os.Args) > 1 {
		_ = shell.Process(os.Args[1:]...)
	} else {
		go beforeOfRun(shell)
		shell.Printf("%s - Xray Shell Client - %s\n", name, version)
		shell.Run()
	}
}
