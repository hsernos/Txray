// cmd/recycle.go 负责 shell 层面回收站管理命令的注册与实现
package cmd

import (
	"Txray/cmd/help"      // 帮助文档内容
	"Txray/core"          // 索引工具
	"Txray/core/manage"   // 节点/回收站管理器
	"github.com/abiosoft/ishell" // shell 框架
	"github.com/olekukonko/tablewriter" // 表格输出
	"os"                  // 系统操作
	"strconv"             // 字符串与数字转换
)

// InitRecycleShell 注册 recycle 命令及其子命令，支持回收站节点展示、恢复、帮助等
func InitRecycleShell(shell *ishell.Shell) {
	recycleCmd := &ishell.Cmd{
		Name: "recycle",
		Func: func(c *ishell.Context) {
			// 展示回收站节点列表
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
	// help 子命令
	recycleCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Recycle)
		},
	})
	// restore 子命令，恢复回收站节点
	recycleCmd.AddCmd(&ishell.Cmd{
		Name: "restore",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.MoveFormRecycle(c.Args[0])
			}
		},
	})
	// clear 子命令，清空回收站
	recycleCmd.AddCmd(&ishell.Cmd{
		Name: "clear",
		Func: func(c *ishell.Context) {
			manage.Manager.ClearRecycle()
		},
	})
	// 注册到 shell
	shell.AddCmd(recycleCmd)
}
