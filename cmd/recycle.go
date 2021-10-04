package cmd

import (
	"Txray/cmd/help"
	"Txray/core"
	"Txray/core/manage"
	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

func InitRecycleShell(shell *ishell.Shell) {
	recycleCmd := &ishell.Cmd{
		Name: "recycle",
		Func: func(c *ishell.Context) {
			var key string
			if len(c.Args) == 1 {
				key = c.Args[0]
			} else {
				key = "all"
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"索引", "协议", "别名", "地址", "端口"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			center := tablewriter.ALIGN_CENTER
			left := tablewriter.ALIGN_LEFT
			table.SetColumnAlignment([]int{center, center, left, center, center, center})
			table.SetColWidth(70)
			for _, index := range core.IndexList(key, manage.Manager.RecycleLen()) {
				n := manage.Manager.GetRecycleNode(index)
				if n != nil {
					table.Append([]string{
						strconv.Itoa(index),
						string(n.GetProtocolMode()),
						n.GetName(),
						n.GetAddr(),
						strconv.Itoa(n.GetPort()),
					})
				}
			}
			table.Render()
		},
	}
	// help
	recycleCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Recycle)
		},
	})
	// restore
	recycleCmd.AddCmd(&ishell.Cmd{
		Name: "restore",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.MoveFormRecycle(c.Args[0])
			}
		},
	})
	// clear
	recycleCmd.AddCmd(&ishell.Cmd{
		Name: "clear",
		Func: func(c *ishell.Context) {
			manage.Manager.ClearRecycle()
		},
	})
	shell.AddCmd(recycleCmd)
}
