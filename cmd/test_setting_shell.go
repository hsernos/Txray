package cmd

import (
	"Txray/core/setting"
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"os"
)

func InitTestShell(shell *ishell.Shell) {
	testCmd := &ishell.Cmd{
		Name: "test",
		Help: "查看test设置",
		Func: func(c *ishell.Context) {
			test := setting.TestSetting()
			data := []string{
				test.Url,
				tools.UintToStr(test.TimeOut),
			}
			format.ShowTest(os.Stdout, data)
		},
	}
	testCmd.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "查看帮助",
		Func: func(c *ishell.Context) {
			c.Println(HelpTest())
		},
	})
	testCmd.AddCmd(&ishell.Cmd{
		Name: "url",
		Help: "设置测试网站",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				setting.SetTestUrl(c.Args[0])
			}
		},
	})
	testCmd.AddCmd(&ishell.Cmd{
		Name: "timeout",
		Help: "设置超时时间",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 && tools.IsUint(c.Args[0]) {
				setting.SetTimeOut(tools.StrToUint(c.Args[0]))
			}
		},
	})
	shell.AddCmd(testCmd)
}
