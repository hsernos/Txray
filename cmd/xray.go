// cmd/xray.go 负责 shell 层面 xray 服务的启动与停止命令
package cmd

import (
	"Txray/core/manage" // 节点管理器
	"Txray/log"         // 日志输出
	"Txray/xray"        // xray 服务控制
	"strconv"           // 字符串与数字转换

	"github.com/abiosoft/ishell" // shell 框架
)

// InitServiceShell 向 shell 注册 run/stop 命令用于启动和停止 xray 服务
// shell: ishell.Shell 指针，交互式 shell 实例
func InitServiceShell(shell *ishell.Shell) {
	// 启动或重启服务命令
	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "启动或重启服务",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"c": "config",
				"d": "confdir",
			})
			if _, ok := argMap["config"]; ok {
				if argMap["config"] == "" {
					log.Warn("需要指定配置文件，例如: run -c config.json")
					return
				}
				xray.RunConfig(argMap["config"])
				return
			}
			if _, ok := argMap["confdir"]; ok {
				if argMap["confdir"] == "" {
					log.Warn("需要指定配置目录，例如: run -d confdir")
					return
				}
				xray.RunConfDir(argMap["confdir"])
				return
			}
			if key, ok := argMap["data"]; ok {
				xray.Start(key)
			} else {
				xray.Start(strconv.Itoa(manage.Manager.SelectedIndex()))
			}
		},
	})
	// 停止服务命令
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "停止服务",
		Func: func(c *ishell.Context) {
			xray.Stop()
		},
	})
}
