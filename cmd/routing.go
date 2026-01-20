// cmd/routing.go 负责 shell 层面路由规则管理命令的注册与实现
package cmd

import (
	"Txray/cmd/help"     // 帮助文档内容
	"Txray/core/routing" // 路由规则管理
	"Txray/log"          // 日志
	"fmt"                // 格式化输出
	"os"                 // 系统操作
	"strings"            // 字符串处理

	"github.com/abiosoft/ishell"        // shell 框架
	"github.com/atotto/clipboard"       // 剪贴板操作
	"github.com/olekukonko/tablewriter" // 表格输出
)

// InitRouteShell 注册 routing 命令及其子命令，支持规则展示、添加、删除、帮助等
func InitRouteShell(shell *ishell.Shell) {
	routingCmd := &ishell.Cmd{
		Name:    "routing",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			shell.Process("routing", "help") // 默认输出帮助
		},
	}
	// help 子命令
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "help",
		Func: func(c *ishell.Context) {
			c.Println(help.Routing)
		},
	})
	// block 子命令，支持从剪贴板、文件、参数添加规则
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "block",
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
				content = strings.ReplaceAll(content, "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				routing.AddRule(mode, strings.Split(content, "\n")...)
			} else if fileArg, ok := argMap["file"]; ok {
				if _, err := os.Stat(fileArg); os.IsNotExist(err) {
					log.Error("open ", fileArg, " : 没有这个文件")
					return
				}
				data, _ := os.ReadFile(fileArg)
				content := strings.ReplaceAll(string(data), "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				count := routing.AddRule(mode, strings.Split(content, "\n")...)
				log.Infof("共添加了%d条规则", count)
			} else if data, ok := argMap["add"]; ok {
				routing.AddRule(mode, data)
			} else if key, ok := argMap["rm"]; ok {
				routing.DelRule(mode, key)
			} else if key, ok := argMap["data"]; ok {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"索引", "类型", "规则"})
				center := tablewriter.ALIGN_CENTER
				table.SetAlignment(center)
				table.AppendBulk(routing.GetRule(mode, key))
				table.SetCaption(true, fmt.Sprintf("共[ %d ]条规则", routing.RuleLen(mode)))
				table.Render()
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"索引", "类型", "规则"})
				center := tablewriter.ALIGN_CENTER
				table.SetAlignment(center)
				table.AppendBulk(routing.GetRule(mode, "0-100"))
				table.SetCaption(true, fmt.Sprintf("共[ %d ]条规则", routing.RuleLen(mode)))
				table.Render()
			}
		},
	})
	// proxy 子命令，支持从剪贴板、文件、参数添加规则
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "proxy",
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
				content = strings.ReplaceAll(content, "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				routing.AddRule(mode, strings.Split(content, "\n")...)
			} else if fileArg, ok := argMap["file"]; ok {
				if _, err := os.Stat(fileArg); os.IsNotExist(err) {
					log.Error("open ", fileArg, " : 没有这个文件")
					return
				}
				data, _ := os.ReadFile(fileArg)
				content := strings.ReplaceAll(string(data), "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				count := routing.AddRule(mode, strings.Split(content, "\n")...)
				log.Infof("共添加了%d条规则", count)
			} else if data, ok := argMap["add"]; ok {
				routing.AddRule(mode, data)
			} else if key, ok := argMap["rm"]; ok {
				routing.DelRule(mode, key)
			} else if key, ok := argMap["data"]; ok {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"索引", "类型", "规则"})
				center := tablewriter.ALIGN_CENTER
				table.SetAlignment(center)
				table.AppendBulk(routing.GetRule(mode, key))
				table.SetCaption(true, fmt.Sprintf("共[ %d ]条规则", routing.RuleLen(mode)))
				table.Render()
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"索引", "类型", "规则"})
				center := tablewriter.ALIGN_CENTER
				table.SetAlignment(center)
				table.AppendBulk(routing.GetRule(mode, "0-100"))
				table.SetCaption(true, fmt.Sprintf("共[ %d ]条规则", routing.RuleLen(mode)))
				table.Render()
			}
		},
	})
	// direct 子命令，支持从剪贴板、文件、参数添加规则
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "direct",
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
				content = strings.ReplaceAll(content, "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				routing.AddRule(mode, strings.Split(content, "\n")...)
			} else if fileArg, ok := argMap["file"]; ok {
				if _, err := os.Stat(fileArg); os.IsNotExist(err) {
					log.Error("open ", fileArg, " : 没有这个文件")
					return
				}
				data, _ := os.ReadFile(fileArg)
				content := strings.ReplaceAll(string(data), "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				count := routing.AddRule(mode, strings.Split(content, "\n")...)
				log.Infof("共添加了%d条规则", count)
			} else if data, ok := argMap["add"]; ok {
				routing.AddRule(mode, data)
			} else if key, ok := argMap["rm"]; ok {
				routing.DelRule(mode, key)
			} else if key, ok := argMap["data"]; ok {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"索引", "类型", "规则"})
				center := tablewriter.ALIGN_CENTER
				table.SetAlignment(center)
				table.AppendBulk(routing.GetRule(mode, key))
				table.SetCaption(true, fmt.Sprintf("共[ %d ]条规则", routing.RuleLen(mode)))
				table.Render()
			} else {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"索引", "类型", "规则"})
				center := tablewriter.ALIGN_CENTER
				table.SetAlignment(center)
				table.AppendBulk(routing.GetRule(mode, "0-100"))
				table.SetCaption(true, fmt.Sprintf("共[ %d ]条规则", routing.RuleLen(mode)))
				table.Render()
			}
		},
	})
	shell.AddCmd(routingCmd)
}
