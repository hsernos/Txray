// cmd/filter.go 负责 shell 层面节点过滤命令的注册与实现
package cmd

import (
	"Txray/cmd/help"
	"Txray/core/manage"
	"Txray/core/node"
	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

// InitFilterShell 初始化过滤器相关的命令
func InitFilterShell(shell *ishell.Shell) {
	filterCmd := &ishell.Cmd{
		Name: "filter",
		Func: func(c *ishell.Context) {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"索引", "规则", "是否启用"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			manage.Manager.FilterForEach(func(i int, filter *node.NodeFilter) {
				table.Append([]string{
					strconv.Itoa(i),
					filter.String(),
					strconv.FormatBool(filter.IsUse),
				})
			})
			table.Render()
		},
	}
	// help
	filterCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Filter)
		},
	})
	// add
	filterCmd.AddCmd(&ishell.Cmd{
		Name: "add",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.AddFilter(c.Args[0])
			_:
				shell.Process("filter")
			}
		},
	})
	// run
	filterCmd.AddCmd(&ishell.Cmd{
		Name: "run",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				manage.Manager.RunFilter("")
			} else {
				manage.Manager.RunFilter(c.Args[0])
			}

		},
	})
	// rm
	filterCmd.AddCmd(&ishell.Cmd{
		Name:    "rm",
		Aliases: []string{"del"},
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.DelFilter(c.Args[0])
			}
			_ = shell.Process("filter")
		},
	})
	// open
	filterCmd.AddCmd(&ishell.Cmd{
		Name: "open",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.SetFilter(c.Args[0], true)
			}
			_ = shell.Process("filter")
		},
	})
	//close
	filterCmd.AddCmd(&ishell.Cmd{
		Name: "close",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.SetFilter(c.Args[0], false)
			}
			_ = shell.Process("filter")
		},
	})
	shell.AddCmd(filterCmd)
}
