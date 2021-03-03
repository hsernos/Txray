package cmd

import (
	"Txray/core/setting"
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"os"
	"strings"
)

func InitDNSShell(shell *ishell.Shell) {
	dnsCmd := &ishell.Cmd{
		Name: "dns",
		Help: "查看DNS设置",
		Func: func(c *ishell.Context) {
			data := []string{
				tools.UintToStr(setting.DNSPort()),
				setting.OutlandDNS(),
				setting.InlandDNS(),
				strings.Join(setting.BackupDNS(), ","),
			}
			format.ShowDNS(os.Stdout, data)
		},
	}
	dnsCmd.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "查看帮助",
		Func: func(c *ishell.Context) {
			c.Println(HelpDNS())
		},
	})
	dnsCmd.AddCmd(&ishell.Cmd{
		Name: "port",
		Help: "设置dns端口",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 && tools.IsNetPort(c.Args[0]) {
				setting.SetDNSPort(tools.StrToUint(c.Args[0]))
			}
		},
	})
	dnsCmd.AddCmd(&ishell.Cmd{
		Name: "inland",
		Help: "设置境内DNS",
		Func: func(c *ishell.Context) {
			setting.SetInlandDNS(c.Args[0])
		},
	})
	dnsCmd.AddCmd(&ishell.Cmd{
		Name: "outland",
		Help: "设置境外DNS",
		Func: func(c *ishell.Context) {
			setting.SetOutlandDNS(c.Args[0])
		},
	})
	dnsCmd.AddCmd(&ishell.Cmd{
		Name: "backup",
		Help: "设置备用DNS",
		Func: func(c *ishell.Context) {
			setting.SetBackupDNS(strings.Join(c.Args, ","))
		},
	})
	shell.AddCmd(dnsCmd)
}
