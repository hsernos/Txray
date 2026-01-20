// cmd/subscribe.go 负责 shell 层面订阅管理相关命令的注册与实现
package cmd

import (
	"Txray/cmd/help"      // 帮助文档内容
	"Txray/core/manage"   // 节点/订阅管理器
	"Txray/core/sub"      // 订阅相关结构体与方法
	"Txray/log"           // 日志输出
	"github.com/abiosoft/ishell" // shell 框架
	"github.com/olekukonko/tablewriter" // 表格输出
	"os"                  // 系统操作
	"strconv"             // 字符串与数字转换
)

// InitSubscribeShell 向 shell 注册 sub 相关命令（展示、添加、删除、帮助等）
func InitSubscribeShell(shell *ishell.Shell) {
	subCmd := &ishell.Cmd{
		Name: "sub",
		Func: func(c *ishell.Context) {
			// 展示所有订阅信息
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
	// help 子命令，输出帮助内容
	subCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Sub)
		},
	})
	// add 子命令，添加订阅
	subCmd.AddCmd(&ishell.Cmd{
		Name: "add",
		Func: func(c *ishell.Context) {
			// 解析参数，支持 -r 备注
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
					_ = shell.Process("sub") // 添加后刷新列表
				} else {
					log.Warn("需要输入一个订阅链接")
				}
			} else if len(c.Args) == 0 {
				log.Warn("还需要输入一个订阅链接")
			}
		},
	})
	// update-node 子命令，更新节点信息
	subCmd.AddCmd(&ishell.Cmd{
		Name: "update-node",
		Func: func(c *ishell.Context) {
			// 解析参数，支持 -s socks  -h http -a addr
			argMap := FlagsParse(c.Args, map[string]string{
				"s": "socks",
				"h": "http",
				"a": "addr",
			})
			opt := sub.UpdataOption{}
			opt.Key = argMap["data"]
			// 判断是 socks 还是 http
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
			// addr 为必填项
			if address, ok := argMap["addr"]; ok {
				opt.Addr = address
			}
			manage.Manager.UpdataNode(opt)
		},
	})
	// rm 子命令，删除订阅
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
	// mv 子命令，修改订阅信息
	subCmd.AddCmd(&ishell.Cmd{
		Name:    "mv",
		Aliases: []string{"set"},
		Func: func(c *ishell.Context) {
			// 解析参数，支持 -r 备注 -u url
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
