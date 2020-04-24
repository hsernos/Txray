package cmd

import (
	"github.com/abiosoft/ishell"
	"v3ray/config"
	"v3ray/tool"
)

var configObj config.Config

// InitConfig 初始化配置
func InitConfig() {
	configObj = config.NewConfig()
}

// Kill 关闭v2ray进程
func Kill() {
	configObj.Stop()
}

// InitShell 初始化
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
			title := []string{"监听端口", "udp转发", "启用流量监听", "多路复用", "允许局域网连接", "绕过局域网和大陆", "路由策略"}
			s := configObj.Settings
			data := []string{tool.UintToStr(s.Port),
				tool.BoolToStr(s.UDP),
				tool.BoolToStr(s.Sniffing),
				tool.BoolToStr(s.Mux),
				tool.BoolToStr(s.AllowLANConn),
				tool.BoolToStr(s.BypassLanAndContinent),
				s.DomainStrategy,
			}
			tool.GetTable(title, data)
		},
	})

	setting.AddCmd(&ishell.Cmd{
		Name: "alter",
		Help: "修改基础设置",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, map[string]string{
				"p": "port",
				"u": "udp",
				"s": "sniffing",
				"l": "lanconn",
				"m": "mux",
				"b": "bypass",
				"r": "route",
			})
			d, ok := r["port"]
			if ok && tool.IsUint(d) {
				configObj.SetPort(tool.StrToUint(d))
			}
			d, ok = r["udp"]
			if ok {
				if d == "true" {
					configObj.SetUDP(true)
				} else if d == "false" {
					configObj.SetUDP(false)
				}
			}
			d, ok = r["sniffing"]
			if ok {
				if d == "true" {
					configObj.SetSniffing(true)
				} else if d == "false" {
					configObj.SetSniffing(false)
				}
			}
			d, ok = r["lanconn"]
			if ok {
				if d == "true" {
					configObj.SetLANConn(true)
				} else if d == "false" {
					configObj.SetLANConn(false)
				}
			}
			d, ok = r["mux"]
			if ok {
				if d == "true" {
					configObj.SetMux(true)
				} else if d == "false" {
					configObj.SetMux(false)
				}
			}
			d, ok = r["bypass"]
			if ok {
				if d == "true" {
					configObj.SetBypassLanAndContinent(true)
				} else if d == "false" {
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
			l := []string{}
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
			})
			d1, ok1 := r["vmess"]
			d2, ok2 := r["file"]
			if ok1 {
				configObj.AddNodeByVmessLinks([]string{d1})
			} else if ok2 {
				configObj.AddNodeByFile(d2)
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
			Help: "查看节点",
			Func: func(c *ishell.Context) {
				r := FlagsParse(c.Args, nil)
				d, ok := r["data"]
				title := []string{"索引", "别名", "地址", "端口", "加密方式", "测试结果"}
				if ok {
					configObj.PingNodes(d)
					tool.GetTable(title, configObj.GetNodes("tcping")...)
					println("当前选定节点索引：", configObj.GetNodeIndex())
				} else {
					configObj.PingNodes("all")
					tool.GetTable(title, configObj.GetNodes("all")...)
					println("当前选定节点索引：", configObj.GetNodeIndex())
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
				title := []string{"索引", "别名", "地址", "端口", "加密方式", "测试结果"}
				if ok {
					tool.GetTable(title, configObj.GetNodes(d)...)
					println("当前选定节点索引：", configObj.GetNodeIndex())
				} else {
					tool.GetTable(title, configObj.GetNodes("all")...)
					println("当前选定节点索引：", configObj.GetNodeIndex())
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
				title := []string{"索引", "别名", "地址", "端口", "加密方式", "测试结果"}
				if ok {
					tool.GetTable(title, configObj.FindNodes(d)...)
					println("当前选定节点索引：", configObj.GetNodeIndex())
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
					if d == "true" {
						using = "true"
					} else if d == "false" {
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
			title := []string{"索引", "别名", "url", "是否启用"}
			if ok {
				tool.GetTable(title, configObj.GetSubs(d)...)
			} else {
				tool.GetTable(title, configObj.GetSubs("all")...)
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
			d, ok := r["proxy"]
			if ok && tool.IsUint(d) {
				configObj.AddNodeBySub(tool.StrToUint(d))
			} else {
				configObj.AddNodeBySub(1000000)
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
			title := []string{"索引", "DNS"}
			tool.GetTable(title, configObj.GetDNS()...)
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
		Name: "show-proxy-ip",
		Help: "查看代理IP规则",
		Func: func(c *ishell.Context) {
			title := []string{"索引", "规则"}
			tool.GetTable(title, configObj.GetProxyIP()...)
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-proxy-domain",
		Help: "查看代理Domain规则",
		Func: func(c *ishell.Context) {
			title := []string{"索引", "规则"}
			tool.GetTable(title, configObj.GetProxyDomain()...)
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-direct-ip",
		Help: "查看直连IP规则",
		Func: func(c *ishell.Context) {
			title := []string{"索引", "规则"}
			tool.GetTable(title, configObj.GetDirectIP()...)
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-direct-domain",
		Help: "查看直连Domain规则",
		Func: func(c *ishell.Context) {
			title := []string{"索引", "规则"}
			tool.GetTable(title, configObj.GetDirectDomain()...)
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-block-ip",
		Help: "查看禁止IP规则",
		Func: func(c *ishell.Context) {
			title := []string{"索引", "规则"}
			tool.GetTable(title, configObj.GetBlockIP()...)
		},
	})
	route.AddCmd(&ishell.Cmd{
		Name: "show-block-domain",
		Help: "查看禁止Domain规则",
		Func: func(c *ishell.Context) {
			title := []string{"索引", "规则"}
			tool.GetTable(title, configObj.GetBlockDomain()...)
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

	service := &ishell.Cmd{
		Name: "service",
		Help: "v2ray服务管理, 使用service查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(serviceHelp)
		},
	}

	service.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "开始服务",
		Func: func(c *ishell.Context) {
			r := FlagsParse(c.Args, nil)
			d, ok := r["data"]
			configObj.Stop()
			if ok && tool.IsInt(d) {
				configObj.Start(tool.StrToInt(d))
			} else {
				configObj.Start(-1)
			}
		},
	},
	)

	service.AddCmd(&ishell.Cmd{
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
	shell.AddCmd(service)
	shell.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "停止服务",
		Func: func(c *ishell.Context) {
			c.Println(help)
		},
	})
}
