// cmd/alias.go 负责 shell 层面别名命令的注册与实现
package cmd

import (
	"Txray/cmd/help"
	"Txray/core/setting"
	"Txray/log"
	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

var names = make(map[string]int)

// InitAliasShell 初始化别名命令
// 将内置命令加载到 names 中，
// 并读取已设置的别名，注册 alias 命令及其子命令
func InitAliasShell(shell *ishell.Shell) {
	// 读取内置命令
	for _, cmd := range shell.Cmds() {
		names[cmd.Name] = 0
	}
	LoadAlias(shell)
	aliasCmd := &ishell.Cmd{
		Name: "alias",
		Func: func(c *ishell.Context) {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"索引", "别名", "命令"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			table.SetColWidth(70)
			center := tablewriter.ALIGN_CENTER
			left := tablewriter.ALIGN_LEFT
			table.SetColumnAlignment([]int{center, center, left})
			for i, alias := range setting.AliasList() {
				table.Append([]string{
					strconv.Itoa(i + 1),
					alias.Name,
					alias.Cmd,
				})
			}
			table.Render()
		},
	}
	aliasCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Alias)
		},
	})
	aliasCmd.AddCmd(&ishell.Cmd{
		Name:    "set",
		Aliases: []string{"add"},
		Func: func(c *ishell.Context) {
			if len(c.Args) == 2 {
				if _, ok := names[c.Args[0]]; !ok {
					setting.AddAlias(c.Args[0], c.Args[1])
					LoadAlias(shell)
					_ = shell.Process("alias")
				} else {
					log.Warnf("'%s' 是内置命令，不能被覆盖", c.Args[0])
				}
			}
		},
	})
	aliasCmd.AddCmd(&ishell.Cmd{
		Name: "rm",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				for _, alias := range setting.DelAlias(c.Args[0]) {
					shell.DeleteCmd(alias)
				}
				_ = shell.Process("alias")
			}
		},
	})
	shell.AddCmd(aliasCmd)
}

// LoadAlias 加载所有别名
// 遍历设置的别名列表，为每个别名调用 AddAliasShell 函数
func LoadAlias(shell *ishell.Shell) {
	for _, a := range setting.AliasList() {
		AddAliasShell(shell, a)
	}
}

// AddAliasShell 为单个别名注册 shell 命令
// 如果别名不为 nil 且不与内置命令冲突，则为该别名注册一个新的 shell 命令
func AddAliasShell(shell *ishell.Shell, a *setting.Alias) {
	if a != nil {
		// 防止覆盖内置命令
		if _, ok := names[a.Name]; !ok {
			shell.AddCmd(&ishell.Cmd{
				Name: a.Name,
				Func: func(c *ishell.Context) {
					for i, line := range a.GetCmd() {
						if i+1 == len(a.GetCmd()) {
							line = append(line, c.Args...)
						}
						shell.Process(line...)
					}
				},
			})
		}
	}
}
