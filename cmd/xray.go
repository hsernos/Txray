// cmd/xray.go 负责 shell 层面 xray 服务的启动与停止命令
package cmd

import (
	"Txray/core/manage" // 节点管理器
	"Txray/xray"        // xray 服务控制
	"github.com/abiosoft/ishell" // shell 框架
	"strconv"           // 字符串与数字转换
)

// InitServiceShell 向 shell 注册 run/stop 命令用于启动和停止 xray 服务
// shell: ishell.Shell 指针，交互式 shell 实例
func InitServiceShell(shell *ishell.Shell) {
	// 启动或重启服务命令
	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "启动或重启服务",
		Func: func(c *ishell.Context) {
			// 若有参数则用参数，否则用当前选中节点
			if len(c.Args) == 1 {
				xray.Start(c.Args[0])
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
