package cmd

import (
	"Txray/core/routing"
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"os"
	"strings"
)

func InitRouteShell(shell *ishell.Shell) {
	routingCmd := &ishell.Cmd{
		Name: "routing",
		Help: "路由管理, 使用route查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(HelpRouting())
		},
	}
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "查看帮助",
		Func: func(c *ishell.Context) {
			c.Println(HelpRouting())
		},
	})
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "block",
		Help: "",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"a": "add",
				"r": "rm",
				"f": "file",
				"c": "clipboard",
			})
			mode := routing.TypeBlock
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
				routing.AddRule(mode, strings.Split(content, "\n")...)
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
				routing.AddRule(mode, list...)
			} else if data, ok := argMap["add"]; ok {
				routing.AddRule(mode, data)
			} else if key, ok := argMap["rm"]; ok {
				routing.DelRule(mode, key)
			} else if key, ok := argMap["data"]; ok {
				format.ShowRouter(os.Stdout, routing.GetRule(mode, key)...)
			} else {
				format.ShowRouter(os.Stdout, routing.GetRule(mode, "all")...)
			}
		},
	})
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "proxy",
		Help: "",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"a": "add",
				"r": "rm",
				"f": "file",
				"c": "clipboard",
			})
			mode := routing.TypeProxy
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
				routing.AddRule(mode, strings.Split(content, "\n")...)
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
				routing.AddRule(mode, list...)
			} else if data, ok := argMap["add"]; ok {
				routing.AddRule(mode, data)
			} else if key, ok := argMap["rm"]; ok {
				routing.DelRule(mode, key)
			} else if key, ok := argMap["data"]; ok {
				format.ShowRouter(os.Stdout, routing.GetRule(mode, key)...)
			} else {
				format.ShowRouter(os.Stdout, routing.GetRule(mode, "all")...)
			}
		},
	})
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "direct",
		Help: "",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"a": "add",
				"r": "rm",
				"f": "file",
				"c": "clipboard",
			})
			mode := routing.TypeDirect
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
				routing.AddRule(mode, strings.Split(content, "\n")...)
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
				routing.AddRule(mode, list...)
			} else if data, ok := argMap["add"]; ok {
				routing.AddRule(mode, data)
			} else if key, ok := argMap["rm"]; ok {
				routing.DelRule(mode, key)
			} else if key, ok := argMap["data"]; ok {
				format.ShowRouter(os.Stdout, routing.GetRule(mode, key)...)
			} else {
				format.ShowRouter(os.Stdout, routing.GetRule(mode, "all")...)
			}
		},
	})
	shell.AddCmd(routingCmd)
}
