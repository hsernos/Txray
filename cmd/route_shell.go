package cmd

import (
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"os"
	"strings"
)

func InitRouteShell(shell *ishell.Shell) {
	route := &ishell.Cmd{
		Name: "route",
		Help: "路由管理, 使用route查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(HelpRoute())
		},
	}
	// 添加路由规则
	route.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "添加路由规则",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"b": "block",
				"d": "direct",
				"p": "proxy",
				"f": "file",
				"c": "clipboard",
			})
			mode := "none"
			if _, ok := argMap["block"]; ok {
				mode = "block"
			} else if _, ok := argMap["direct"]; ok {
				mode = "direct"
			} else if _, ok := argMap["proxy"]; ok {
				mode = "proxy"
			} else {
				log.Warn("必须使用 --block、--direct、--proxy 中的一个")
				return
			}
			if _, ok := argMap["clipboard"]; ok {
				content, err := clipboard.ReadAll()
				if err != nil {
					log.Error(err)
					return
				}
				c.Println("剪贴板内容如下：")
				c.Println("========================================================================")
				c.Println(content)
				c.Println("========================================================================")
				switch mode {
				case "block":
					coreService.AddBlock(strings.Split(content, "\n")...)
				case "direct":
					coreService.AddDirect(strings.Split(content, "\n")...)
				case "proxy":
					coreService.AddProxy(strings.Split(content, "\n")...)
				}
			} else if fileArg, ok := argMap["file"]; ok {
				if !tools.IsFile(fileArg) {
					log.Error("open ", fileArg, " : 没有这个文件")
					return
				}
				list := tools.ReadFile(fileArg)
				c.Println("文件内容如下：")
				c.Println("========================================================================")
				c.Println(strings.Join(list, "\n"))
				c.Println("========================================================================")
				switch mode {
				case "block":
					coreService.AddBlock(list...)
				case "direct":
					coreService.AddDirect(list...)
				case "proxy":
					coreService.AddProxy(list...)
				}
			} else if data, ok := argMap["data"]; ok {
				switch mode {
				case "block":
					coreService.AddBlock(data)
				case "direct":
					coreService.AddDirect(data)
				case "proxy":
					coreService.AddProxy(data)
				}
			}
		},
	})
	// 查看路由规则
	route.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看路由规则",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"b": "block",
				"d": "direct",
				"p": "proxy",
			})
			key := "all"
			var data [][]string
			if k, ok := argMap["block"]; ok {
				if k != "" {
					key = k
				}
				data = coreService.GetBlock(key)
			} else if k, ok := argMap["direct"]; ok {
				if k != "" {
					key = k
				}
				data = coreService.GetDirect(key)
			} else if k, ok := argMap["proxy"]; ok {
				if k != "" {
					key = k
				}
				data = coreService.GetProxy(key)
			} else {
				log.Warn("必须使用 --block、--direct、--proxy 中的一个")
				return
			}
			format.ShowRouter(os.Stdout, data...)
		},
	})
	// 删除路由规则
	route.AddCmd(&ishell.Cmd{
		Name: "del",
		Help: "删除路由规则",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"b": "block",
				"d": "direct",
				"p": "proxy",
			})
			if key, ok := argMap["block"]; ok {
				coreService.DelBlock(key)
			} else if key, ok := argMap["direct"]; ok {
				coreService.DelDirect(key)
			} else if key, ok := argMap["proxy"]; ok {
				coreService.DelProxy(key)
			} else {
				log.Warn("必须使用 --block、--direct、--proxy 中的一个")
				return
			}
		},
	})
	shell.AddCmd(route)
}
