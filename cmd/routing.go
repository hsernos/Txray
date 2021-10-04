package cmd

import (
	"Txray/cmd/help"
	"Txray/core/routing"
	"Txray/log"
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"strings"
)

func InitRouteShell(shell *ishell.Shell) {
	routingCmd := &ishell.Cmd{
		Name:    "routing",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			shell.Process("routing", "help")
		},
	}
	routingCmd.AddCmd(&ishell.Cmd{
		Name: "help",
		Func: func(c *ishell.Context) {
			c.Println(help.Routing)
		},
	})
	// block
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
				data, _ := ioutil.ReadFile(fileArg)
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
	// proxy
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
				data, _ := ioutil.ReadFile(fileArg)
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
	// direct
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
				data, _ := ioutil.ReadFile(fileArg)
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
