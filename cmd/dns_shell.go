package cmd

import (
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"os"
)

func InitDNSShell(shell *ishell.Shell) {
	dns := &ishell.Cmd{
		Name: "dns",
		Help: "DNS设置, 使用dns查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(HelpDNS())
		},
	}
	// 查看DNS设置
	dns.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看DNS设置",
		Func: func(c *ishell.Context) {
			data := []string{
				tools.UintToStr(coreService.DNS.Port),
				coreService.DNS.Outland,
				coreService.DNS.Inland,
				coreService.DNS.Backup,
			}
			format.ShowDNS(os.Stdout, data)
		},
	})
	// 修改DNS设置
	dns.AddCmd(&ishell.Cmd{
		Name: "alter",
		Help: "修改DNS设置",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"p": "port",
				"i": "inland",
				"o": "outland",
				"b": "backup",
			})
			if port, ok := r["port"]; ok && tools.IsNetPort(port) {
				coreService.SetDNSPort(tools.StrToUint(port))
			}
			if inland, ok := r["inland"]; ok {
				coreService.SetInlandDNS(inland)
			}
			if outland, ok := r["outland"]; ok {
				coreService.SetOutlandDNS(outland)
			}
			if backup, ok := r["backup"]; ok {
				coreService.SetBackupDNS(backup)
			}
		},
	})
	shell.AddCmd(dns)
}
