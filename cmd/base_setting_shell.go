package cmd

import (
	"Txray/core/setting"
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"os"
)

func InitSettingShell(shell *ishell.Shell) {
	baseSettingCmd := &ishell.Cmd{
		Name: "base",
		Help: "查看基础设置",
		Func: func(c *ishell.Context) {
			base := setting.Base()
			data := []string{tools.UintToStr(base.Socks),
				tools.UintToStr(base.Http),
				tools.BoolToStr(base.UDP),
				tools.BoolToStr(base.Sniffing),
				tools.BoolToStr(base.Mux),
				tools.BoolToStr(base.AllowLANConn),
				tools.BoolToStr(base.BypassLanAndContinent),
				base.DomainStrategy,
			}
			format.ShowSetting(os.Stdout, data)
		},
	}
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "基础设置help文档",
		Func: func(c *ishell.Context) {
			c.Println(HelpSetting())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "socks",
		Help: "修改socks端口",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 && tools.IsNetPort(c.Args[0]) {
				setting.SetSocksPort(tools.StrToUint(c.Args[0]))
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "http",
		Help: "修改http端口",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 && tools.IsNetPort(c.Args[0]) {
				setting.SetHttpPort(tools.StrToUint(c.Args[0]))
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "udp",
		Help: "是否启用udp",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				if key == "y" {
					setting.SetUDP(true)
				} else if key == "n" {
					setting.SetUDP(false)
				}
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "sniffing",
		Help: "是否启用流量监听",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				if key == "y" {
					setting.SetSniffing(true)
				} else if key == "n" {
					setting.SetSniffing(false)
				}
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "lanconn",
		Help: "是否启用局域网连接",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				if key == "y" {
					setting.SetLANConn(true)
				} else if key == "n" {
					setting.SetLANConn(false)
				}
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "mux",
		Help: "是否启用多路复用",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				if key == "y" {
					setting.SetMux(true)
				} else if key == "n" {
					setting.SetMux(false)
				}
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "bypass",
		Help: "是否绕过局域网及大陆",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				if key == "y" {
					setting.SetBypassLanAndContinent(true)
				} else if key == "n" {
					setting.SetBypassLanAndContinent(false)
				}
			}
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "routing",
		Help: "设置路由策略",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				if key == "1" {
					setting.SetDomainStrategy(1)
				} else if key == "2" {
					setting.SetDomainStrategy(2)
				} else if key == "3" {
					setting.SetDomainStrategy(3)
				}
			}
		},
	})
	shell.AddCmd(baseSettingCmd)
}
