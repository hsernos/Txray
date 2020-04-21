package main

import (
	"github.com/abiosoft/ishell"
	"os"
	"v3ray/cmd"
	log "v3ray/logger"
	"v3ray/tool"
)


const (
    version = "v1.0.0"
    name = "v3ray"
)


func main(){
    shell := ishell.New()
    shell.Printf("%s - V2ray Shell Client\n", name)
    if !tool.CheckPATH("v2ray") {
    	log.Error("请将v2ray程序所在目录添加到PATH环境变量中")
    	os.Exit(1)
	}
	cmd.InitShell(shell)
	cmd.InitConfig()
    // run shell
	shell.Run()
	defer cmd.Kill()
}


