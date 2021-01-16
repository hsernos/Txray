package cmd

import (
	"Txray/core"
	"github.com/abiosoft/ishell"
)

var coreService core.Core

func Kill() {
	coreService.Stop()
}

func InitShell(shell *ishell.Shell) {
	coreService = core.New()
	shell.AddCmd(&ishell.Cmd{
		Name:    "version",
		Aliases: []string{"-v", "--version"},
		Help:    "程序版本",
		Func: func(c *ishell.Context) {
			c.Printf("%s version \"%s\"\n", c.Get("name"), c.Get("version"))
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Help:    "帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(Help())
		},
	})
	InitSettingShell(shell)
	InitNodeShell(shell)
	InitSubscribeShell(shell)
	InitDNSShell(shell)
	InitRouteShell(shell)
	InitServiceShell(shell)
}

// 参数解析
func FlagsParse(args []string, keys map[string]string) map[string]string {
	resultMap := make(map[string]string)
	key := "data"
	for _, x := range args {
		if len(x) >= 2 {
			if x[:2] == "--" {
				key = x[2:]
				resultMap[key] = ""
			} else if x[:1] == "-" {
				if x[1] >= 48 && x[1] <= 57 {
					resultMap[key] = x
				} else if len(x) == 2 {
					d, ok := keys[x[1:]]
					if ok {
						key = d
					} else {
						key = x[1:]
					}
					resultMap[key] = ""
				}
			} else {
				resultMap[key] = x
			}
		} else if len(x) > 0 {
			resultMap[key] = x
		}
	}
	return resultMap
}
