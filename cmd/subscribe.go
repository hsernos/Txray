package cmd

import (
	"Txray/cmd/help"
	"Txray/core/manage"
	"Txray/core/sub"
	"Txray/log"
	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

func InitSubscribeShell(shell *ishell.Shell) {
	subCmd := &ishell.Cmd{
		Name: "sub",
		Func: func(c *ishell.Context) {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"索引", "别名", "订阅地址", "是否启用"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			manage.Manager.SubForEach(func(i int, subscirbe *sub.Subscirbe) {
				table.Append([]string{
					strconv.Itoa(i),
					subscirbe.Name,
					subscirbe.Url,
					strconv.FormatBool(subscirbe.Using),
				})
			})
			table.Render()
		},
	}
	// help
	subCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Sub)
		},
	})
	// add
	subCmd.AddCmd(&ishell.Cmd{
		Name: "add",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"r": "remarks",
			})
			if len(c.Args) >= 1 {
				if sublink, ok := argMap["data"]; ok {
					if remarksArg, ok := argMap["remarks"]; ok {
						manage.Manager.AddSubscirbe(sub.NewSubscirbe(sublink, remarksArg))
					} else {
						manage.Manager.AddSubscirbe(sub.NewSubscirbe(sublink, "remarks"))
					}
					_ = shell.Process("sub")
				} else {
					log.Warn("需要输入一个订阅链接")
				}
			} else if len(c.Args) == 0 {
				log.Warn("还需要输入一个订阅链接")
			}
		},
	})
	// update-node
	subCmd.AddCmd(&ishell.Cmd{
		Name: "update-node",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"s": "socks",
				"h": "http",
				"a": "addr",
			})
			opt := sub.UpdataOption{}
			opt.Key = argMap["data"]
			if socks, ok := argMap["socks"]; ok {
				if v, err := strconv.Atoi(socks); err == nil {
					if 0 < v && v <= 65535 {
						opt.Port = v
					}
				}
				opt.ProxyMode = sub.SOCKS
			} else if http, ok := argMap["http"]; ok {
				if v, err := strconv.Atoi(http); err == nil {
					if 0 < v && v <= 65535 {
						opt.Port = v
					}
				}
				opt.ProxyMode = sub.HTTP
			}
			if address, ok := argMap["addr"]; ok {
				opt.Addr = address
			}
			manage.Manager.UpdataNode(opt)
		},
	})
	// rm
	subCmd.AddCmd(&ishell.Cmd{
		Name:    "rm",
		Aliases: []string{"del"},
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.DelSub(c.Args[0])
				_ = shell.Process("sub")
			}
		},
	})
	// mv
	subCmd.AddCmd(&ishell.Cmd{
		Name:    "mv",
		Aliases: []string{"set"},
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"r": "remarks",
				"u": "url",
			})
			if key, ok := argMap["data"]; ok {

				url := argMap["url"]
				remarks := argMap["remarks"]

				using := ""
				if value, ok := argMap["using"]; ok {
					using = value
				}
				manage.Manager.SetSub(key, using, url, remarks)
				_ = shell.Process("sub")
			}
		},
	})
	shell.AddCmd(subCmd)
}
