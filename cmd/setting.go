package cmd

import (
	"Txray/cmd/help"
	"Txray/core/setting"
	"Txray/log"
	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
)

func InitSettingShell(shell *ishell.Shell) {
	baseSettingCmd := &ishell.Cmd{
		Name: "setting",
		Func: func(c *ishell.Context) {
			// 连接设置
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"socks端口", "http端口", "udp转发", "流量地址监听", "允许来自局域网连接", "多路复用"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			data := []string{
				strconv.Itoa(setting.Socks()),
				strconv.Itoa(setting.Http()),
				strconv.FormatBool(setting.UDP()),
				strconv.FormatBool(setting.Sniffing()),
				strconv.FormatBool(setting.FromLanConn()),
				strconv.FormatBool(setting.Mux()),
			}
			table.Append(data)
			table.Render()

			// DNS及路由设置
			table = tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"DNS端口", "国外DNS", "国内DNS", "备用国内DNS", "路由策略", "绕过局域网和大陆"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			data = []string{
				strconv.Itoa(setting.DNSPort()),
				setting.DNSForeign(),
				setting.DNSDomestic(),
				setting.DNSBackup(),
				setting.RoutingStrategy(),
				strconv.FormatBool(setting.RoutingBypass()),
			}
			table.Append(data)
			table.Render()

			table = tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"测试国外URL", "测试超时时间 (秒)", "批量测试终止时间 (毫秒)", "启动时执行"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			data = []string{
				setting.TestUrl(),
				strconv.Itoa(setting.TestTimeout()),
				strconv.Itoa(setting.TestMinTime()),
				setting.RunBefore(),
			}
			table.Append(data)
			table.Render()
		},
	}
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Func: func(c *ishell.Context) {
			c.Println(help.Setting)
		},
	})

	// 本地连接设置
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "socks",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				v, err := strconv.Atoi(c.Args[0])
				if err != nil {
					log.Warn("非法输入")
					return
				}
				err = setting.SetSocks(v)
				if err != nil {
					log.Error(err)
					return
				}
			}
			log.Info("socks端口: ", setting.Socks())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "http",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				v, err := strconv.Atoi(c.Args[0])
				if err != nil {
					log.Error("非法输入")
					return
				}
				err = setting.SetHttp(v)
				if err != nil {
					log.Error(err)
					return
				}
			}
			log.Info("http端口: ", setting.Http())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "udp",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				str := strings.ToLower(c.Args[0])
				switch str {
				case "y", "yes", "true", "t":
					setting.SetUDP(true)
				case "n", "no", "false", "f":
					setting.SetUDP(false)
				}
			}
			log.Info("UDP转发: ", setting.UDP())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "sniffing",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				str := strings.ToLower(c.Args[0])
				switch str {
				case "y", "yes", "true", "t":
					setting.SetSniffing(true)
				case "n", "no", "false", "f":
					setting.SetSniffing(false)
				}
			}
			log.Info("流量地址监听: ", setting.Sniffing())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "mux",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				str := strings.ToLower(c.Args[0])
				switch str {
				case "y", "yes", "true", "t":
					setting.SetMux(true)
				case "n", "no", "false", "f":
					setting.SetMux(false)
				}
			}
			log.Info("多路复用: ", setting.Mux())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "from_lan_conn",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				str := strings.ToLower(c.Args[0])
				switch str {
				case "y", "yes", "true", "t":
					setting.SetFromLanConn(true)
				case "n", "no", "false", "f":
					setting.SetFromLanConn(false)
				}
			}
			log.Info("来自局域网连接: ", setting.FromLanConn())
		},
	})

	// 路由
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "routing.strategy",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				switch c.Args[0] {
				case "1", "AsIs":
					setting.SetRoutingStrategy(1)
				case "2", "RoutingStrategy":
					setting.SetRoutingStrategy(2)
				case "3", "IPOnDemand":
					setting.SetRoutingStrategy(3)
				}
			}
			log.Info("路由策略: ", setting.RoutingStrategy())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "routing.bypass",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				str := strings.ToLower(c.Args[0])
				switch str {
				case "y", "yes", "true", "t":
					setting.SetRoutingBypass(true)
				case "n", "no", "false", "f":
					setting.SetRoutingBypass(false)
				}
			}
			log.Info("绕过局域网和大陆: ", setting.RoutingBypass())
		},
	})

	// DNS
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "dns.port",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				v, err := strconv.Atoi(c.Args[0])
				if err != nil {
					log.Error("非法输入")
					return
				}
				err = setting.SetDNSPort(v)
				if err != nil {
					log.Error(err)
					return
				}
			}
			log.Info("DNS端口: ", setting.DNSPort())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "dns.foreign",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				err := setting.SetDNSForeign(c.Args[0])
				if err != nil {
					log.Warn(err)
					return
				}
			}
			log.Info("国外DNS: ", setting.DNSForeign())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "dns.domestic",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				err := setting.SetDNSDomestic(c.Args[0])
				if err != nil {
					log.Warn(err)
					return
				}
			}
			log.Info("国内DNS: ", setting.DNSDomestic())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "dns.backup",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				err := setting.SetDNSBackup(c.Args[0])
				if err != nil {
					log.Warn(err)
					return
				}
			}
			log.Info("备用国内DNS: ", setting.DNSBackup())
		},
	})

	// 外网测试设置
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "test.timeout",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				v, err := strconv.Atoi(c.Args[0])
				if err != nil {
					log.Error("非法输入")
					return
				}
				err = setting.SetTestTimeout(v)
				if err != nil {
					log.Error(err)
					return
				}
			}
			log.Info("外网测试超时时间 (秒): ", setting.TestTimeout())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "test.mintime",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				v, err := strconv.Atoi(c.Args[0])
				if err != nil {
					log.Error("非法输入")
					return
				}
				err = setting.SetTestMinTime(v)
				if err != nil {
					log.Error(err)
					return
				}
			}
			log.Info("批量测试终止时间 (毫秒): ", setting.TestMinTime())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "test.url",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				err := setting.SetTestUrl(c.Args[0])
				if err != nil {
					log.Warn(err)
					return
				}
			}
			log.Info("外网测试URL: ", setting.TestUrl())
		},
	})
	baseSettingCmd.AddCmd(&ishell.Cmd{
		Name: "run_before",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"c": "close",
			})
			if _, ok := argMap["close"]; ok {
				err := setting.SetRunBefore("")
				if err != nil {
					log.Warn(err)
					return
				}
			} else if _, ok := argMap["data"]; ok {
				err := setting.SetRunBefore(argMap["data"])
				if err != nil {
					log.Warn(err)
					return
				}
			}
			log.Info("启动时执行: ", setting.RunBefore())
		},
	})
	shell.AddCmd(baseSettingCmd)
}
