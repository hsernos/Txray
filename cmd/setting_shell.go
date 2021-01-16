package cmd

import (
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"os"
)

func InitSettingShell(shell *ishell.Shell) {
	setting := &ishell.Cmd{
		Name: "setting",
		Help: "基础设置, 使用setting查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(HelpSetting())
		},
	}
	// 查看基础设置
	setting.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看基础设置",
		Func: func(c *ishell.Context) {
			s := coreService.Settings
			data := []string{tools.UintToStr(s.Socks),
				tools.UintToStr(s.Http),
				tools.BoolToStr(s.UDP),
				tools.BoolToStr(s.Sniffing),
				tools.BoolToStr(s.Mux),
				tools.BoolToStr(s.AllowLANConn),
				tools.BoolToStr(s.BypassLanAndContinent),
				s.DomainStrategy,
			}
			format.ShowSetting(os.Stdout, data)
		},
	})
	// 修改基础设置
	setting.AddCmd(&ishell.Cmd{
		Name: "alter",
		Help: "修改基础设置",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"p": "port",
				"h": "http",
				"u": "udp",
				"s": "sniffing",
				"l": "lanconn",
				"m": "mux",
				"b": "bypass",
				"r": "route",
			})

			if port, ok := r["port"]; ok && tools.IsNetPort(port) {
				coreService.SetSocksPort(tools.StrToUint(port))
			}
			if port, ok := r["http"]; ok && tools.IsNetPort(port) {
				coreService.SetHttpPort(tools.StrToUint(port))
			}
			if key, ok := r["udp"]; ok {
				if key == "y" {
					coreService.SetUDP(true)
				} else if key == "n" {
					coreService.SetUDP(false)
				}
			}
			if key, ok := r["sniffing"]; ok {
				if key == "y" {
					coreService.SetSniffing(true)
				} else if key == "n" {
					coreService.SetSniffing(false)
				}
			}
			if key, ok := r["lanconn"]; ok {
				if key == "y" {
					coreService.SetLANConn(true)
				} else if key == "n" {
					coreService.SetLANConn(false)
				}
			}
			if key, ok := r["mux"]; ok {
				if key == "y" {
					coreService.SetMux(true)
				} else if key == "n" {
					coreService.SetMux(false)
				}
			}
			if key, ok := r["bypass"]; ok {
				if key == "y" {
					coreService.SetBypassLanAndContinent(true)
				} else if key == "n" {
					coreService.SetBypassLanAndContinent(false)
				}
			}
			if key, ok := r["route"]; ok {
				if key == "1" {
					coreService.SetDomainStrategy(1)
				} else if key == "2" {
					coreService.SetDomainStrategy(2)
				} else if key == "3" {
					coreService.SetDomainStrategy(3)
				}
			}
		},
	})
	shell.AddCmd(setting)
}
