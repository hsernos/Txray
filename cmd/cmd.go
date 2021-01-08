package cmd

import (
	"Tv2ray/config"
	log "Tv2ray/logger"
	"Tv2ray/tools"
	"fmt"
	"github.com/abiosoft/ishell"
	"os"
	"strings"
)

var configObj config.Config

// 初始化配置
func InitConfig() {
	configObj = config.NewConfig()
}

// 关闭v2ray进程
func Kill() {
	configObj.Stop()
}

// 初始化
func InitShell(shell *ishell.Shell) {
	setting := &ishell.Cmd{
		Name: "setting",
		Help: "基础设置, 使用setting查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(settingHelp)
		},
	}
	setting.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看基础设置",
		Func: func(c *ishell.Context) {
			s := configObj.Settings
			data := []string{tools.UintToStr(s.Port),
				tools.UintToStr(s.Http),
				tools.BoolToStr(s.UDP),
				tools.BoolToStr(s.Sniffing),
				tools.BoolToStr(s.Mux),
				tools.BoolToStr(s.AllowLANConn),
				tools.BoolToStr(s.BypassLanAndContinent),
				s.DomainStrategy,
			}
			ShowSetting(os.Stdout, data)
		},
	})

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
			d, ok := r["port"]
			if ok && tools.IsUint(d) {
				configObj.SetPort(tools.StrToUint(d))
			}
			d, ok = r["http"]
			if ok && tools.IsUint(d) {
				configObj.SetHttpPort(tools.StrToUint(d))
			}
			d, ok = r["udp"]
			if ok {
				if d == "y" {
					configObj.SetUDP(true)
				} else if d == "n" {
					configObj.SetUDP(false)
				}
			}
			d, ok = r["sniffing"]
			if ok {
				if d == "y" {
					configObj.SetSniffing(true)
				} else if d == "n" {
					configObj.SetSniffing(false)
				}
			}
			d, ok = r["lanconn"]
			if ok {
				if d == "y" {
					configObj.SetLANConn(true)
				} else if d == "n" {
					configObj.SetLANConn(false)
				}
			}
			d, ok = r["mux"]
			if ok {
				if d == "y" {
					configObj.SetMux(true)
				} else if d == "n" {
					configObj.SetMux(false)
				}
			}
			d, ok = r["bypass"]
			if ok {
				if d == "y" {
					configObj.SetBypassLanAndContinent(true)
				} else if d == "n" {
					configObj.SetBypassLanAndContinent(false)
				}
			}
			d, ok = r["route"]
			if ok {
				if d == "1" {
					configObj.SetDomainStrategy(1)
				} else if d == "2" {
					configObj.SetDomainStrategy(2)
				} else if d == "3" {
					configObj.SetDomainStrategy(3)
				}
			}
		},
	})

	node := &ishell.Cmd{
		Name: "node",
		Help: "节点管理, 使用node查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(nodeHelp)
		},
	}

	node.AddCmd(&ishell.Cmd{
		Name: "del",
		Help: "删除节点",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelNodes(d)
			} else {
				configObj.DelNodes("")
			}
		},
	})
	node.AddCmd(&ishell.Cmd{
		Name: "export",
		Help: "导出节点",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			var l []string
			if ok {
				l = configObj.ExportNodes(d)
			} else {
				l = configObj.ExportNodes("")
			}
			for _, x := range l {
				c.Println(x)
				c.Println()
			}
		},
	})

	node.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "添加节点",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"v": "vmess",
				"f": "file",
				"s": "subfile",
			})
			vmess, ok1 := r["vmess"]
			file, ok2 := r["file"]
			subfile, ok3 := r["subfile"]
			if ok1 {
				configObj.AddNodeByVmessLinks([]string{vmess})
			} else if ok2 {
				configObj.AddNodeByFile(file)
			} else if ok3 {
				configObj.AddNodeBySubFile(subfile)
			} else {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)
				c.Print("别名（remarks）: ")
				remarks := c.ReadLine()
				c.Print("地址（address）: ")
				address := c.ReadLine()
				c.Print("端口（port）: ")
				port := c.ReadLine()
				c.Print("用户ID（id）: ")
				id := c.ReadLine()
				c.Print("额外ID（alterID）: ")
				alterID := c.ReadLine()
				securityList := []string{
					"chacha20-poly1305",
					"aes-128-gcm",
					"auto",
					"none",
				}
				security := c.MultiChoice(securityList, "加密方式（security） ?")
				networkList := []string{
					"tcp",
					"kcp",
					"ws",
					"h2",
					"quic",
				}
				network := c.MultiChoice(networkList, "传输协议（network） ?")
				typeList := []string{
					"none",
					"http",
					"srtp",
					"utp",
					"wechat-video",
					"dtls",
					"wireguard",
				}
				types := c.MultiChoice(typeList, "伪装类型（type） ?")
				c.Print("伪装域名 host（host/ws host/h2 host）/QUIC 加密方式: ")
				host := c.ReadLine()
				c.Print("path（ws path/h2 path)/QUIC 加密秘钥: ")
				path := c.ReadLine()
				tlsList := []string{
					"tls",
					"",
				}
				tls := c.MultiChoice(tlsList, "底层安全传输 ?")
				configObj.AddNode(remarks, address, port, id, alterID, securityList[security], networkList[network], typeList[types], host, path, tlsList[tls])
			}
		},
	})

	node.AddCmd(
		&ishell.Cmd{
			Name: "tcping",
			Help: "tcping节点",
			Func: func(c *ishell.Context) {
				r := FlagsParse(c.Args, nil)
				d, ok := r["data"]
				if ok {
					configObj.PingNodes(d)
				} else {
					configObj.PingNodes("all")
				}
				ShowSimpleNode(os.Stdout, configObj.GetNodes("test")...)
				c.Println("当前选定节点索引：", configObj.GetNodeIndex())
			},
		},
	)

	node.AddCmd(
		&ishell.Cmd{
			Name: "info",
			Help: "查看节点的详细信息",
			Func: func(c *ishell.Context) {
				r := FlagsParse(c.Args, nil)
				d, ok := r["data"]
				if ok {
					if tools.IsUint(d) && tools.StrToUint(d) > 0 {

						n := configObj.GetNode(tools.StrToUint(d) - 1)
						if n == nil {
							log.Warn(d, "超出了索引范围")
						} else {
							c.Println()
							c.Printf("%7s: %s\n", "索引", d)
							c.Printf("%7s: %s\n", "别名", n.Remarks)
							c.Printf("%7s: %s\n", "地址", n.Address)
							c.Printf("%7s: %d\n", "端口", n.Port)
							c.Printf("%7s: %s\n", "用户ID", n.ID)
							c.Printf("%7s: %d\n", "额外ID", n.AlterID)
							c.Printf("%5s: %s\n", "加密方式", n.Security)
							c.Printf("%5s: %s\n", "伪装类型", n.HeaderType)
							c.Printf("%5s: %s\n", "伪装域名", n.RequestHost)
							c.Printf("%9s: %s\n", "path", n.Path)
							c.Printf("%5s: %s\n", "安全传输", n.StreamSecurity)
							c.Println()
						}
					} else {
						log.Warn(d, "没有此索引对应的节点")
					}
				} else {
					log.Info("还需要输入一个索引")
				}

			},
		},
	)

	node.AddCmd(
		&ishell.Cmd{
			Name: "show",
			Help: "查看节点",
			Func: func(c *ishell.Context) {
				r := FlagsParse(c.Args, nil)
				d, ok := r["data"]
				if ok {
					ShowSimpleNode(os.Stdout, configObj.GetNodes(d)...)
					c.Println("当前选定节点索引：", configObj.GetNodeIndex())
				} else {
					ShowSimpleNode(os.Stdout, configObj.GetNodes("all")...)
					c.Println("当前选定节点索引：", configObj.GetNodeIndex())
				}
			},
		},
	)

	node.AddCmd(
		&ishell.Cmd{
			Name: "find",
			Help: "查找节点",
			Func: func(c *ishell.Context) {
				r := FlagsParse(c.Args, nil)
				d, ok := r["data"]
				if ok {
					ShowSimpleNode(os.Stdout, configObj.FindNodes(d)...)
					c.Println("当前选定节点索引：", configObj.GetNodeIndex())
				}
			},
		},
	)

	sub := &ishell.Cmd{
		Name: "sub",
		Help: "订阅管理, 使用sub查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(subHelp)
		},
	}

	sub.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "添加订阅",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"r": "remarks",
			})
			url, ok := r["data"]
			if ok {
				d, ok := r["remarks"]
				if ok {
					configObj.AddSub(url, d)
				} else {
					configObj.AddSub(url, "remarks")
				}
			}
		},
	},
	)

	sub.AddCmd(&ishell.Cmd{
		Name: "del",
		Help: "删除订阅",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelSubs(d)
			} else {
				configObj.DelSubs("")
			}
		},
	},
	)
	sub.AddCmd(&ishell.Cmd{
		Name: "atler",
		Help: "修改订阅",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"r": "remarks",
				"u": "url",
			})
			key, ok := r["data"]
			if ok {
				using := ""
				url := ""
				remarks := ""
				d, ok := r["remarks"]
				if ok {
					remarks = d
				}
				d, ok = r["url"]
				if ok {
					url = d
				}
				d, ok = r["using"]
				if ok {
					if d == "y" {
						using = "true"
					} else if d == "n" {
						using = "false"
					}
				}
				configObj.SetSubs(key, using, url, remarks)
			}
		},
	},
	)
	sub.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看订阅",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				ShowSub(os.Stdout, configObj.GetSubs(d)...)
			} else {
				ShowSub(os.Stdout, configObj.GetSubs("all")...)
			}
		},
	},
	)

	sub.AddCmd(&ishell.Cmd{
		Name: "update-node",
		Help: "更新节点",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"p": "proxy",
			})
			key, _ := r["data"]
			proxyPort, isProxy := r["proxy"]
			if isProxy {
				if tools.IsUint(proxyPort) {
					configObj.AddNodeBySub(key, tools.StrToUint(proxyPort))
				} else {
					configObj.AddNodeBySub(key, configObj.Settings.Port)
				}
			} else {
				configObj.AddNodeBySub(key, 1000000)
			}

		},
	},
	)

	dns := &ishell.Cmd{
		Name: "dns",
		Help: "DNS管理, 使用dns查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(dnsHelp)
		},
	}
	dns.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "添加DNS",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddDNS(data)
			}
		},
	})

	dns.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看DNS",
		Func: func(c *ishell.Context) {
			ShowDNS(os.Stdout, configObj.GetDNS()...)
		},
	})

	dns.AddCmd(&ishell.Cmd{
		Name: "del",
		Help: "删除DNS",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelDNS(d)
			} else {
				configObj.DelDNS("")
			}
		},
	})

	route := &ishell.Cmd{
		Name: "route",
		Help: "路由管理, 使用route查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(routeHelp)
		},
	}

	route.AddCmd(&ishell.Cmd{
		Name: "add-block-ip",
		Help: "添加禁止规则",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddBlockIP(data)
			}
		},
	})

	route.AddCmd(&ishell.Cmd{
		Name: "add-block-domain",
		Help: "添加禁止规则",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddBlockDomain(data)
			}
		},
	})

	route.AddCmd(&ishell.Cmd{
		Name: "add-direct-ip",
		Help: "添加直连规则",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddDirectIP(data)
			}
		},
	})

	route.AddCmd(&ishell.Cmd{
		Name: "add-direct-domain",
		Help: "添加直连规则",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddDirectDomain(data)
			}
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "add-proxy-ip",
		Help: "添加代理规则",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddProxyIP(data)
			}
		},
	})

	route.AddCmd(&ishell.Cmd{
		Name: "add-proxy-domain",
		Help: "添加代理规则",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			data, ok := r["data"]
			if ok {
				configObj.AddProxyDomain(data)
			}
		},
	})

	route.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看全部规则",
		Func: func(c *ishell.Context) {
			c.Println()
			c.Println("代理路由规则：")
			c.Println("==============================================")
			ShowIPRouter(os.Stdout, configObj.GetProxyIP()...)
			c.Println()
			ShowDomainRouter(os.Stdout, configObj.GetProxyDomain()...)
			c.Println("==============================================\n")

			c.Println("直连路由规则")
			c.Println("==============================================")
			ShowIPRouter(os.Stdout, configObj.GetDirectIP()...)
			c.Println()
			ShowDomainRouter(os.Stdout, configObj.GetDirectDomain()...)
			c.Println("==============================================\n")

			c.Println("禁止路由规则")
			c.Println("==============================================")
			ShowIPRouter(os.Stdout, configObj.GetBlockIP()...)
			c.Println()
			ShowDomainRouter(os.Stdout, configObj.GetBlockDomain()...)
			c.Println("==============================================\n")
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-proxy",
		Help: "查看代理规则",
		Func: func(c *ishell.Context) {
			c.Println()
			c.Println("代理路由规则：")
			c.Println("==============================================")
			ShowIPRouter(os.Stdout, configObj.GetProxyIP()...)
			c.Println()
			ShowDomainRouter(os.Stdout, configObj.GetProxyDomain()...)
			c.Println("==============================================\n")
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-direct",
		Help: "查看直连规则",
		Func: func(c *ishell.Context) {
			c.Println()
			c.Println("直连路由规则")
			c.Println("==============================================")
			ShowIPRouter(os.Stdout, configObj.GetDirectIP()...)
			c.Println()
			ShowDomainRouter(os.Stdout, configObj.GetDirectDomain()...)
			c.Println("==============================================\n")
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-block",
		Help: "查看禁止规则",
		Func: func(c *ishell.Context) {
			c.Println()
			c.Println("禁止路由规则")
			c.Println("==============================================")
			ShowIPRouter(os.Stdout, configObj.GetBlockIP()...)
			c.Println()
			ShowDomainRouter(os.Stdout, configObj.GetBlockDomain()...)
			c.Println("==============================================\n")
		},
	})

	route.AddCmd(&ishell.Cmd{
		Name: "del-proxy-ip",
		Help: "删除",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelDirectIP(d)
			} else {
				configObj.DelDirectIP("")
			}
		},
	},
	)

	route.AddCmd(&ishell.Cmd{
		Name: "del-proxy-domain",
		Help: "删除",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelProxyDomain(d)
			} else {
				configObj.DelProxyDomain("")
			}
		},
	},
	)

	route.AddCmd(&ishell.Cmd{
		Name: "del-direct-ip",
		Help: "删除",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelDirectIP(d)
			} else {
				configObj.DelDirectIP("")
			}
		},
	},
	)

	route.AddCmd(&ishell.Cmd{
		Name: "del-direct-domain",
		Help: "删除",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelDirectDomain(d)
			} else {
				configObj.DelDirectDomain("")
			}
		},
	},
	)
	route.AddCmd(&ishell.Cmd{
		Name: "del-block-ip",
		Help: "删除",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelBlockIP(d)
			} else {
				configObj.DelBlockIP("")
			}
		},
	},
	)

	route.AddCmd(&ishell.Cmd{
		Name: "del-block-domain",
		Help: "删除",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			if ok {
				configObj.DelBlockDomain(d)
			} else {
				configObj.DelBlockDomain("")
			}
		},
	},
	)

	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Help: "开始服务",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			configObj.Stop()
			if ok && tools.IsUint(d) {
				if configObj.Start(tools.StrToInt(d) - 1) {
					t, status := configObj.TestNode("https://www.youtube.com")
					log.Info(status, " [ https://www.youtube.com ] 延迟：", strings.Trim(t, " "))
				}

			} else if ok && (len(tools.IndexDeal(d, len(configObj.Nodes))) >= 1 || (len(d) > 1 && d[0] == 't' && tools.IsUint(d[1:]))) {
				fmt.Println(d)
				min := 100000
				index := -1
				var indexs []int
				if d[0] == 't' {
					num := tools.StrToInt(d[1:])
					if num != 0 {
						indexs = configObj.GetNodeTestSort(num)
					}
				} else {
					l := len(configObj.Nodes)
					indexs = tools.IndexDeal(d, l)
				}

				log.Info("=====================================================================")
				for _, x := range indexs {
					if configObj.Start(x) {
						t, status := configObj.TestNode("https://www.youtube.com")
						t = strings.Trim(t, " ")
						if t != "-1ms" && min > tools.StrToInt(strings.TrimRight(t, "ms")) {
							index = x
							min = tools.StrToInt(strings.TrimRight(t, "ms"))
						}
						log.Info(status, " [ https://www.youtube.com ] 延迟：", t)
					}
					configObj.Stop()
				}
				log.Info("=====================================================================")
				if index != -1 {
					log.Info("延迟最小的节点索引为：", index+1, "，延迟：", min, "ms")
					configObj.Start(index)
				} else {
					log.Info("所选节点全部不能访问外网")
				}
			} else {
				if configObj.Start(-1) {
					t, status := configObj.TestNode("https://www.youtube.com")
					log.Info(status, " [ https://www.youtube.com ] 延迟：", strings.Trim(t, " "))
				}
			}
		},
	},
	)

	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "停止服务",
		Func: func(c *ishell.Context) {
			configObj.Stop()
		},
	},
	)

	shell.AddCmd(setting)
	shell.AddCmd(node)
	shell.AddCmd(sub)
	shell.AddCmd(dns)
	shell.AddCmd(route)
	shell.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "帮助",
		Func: func(c *ishell.Context) {
			c.Println(help)
		},
	})
}
