// cmd/shell.go 负责注册所有 shell 命令入口及参数解析工具
package cmd

import (
	"Txray/cmd/help" // 帮助文档内容
	"Txray/xray"     // xray 日志展示
	"github.com/abiosoft/ishell" // shell 框架
)

// InitShell 注册所有一级命令到 shell，包括版本、帮助、设置、节点、订阅、过滤、回收站、路由、服务、别名、日志等
func InitShell(shell *ishell.Shell) {
	shell.AddCmd(&ishell.Cmd{
		Name:    "version",
		Aliases: []string{"-v", "--version"},
		Help:    "程序版本",
		Func: func(c *ishell.Context) {
			c.Printf("%s version \"%s\"\n", c.Get("name"), c.Get("version"))
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Help:    "帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(help.Help)
		},
	})
	// 注册各功能子命令
	InitSettingShell(shell)
	InitNodeShell(shell)
	InitSubscribeShell(shell)
	InitFilterShell(shell)
	InitRecycleShell(shell)
	InitRouteShell(shell)
	InitServiceShell(shell)
	InitAliasShell(shell)
	// 日志查看命令
	shell.AddCmd(&ishell.Cmd{
		Name:    "log",
		Func: func(c *ishell.Context) {
			xray.ShowLog()
		},
	})
}

// FlagsParse 解析 shell 命令参数，支持短/长选项与 data 键
// args: 参数列表，keys: 短选项到长选项映射
// 返回 map，key 为参数名，value 为参数值
func FlagsParse(args []string, keys map[string]string) map[string]string {
	resultMap := make(map[string]string)
	key := "data"
	for _, x := range args {
		if len(x) >= 2 {
			if x[:2] == "--" {
				key = x[2:]
				resultMap[key] = ""
			} else if x[:1] == "-" {
				if x[1] >= 48 && x[1] <= 57 {
					resultMap[key] = x
				} else if len(x) == 2 {
					d, ok := keys[x[1:]]
					if ok {
						key = d
					} else {
						key = x[1:]
					}
					resultMap[key] = ""
				}
			} else {
				resultMap[key] = x
			}
		} else if len(x) > 0 {
			resultMap[key] = x
		}
	}
	return resultMap
}
