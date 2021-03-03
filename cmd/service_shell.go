package cmd

import (
	nas "Txray/core/node"
	"Txray/core/service"
	"fmt"
	"github.com/abiosoft/ishell"
)

func InitServiceShell(shell *ishell.Shell) {
	// 启动或重启服务
	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "启动或重启服务",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"t": "tcp",
			})
			key := fmt.Sprintf("%d", nas.GetSelectedIndex())
			isTcpSort := false
			if k, ok := argMap["tcp"]; ok {
				if k == "" {
					key = "1"
				} else {
					key = k
				}
				isTcpSort = true
			} else if k, ok := argMap["data"]; ok {
				if k != "" {
					key = k
				}
			}
			service.Start(key, isTcpSort)
		},
	})
	// 停止服务
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "停止服务",
		Func: func(c *ishell.Context) {
			service.Stop()
		},
	})
}
