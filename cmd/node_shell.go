package cmd

import (
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"os"
	"strings"
)

func InitNodeShell(shell *ishell.Shell) {

	node := &ishell.Cmd{
		Name: "node",
		Help: "节点管理, 使用node查看帮助信息",
		Func: func(c *ishell.Context) {
			c.Println(HelpNode())
		},
	}
	// 添加节点
	node.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "添加节点",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"l": "link",
				"c": "clipboard",
				"f": "file",
			})

			// 从剪贴板导入节点
			if _, ok := argMap["clipboard"]; ok {
				content, err := clipboard.ReadAll()
				if err != nil {
					log.Error(err)
					return
				}
				c.Println("剪贴板内容如下：")
				c.Println("========================================================================")
				c.Println(content)
				c.Println("========================================================================")
				if strings.Contains(content, "://") {
					links := make([]string, 0)
					for _, link := range strings.Split(content, "\n") {
						if protocols.IsSupportLinkFormat(link) {
							links = append(links, link)
						}
					}
					coreService.AddNodeByLinks(links...)
				} else {
					coreService.AddNodeBySubText(content)
				}
			} else if fileArg, ok := argMap["file"]; ok {
				if !tools.IsFile(fileArg) {
					log.Error("open ", fileArg, " : 没有这个文件")
					return
				}
				data, _ := ioutil.ReadFile(fileArg)
				content := strings.ReplaceAll(string(data), "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				c.Println("文件内容如下：")
				c.Println("========================================================================")
				c.Println(content)
				c.Println("========================================================================")
				if strings.Contains(content, "://") {
					links := make([]string, 0)
					for _, link := range strings.Split(content, "\n") {
						if protocols.IsSupportLinkFormat(link) {
							links = append(links, link)
						}
					}
					coreService.AddNodeByLinks(links...)
				} else {
					coreService.AddNodeBySubText(content)
				}
			} else if linkArg, ok := argMap["link"]; ok {
				if protocols.IsSupportLinkFormat(linkArg) {
					coreService.AddNodeByLinks(linkArg)
				} else {
					log.Error("格式错误, 应以 'ss://' 'vmess://' 'trojan://' 或 'vless://' 开头的链接")
					return
				}
			} else {
				c.ShowPrompt(false)
				modeList := []string{
					"VMess",
					"VLESS (暂不支持)",
					"Trojan",
					"Shadowsocks",
					"退出",
				}
				protocolMode := c.MultiChoice(modeList, "手动添加何种协议的节点?")
				c.ShowPrompt(true)

				switch protocolMode {
				// VMess
				case 0:
					c.Println("VMess")
					c.Println("========================")
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
					c.Println("========================")
					coreService.AddVMessNode(remarks, address, port, id, alterID,
						networkList[network], typeList[types], host, path, tlsList[tls])
					c.Println("添加完成")
				// VLESS
				case 1:
					c.Println("========================")
					c.Println("VLESS")
					c.Println("========================")
					log.Info("暂不支持VLESS协议")
					// TODO
				// Trojan
				case 2:
					c.Println("========================")
					c.Println("Trojan")
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					addr := c.ReadLine()
					c.Print("端口（port）: ")
					port := c.ReadLine()
					c.Print("密码（password）: ")
					password := c.ReadLine()
					c.Println("========================")
					coreService.AddTrojanNode(remarks, addr, port, password)
					c.Println("添加完成")
				// Shadowsocks
				case 3:
					c.Println("========================")
					c.Println("Shadowsocks")
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					addr := c.ReadLine()
					c.Print("端口（port）: ")
					port := c.ReadLine()
					c.Print("密码（password）: ")
					password := c.ReadLine()
					security := []string{
						"aes-256-cfb",
						"aes-128-cfb",
						"chacha20",
						"chacha20-ietf",
						"aes-256-gcm",
						"aes-256-gcm",
						"chacha20-poly1305",
						"chacha20-ietf-poly1305",
					}
					sIndex := c.MultiChoice(security, "加密方式（security） ?")
					c.Println("========================")
					coreService.AddShadowSocksNode(remarks, addr, port, password, security[sIndex])
					c.Println("添加完成")
				}
				c.Println()
			}

		},
	})
	// 查看节点
	node.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "查看节点",
		Func: func(c *ishell.Context) {
			var key string
			if len(c.Args) == 1 {
				key = c.Args[0]
			} else {
				key = "all"
			}
			format.ShowSimpleNode(os.Stdout, coreService.GetNodes(key)...)
			c.Println("当前选定节点索引:", coreService.GetNodeIndex())
		},
	})
	// 节点详情
	node.AddCmd(&ishell.Cmd{
		Name: "info",
		Help: "节点详情",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				if tools.IsInt(c.Args[0]) {
					index := tools.StrToInt(c.Args[0])
					if coreService.HasNode(index) {

						coreService.Show(index)

					} else {
						log.Warn(c.Args[0], ": 没有该索引")
					}
				} else {
					log.Warn(c.Args[0], ": 不是一个数字")
				}
			} else if len(c.Args) == 0 {
				log.Warn("需要输入一个索引")
			} else {
				log.Warn(strings.Join(c.Args, " "), ": 参数过多")
			}
		},
	})
	// 删除节点
	node.AddCmd(&ishell.Cmd{
		Name: "del",
		Help: "删除节点",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				coreService.DelNodes(key)
			} else if len(c.Args) == 0 {
				log.Warn("还需要输入一个索引")
			} else {
				log.Warn(strings.Join(c.Args, " "), ": 参数过多")
			}
		},
	})
	// 导出节点
	node.AddCmd(&ishell.Cmd{
		Name: "export",
		Help: "导出节点",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"c": "clipboard",
			})
			key := "all"
			if _, ok := argMap["data"]; ok {
				key, _ = argMap["data"]
			}
			links := coreService.ExportNodes(key)
			if _, ok := argMap["clipboard"]; ok {
				err := clipboard.WriteAll(strings.Join(links, "\n"))
				if err != nil {
					log.Error(err)
					return
				}
				c.Println("共导出", len(links), "条至剪贴板")
			} else {
				c.Println("========================================================================")
				c.Println(strings.Join(links, "\n"))
				c.Println("========================================================================")
				c.Println("共导出", len(links), "条")
			}
		},
	})
	// 测试节点tcp延迟
	node.AddCmd(&ishell.Cmd{
		Name: "tcping",
		Help: "测试节点tcp延迟",
		Func: func(c *ishell.Context) {
			var key string
			if len(c.Args) == 1 {
				key = c.Args[0]
			} else {
				key = "all"
			}
			coreService.PingNodes(key)
			err := shell.Process("node", "show", "t")
			if err != nil {
				log.Error(err)
			}
		},
	})
	// 查找节点
	node.AddCmd(&ishell.Cmd{
		Name: "find",
		Help: "查找节点",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				key := c.Args[0]
				coreService.DelNodes(key)
				format.ShowSimpleNode(os.Stdout, coreService.FindNodes(key)...)
				c.Println("当前选定节点索引:", coreService.GetNodeIndex())
			} else if len(c.Args) == 0 {
				log.Warn("还需要输入一个查找关键字")
			} else {
				log.Warn(strings.Join(c.Args, " "), ": 参数过多")
			}

		},
	})
	shell.AddCmd(node)
}
