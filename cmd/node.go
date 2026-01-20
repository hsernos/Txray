// cmd/node.go 负责 shell 层面节点管理命令的注册与实现
package cmd

import (
	"Txray/cmd/help"       // 帮助文档内容
	"Txray/core"           // 索引工具
	"Txray/core/manage"    // 节点管理器
	"Txray/core/node"      // 节点结构体
	"Txray/core/protocols" // 协议定义
	"Txray/core/sub"       // 订阅相关
	"Txray/log"            // 日志
	"fmt"                  // 格式化输出
	"net/url"              // URL 解析
	"os"                   // 系统操作
	"strconv"              // 字符串与数字转换
	"strings"              // 字符串处理

	"github.com/abiosoft/ishell"        // shell 框架
	"github.com/atotto/clipboard"       // 剪贴板操作
	"github.com/olekukonko/tablewriter" // 表格输出
)

// InitNodeShell 注册 node 命令及其子命令，支持节点展示、添加、删除、测试、帮助等
func InitNodeShell(shell *ishell.Shell) {
	nodeCmd := &ishell.Cmd{
		Name: "node",
		Func: func(c *ishell.Context) {
			// 解析参数，支持 -d 降序
			argMap := FlagsParse(c.Args, map[string]string{
				"d": "desc",
			})
			key := "all"
			if _, ok := argMap["data"]; ok {
				key = argMap["data"]
			}
			_, isDesc := argMap["desc"]
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"索引", "协议", "别名", "地址", "端口", "测试结果"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			center := tablewriter.ALIGN_CENTER
			left := tablewriter.ALIGN_LEFT
			table.SetColumnAlignment([]int{center, center, left, center, center, center})
			table.SetColWidth(70)
			indexList := core.IndexList(key, manage.Manager.NodeLen())
			if isDesc {
				indexList = core.Reverse(indexList)
			}
			for _, index := range indexList {
				n := manage.Manager.GetNode(index)
				if n != nil {
					table.Append([]string{
						strconv.Itoa(index),
						string(n.GetProtocolMode()),
						n.GetName(),
						n.GetAddr(),
						strconv.Itoa(n.GetPort()),
						n.TestResultStr(),
					})
				}
			}
			if n := manage.Manager.SelectedNode(); n == nil {
				table.SetCaption(true, fmt.Sprintf("[ %d/%d ] %s",
					manage.Manager.SelectedIndex(),
					manage.Manager.NodeLen(),
					"无节点",
				))
			} else {
				table.SetCaption(true, fmt.Sprintf("[ %d/%d ] %s",
					manage.Manager.SelectedIndex(),
					manage.Manager.NodeLen(),
					n.GetName(),
				))
			}
			table.Render()
		},
	}
	// help
	nodeCmd.AddCmd(&ishell.Cmd{
		Name:    "help",
		Aliases: []string{"-h", "--help"},
		Help:    "查看帮助",
		Func: func(c *ishell.Context) {
			c.Println(help.Node)
		},
	})
	// tcping
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "tcping",
		Func: func(c *ishell.Context) {
			manage.Manager.Tcping()
			_ = shell.Process("node", "-d")
		},
	})
	// info
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "info",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				v, err := strconv.Atoi(c.Args[0])
				if err != nil {
					log.Warn("非法参数")
					return
				}
				n := manage.Manager.GetNode(v)
				if n == nil {
					log.Warn("不存在该节点")
					return
				}
				n.Show()
			}
		},
	})
	// rm
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "rm",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				manage.Manager.DelNode(c.Args[0])
			}
		},
	})
	// sort
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "sort",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				var mode int
				switch c.Args[0] {
				case "0":
					mode = 0
				case "1":
					mode = 1
				case "2":
					mode = 2
				case "3":
					mode = 3
				case "4":
					mode = 4
				case "5":
					mode = 5
				default:
					return
				}
				manage.Manager.Sort(mode)
				_ = shell.Process("node")
			}
		},
	})
	// export
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "export",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"c": "clipboard",
			})
			key := "all"
			if _, ok := argMap["data"]; ok {
				key = argMap["data"]
			}
			links := manage.Manager.GetNodeLink(key)
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
	// find
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "find",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				return
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"索引", "协议", "别名", "地址", "端口", "测试结果"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			center := tablewriter.ALIGN_CENTER
			left := tablewriter.ALIGN_LEFT
			table.SetColumnAlignment([]int{center, center, left, center, center, center})
			table.SetColWidth(70)
			manage.Manager.NodeForEach(func(i int, n *node.Node) {
				if n != nil && strings.Contains(n.GetName(), c.Args[0]) {
					defer table.Append([]string{
						strconv.Itoa(i),
						string(n.GetProtocolMode()),
						n.GetName(),
						n.GetAddr(),
						strconv.Itoa(n.GetPort()),
						n.TestResultStr(),
					})
				}
			})
			if n := manage.Manager.SelectedNode(); n == nil {
				table.SetCaption(true, fmt.Sprintf("[ %d/%d ] %s",
					manage.Manager.SelectedIndex(),
					manage.Manager.NodeLen(),
					"无节点",
				))
			} else {
				table.SetCaption(true, fmt.Sprintf("[ %d/%d ] %s",
					manage.Manager.SelectedIndex(),
					manage.Manager.NodeLen(),
					n.GetName(),
				))
			}
			table.Render()
		},
	})
	// add
	nodeCmd.AddCmd(&ishell.Cmd{
		Name: "add",
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
				content = strings.ReplaceAll(content, "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				c.Println("剪贴板内容如下：")
				c.Println("========================================================================")
				c.Println(content)
				c.Println("========================================================================")
				if strings.Contains(content, "://") {
					for _, link := range strings.Split(content, "\n") {
						manage.Manager.AddNode(node.NewNode(link, ""))
					}
				} else {
					for _, link := range sub.Sub2links(content) {
						manage.Manager.AddNode(node.NewNode(link, ""))
					}
				}
			} else if fileArg, ok := argMap["file"]; ok {
				if _, err := os.Stat(fileArg); os.IsNotExist(err) {
					log.Error("open ", fileArg, " : 没有这个文件")
					return
				}
				data, _ := os.ReadFile(fileArg)
				content := strings.ReplaceAll(string(data), "\r\n", "\n")
				content = strings.ReplaceAll(content, "\r", "\n")
				c.Println("文件内容如下：")
				c.Println("========================================================================")
				c.Println(content)
				c.Println("========================================================================")
				if strings.Contains(content, "://") {
					for _, link := range strings.Split(content, "\n") {
						manage.Manager.AddNode(node.NewNode(link, ""))
					}
				} else {
					for _, link := range sub.Sub2links(content) {
						manage.Manager.AddNode(node.NewNode(link, ""))
					}
				}
			} else if linkArg, ok := argMap["link"]; ok {
				manage.Manager.AddNode(node.NewNode(linkArg, ""))
			} else {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)
				modeList := []string{
					protocols.ModeVMess.String(),
					protocols.ModeVLESS.String(),
					protocols.ModeVMessAEAD.String(),
					protocols.ModeTrojan.String(),
					protocols.ModeShadowSocks.String(),
					protocols.ModeSocks.String(),
					"退出",
				}
				i := c.MultiChoice(modeList, "手动添加何种协议的节点?")
				protocolMode := modeList[i]
				switch protocolMode {
				case protocols.ModeVMessAEAD.String():
					c.Println("========================")
					c.Println(protocolMode)
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					address := c.ReadLine()
					c.Print("端口（port）: ")
					port, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("端口为数字，且取值为1~65535")
						return
					}
					c.Print("用户ID（id）: ")
					id := c.ReadLine()
					data := make(map[string]string)
					networkList := []string{
						"tcp",
						"kcp",
						"ws",
						"xhttp",
						"h2",
						"quic",
						"grpc",
					}
					index := c.MultiChoice(networkList, "传输协议（network）?")
					network := networkList[index]
					if network != "tcp" {
						data["type"] = network
					}
					switch network {
					case "kcp":
						typeList := []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
						index := c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[index] != "none" {
							data["headerType"] = typeList[index]
						}
						c.Print("KCP种子（seed）: ")
						seed := c.ReadLine()
						if seed != "" {
							data["seed"] = seed
						}
					case "ws":
						c.Print("WebSocket 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("WebSocket Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}

					case "xhttp":
						typeList := []string{
							"auto",
							"packet-up",
							"stream-up",
							"stream-one",
						}
						index := c.MultiChoice(typeList, "xhttp 的模式（mode）?")
						mode := typeList[index]
						data["mode"] = mode

						c.Print("xhttp 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("xhttp Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}

					case "h2":
						c.Print("HTTP/2 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("HTTP/2 Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}
					case "quic":
						typeList := []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
						index := c.MultiChoice(typeList, "伪装类型（headerType）?")
						data["headerType"] = typeList[index]
						quicSecurityList := []string{
							"none",
							"aes-128-gcm",
							"chacha20-poly1305",
						}
						index = c.MultiChoice(quicSecurityList, "QUIC 的加密方式（quicSecurity）?")
						quicSecurity := quicSecurityList[index]
						data["quicSecurity"] = quicSecurity
						if quicSecurity != "none" {
							c.Print("当 QUIC 的加密方式不为 none 时的加密密钥（key）: ")
							key := c.ReadLine()
							if key == "" {
								log.Warn("加密密钥不能为空")
								return
							}
							data["key"] = key
						}
					case "grpc":
						c.Print("gRPC 的 ServiceName（serviceName）: ")
						serviceName := c.ReadLine()
						if serviceName != "" {
							data["serviceName"] = serviceName
						}
						grpcModeList := []string{
							"gun",
							"multi",
						}
						index := c.MultiChoice(grpcModeList, "gRPC 的传输模式（mode）?")
						mode := grpcModeList[index]
						if mode != "gun" {
							data["mode"] = mode
						}
					}
					securityList := []string{
						"",
						"tls",
						"reality",
					}
					security_index := c.MultiChoice(securityList, "底层传输安全（security）?")
					security := securityList[security_index]
					switch security {
					case "":
					case "tls":
						data["security"] = security
						c.Print("SNI（sni）: ")
						sni := c.ReadLine()
						if sni != "" {
							data["sni"] = sni
						}
						alpnList := []string{
							"",
							"h3",
							"h2",
							"http/1.1",
							"h3,h2",
							"h2,http/1.1",
							"h3,h2,http/1.1",
						}
						index := c.MultiChoice(alpnList, "Alpn ?")
						if alpnList[index] != "" {
							data["alpn"] = alpnList[index]
						}
						c.Print("EchConfigList(默认为空)：")
						data["echConfigList"] = c.ReadLine()
						echForceQueryList := []string{
							"",
							"full",
							"half",
							"none",
						}
						index = c.MultiChoice(echForceQueryList, "ECH强制查询，可留空（echForceQuery）?")
						if echForceQueryList[index] != "" {
							data["echForceQuery"] = echForceQueryList[index]
						}
					case "reality":
						data["security"] = security
						c.Print("SNI（sni）: ")
						sni := c.ReadLine()
						if sni != "" {
							data["sni"] = sni
						}
						fpList := []string{
							"",
							"chrome",
							"firefox",
							"safari",
							"ios",
							"android",
							"edge",
							"360",
							"qq",
							"random",
							"randomized",
						}
						index := c.MultiChoice(fpList, "指纹（FingerPrint） ?")
						if fpList[index] != "" {
							data["fp"] = fpList[index]
						}
						c.Print("公钥（PublicKey）: ")
						data["pbk"] = c.ReadLine()
						c.Print("ShortId: ")
						data["sid"] = c.ReadLine()
						c.Print("SpiderX: ")
						data["spx"] = c.ReadLine()
						c.Print("客户端公钥可留空(mldsa65Verify): ")
						data["pqv"] = c.ReadLine()
					}
					vmessAEAD := &protocols.VMessAEAD{
						ID:      id,
						Address: address,
						Port:    port,
						Remarks: remarks,
						Values:  url.Values{},
					}
					for k, v := range data {
						vmessAEAD.Values[k] = []string{v}
					}
					if manage.Manager.AddNode(node.NewNodeByData(vmessAEAD)) {
						c.Println("添加成功")
					}
				case protocols.ModeVLESS.String():
					c.Println("========================")
					c.Println(protocolMode)
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					address := c.ReadLine()
					c.Print("端口（port）: ")
					port, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("端口为数字，且取值为1~65535")
						return
					}
					c.Print("用户ID（id）: ")
					id := c.ReadLine()
					data := make(map[string]string)
					flowList := []string{
						"",
						"xtls-rprx-vision",
						"xtls-rprx-vision-udp443",
						"xtls-rprx-origin",
						"xtls-rprx-origin-udp443",
						"xtls-rprx-direct",
						"xtls-rprx-direct-udp443",
						"xtls-rprx-splice",
						"xtls-rprx-splice-udp443",
					}
					index := c.MultiChoice(flowList, "流控（flow）?")
					flow := flowList[index]
					data["flow"] = flow
					networkList := []string{
						"tcp",
						"kcp",
						"ws",
						"h2",
						"quic",
						"grpc",
						"xhttp",
						"splithttp",
					}
					index = c.MultiChoice(networkList, "传输协议（network）?")
					network := networkList[index]
					if network != "tcp" {
						data["type"] = network
					}
					switch network {
					case "tcp":
						typeList := []string{
							"none",
							"http",
						}
						index = c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[index] != "none" {
							data["headerType"] = typeList[index]
						}
					case "kcp":
						typeList := []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
						index = c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[index] != "none" {
							data["headerType"] = typeList[index]
						}
						c.Print("KCP种子（seed）: ")
						seed := c.ReadLine()
						if seed != "" {
							data["seed"] = seed
						}
					case "ws":
						c.Print("WebSocket 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("WebSocket Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}
					case "h2":
						c.Print("HTTP/2 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("HTTP/2 Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}
					case "quic":
						typeList := []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
						index = c.MultiChoice(typeList, "伪装类型（headerType）?")
						data["headerType"] = typeList[index]
						quicSecurityList := []string{
							"none",
							"aes-128-gcm",
							"chacha20-poly1305",
						}
						index = c.MultiChoice(quicSecurityList, "QUIC 的加密方式（quicSecurity）?")
						quicSecurity := quicSecurityList[index]
						data["quicSecurity"] = quicSecurity
						if quicSecurity != "none" {
							c.Print("当 QUIC 的加密方式不为 none 时的加密密钥（key）: ")
							key := c.ReadLine()
							if key == "" {
								log.Warn("加密密钥不能为空")
								return
							}
							data["key"] = key
						}
					case "grpc":
						c.Print("gRPC 的 ServiceName（serviceName）: ")
						serviceName := c.ReadLine()
						if serviceName != "" {
							data["serviceName"] = serviceName
						}
						grpcModeList := []string{
							"gun",
							"multi",
						}
						index = c.MultiChoice(grpcModeList, "gRPC 的传输模式（mode）?")
						mode := grpcModeList[index]
						if mode != "gun" {
							data["mode"] = mode
						}
					case "splithttp":
						c.Print("SplitHTTP 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("SplitHTTP Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}
						c.Print("SplitHTTP Mode（mode）: ")
						mode := c.ReadLine()
						if mode != "" {
							data["mode"] = mode
						}
						c.Print("SplitHTTP 额外信息（extra）: ")
						extra := c.ReadLine()
						if extra != "" {
							data["extra"] = mode
						}
					case "xhttp":
						c.Print("xhttp 的路径（path）: ")
						path := c.ReadLine()
						if path != "" {
							data["path"] = path
						}
						c.Print("xhttp Host（host）: ")
						host := c.ReadLine()
						if host != "" {
							data["host"] = host
						}
						c.Print("xhttp Mode（mode）: ")
						mode := c.ReadLine()
						if mode != "" {
							data["mode"] = mode
						}
						c.Print("xhttp 额外信息（extra）: ")
						extra := c.ReadLine()
						if extra != "" {
							data["extra"] = mode
						}
					}
					securityList := []string{
						"",
						"tls",
						"xtls",
						"reality",
					}
					index = c.MultiChoice(securityList, "底层传输安全（security）?")
					security := securityList[index]
					switch security {
					case "":
					case "tls":
						data["security"] = security
						c.Print("SNI（sni）: ")
						sni := c.ReadLine()
						if sni != "" {
							data["sni"] = sni
						}
						//  ALPN参数选取，用于 TLS/QUIC 等协议的握手，告诉服务器客户端支持哪些应用层协议（比如 HTTP/2、HTTP/3）。

						alpnList := []string{
							"",
							"h3",
							"h2",
							"http/1.1",
							"h3,h2",
							"h2,http/1.1",
							"h3,h2,http/1.1",
						}
						index := c.MultiChoice(alpnList, "Alpn ?")
						if alpnList[index] != "" {
							data["alpn"] = alpnList[index]
						}
						c.Print("EchConfigList(默认为空)：")
						echConfigList := c.ReadLine()
						if echConfigList != "" {
							data["echConfigList"] = echConfigList
						}
						echForceQueryList := []string{
							"",
							"full",
							"half",
							"none",
						}
						index = c.MultiChoice(echForceQueryList, "ECH强制查询，可留空（echForceQuery）?")
						if echForceQueryList[index] != "" {
							data["echForceQuery"] = echForceQueryList[index]
						}
					case "xtls":
						data["security"] = security
						c.Print("SNI（sni）: ")
						sni := c.ReadLine()
						if sni != "" {
							data["sni"] = sni
						}
						alpnList := []string{
							"",
							"h3",
							"h2",
							"http/1.1",
							"h3,h2",
							"h2,http/1.1",
							"h3,h2,http/1.1",
						}
						index := c.MultiChoice(alpnList, "Alpn ?")
						if alpnList[index] != "" {
							data["alpn"] = alpnList[index]
						}
						c.Print("EchConfigList(默认为空)：")
						data["echConfigList"] = c.ReadLine()
						echForceQueryList := []string{
							"",
							"full",
							"half",
							"none",
						}
						index = c.MultiChoice(echForceQueryList, "ECH强制查询，可留空（echForceQuery）?")
						if echForceQueryList[index] != "" {
							data["echForceQuery"] = echForceQueryList[index]
						}
					case "reality":
						data["security"] = security
						c.Print("SNI（sni）: ")
						sni := c.ReadLine()
						if sni != "" {
							data["sni"] = sni
						}
						fpList := []string{
							"",
							"chrome",
							"firefox",
							"safari",
							"ios",
							"android",
							"edge",
							"360",
							"qq",
							"random",
							"randomized",
						}
						index := c.MultiChoice(fpList, "指纹（FingerPrint） ?")
						if fpList[index] != "" {
							data["fp"] = fpList[index]
						}
						c.Print("公钥（PublicKey）: ")
						data["pbk"] = c.ReadLine()
						c.Print("ShortId: ")
						data["sid"] = c.ReadLine()
						c.Print("SpiderX: ")
						data["spx"] = c.ReadLine()
						c.Print("MLDSA65验证,可留空（Mldsa65Verify）: ")
						data["pqv"] = c.ReadLine()
					}
					vless := &protocols.VLess{
						ID:      id,
						Address: address,
						Port:    port,
						Remarks: remarks,
						Values:  url.Values{},
					}
					for k, v := range data {
						vless.Values[k] = []string{v}
					}
					if manage.Manager.AddNode(node.NewNodeByData(vless)) {
						c.Println("添加成功")
					}
				case protocols.ModeVMess.String():
					vmess := &protocols.VMess{V: "2"}
					c.Println("========================")
					c.Println(protocolMode)
					c.Println("========================")
					c.Print("别名（remarks）: ")
					vmess.Ps = c.ReadLine()
					c.Print("地址（address）: ")
					vmess.Add = c.ReadLine()
					c.Print("端口（port）: ")
					port, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("端口为数字，且取值为1~65535")
						return
					}
					vmess.Port = port
					c.Print("用户ID（id）: ")
					vmess.Id = c.ReadLine()
					c.Print("额外ID（alterID）: ")
					alterID, err := strconv.Atoi(c.ReadLine())
					if err != nil {
						log.Warn("额外ID为数字")
						return
					}
					vmess.Aid = alterID
					securityList := []string{
						"auto",
						"aes-128-gcm",
						"chacha20-poly1305",
						"none",
						"zero",
					}
					index := c.MultiChoice(securityList, "加密方式（security）?")
					vmess.Scy = securityList[index]
					networkList := []string{
						"tcp",
						"kcp",
						"ws",
						"h2",
						"quic",
						"grpc",
					}
					index = c.MultiChoice(networkList, "传输协议（network）?")
					vmess.Net = networkList[index]
					switch networkList[index] {
					case "tcp":
						vmess.Type = "none"
					case "kcp":
						typeList := []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
						index = c.MultiChoice(typeList, "伪装头部类型（type）?")
						vmess.Type = typeList[index]
						c.Print("mKCP 种子（path）: ")
						vmess.Path = c.ReadLine()
					case "quic":
						typeList := []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
						index = c.MultiChoice(typeList, "伪装类型（type）?")
						vmess.Type = typeList[index]
						quicSecurityList := []string{
							"none",
							"aes-128-gcm",
							"chacha20-poly1305",
						}
						index = c.MultiChoice(quicSecurityList, "QUIC的加密方式（host）?")
						vmess.Host = quicSecurityList[index]
						if vmess.Host != "none" {
							c.Print("QUIC的加密key（path）: ")
							vmess.Path = c.ReadLine()
						}

					case "ws", "h2":
						c.Print("Host（host）: ")
						vmess.Host = c.ReadLine()
						c.Print("Path（path）: ")
						vmess.Path = c.ReadLine()
					case "grpc":
						typeList := []string{
							"gun",
							"multi",
						}
						index = c.MultiChoice(typeList, "gRPC的传输模式（type）?")
						vmess.Type = typeList[index]
						c.Print("gRPC的ServiceName（path）: ")
						vmess.Path = c.ReadLine()
					}
					tlsList := []string{
						"",
						"tls",
					}
					index = c.MultiChoice(tlsList, "底层安全传输 （tls）?")
					vmess.Tls = tlsList[index]
					if vmess.Tls != "" {
						c.Print("SNI: ")
						vmess.Sni = c.ReadLine()
						alpnList := []string{
							"",
							"h3",
							"h2",
							"http/1.1",
							"h3,h2",
							"h2,http/1.1",
							"h3,h2,http/1.1",
						}
						index := c.MultiChoice(alpnList, "Alpn ?")
						vmess.Alpn = alpnList[index]
						c.Print("EchConfigList(默认为空)：")
						vmess.EchConfigList = c.ReadLine()
						echForceQueryList := []string{
							"",
							"full",
							"half",
							"none",
						}						
						index = c.MultiChoice(echForceQueryList, "ECH强制查询，可留空（echForceQuery）?")
						vmess.EchForceQuery = echForceQueryList[index]
					}
					c.Println("========================")
					if manage.Manager.AddNode(node.NewNodeByData(vmess)) {
						c.Println("添加成功")
					}
				case protocols.ModeShadowSocks.String():
					c.Println("========================")
					c.Println(protocolMode)
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					addr := c.ReadLine()
					c.Print("端口（port）: ")
					port, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("端口为数字，且取值为1~65535")
						return
					}
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
					ss := &protocols.ShadowSocks{
						Remarks:  remarks,
						Password: password,
						Address:  addr,
						Port:     port,
						Method:   security[sIndex],
					}
					if manage.Manager.AddNode(node.NewNodeByData(ss)) {
						c.Println("添加成功")
					}
				case protocols.ModeTrojan.String():
					c.Println("========================")
					c.Println(protocolMode)
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					addr := c.ReadLine()
					c.Print("端口（port）: ")
					port, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("端口为数字，且取值为1~65535")
						return
					}
					c.Print("密码（password）: ")
					password := c.ReadLine()
					c.Print("SNI（sni）（可选）: ")
					sni := c.ReadLine()
					c.Print("EchConfigList(默认为空)：")
					echConfigList := c.ReadLine()
					echForceQueryList := []string{
						"",
						"full",
						"half",
						"none",
					}
					index := c.MultiChoice(echForceQueryList, "ECH强制查询（echForceQuery）?")
					echForceQuery := echForceQueryList[index]
					c.Println("========================")
					trojan := &protocols.Trojan{
						Remarks:  remarks,
						Password: password,
						Address:  addr,
						Port:     port,
					}
					values := url.Values{}
					if sni != "" {
						values["sni"] = []string{sni}
					}
					if echConfigList != "" {
						values["echConfigList"] = []string{echConfigList}
					}
					if echForceQuery != "" {
						values["echForceQuery"] = []string{echForceQuery}
					}
					if len(values) > 0 {
						trojan.Values = values
					}
					if manage.Manager.AddNode(node.NewNodeByData(trojan)) {
						c.Println("添加成功")
					}
				case protocols.ModeSocks.String():
					c.Println("========================")
					c.Println(protocolMode)
					c.Println("========================")
					c.Print("别名（remarks）: ")
					remarks := c.ReadLine()
					c.Print("地址（address）: ")
					addr := c.ReadLine()
					c.Print("端口（port）: ")
					port, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("端口为数字，且取值为1~65535")
						return
					}
					c.Print("用户名（username）（可选）: ")
					username := c.ReadLine()
					c.Print("密码（password）（可选）: ")
					password := c.ReadLine()
					c.Println("========================")
					socks := &protocols.Socks{
						Remarks: remarks,
						Address: addr,
						Port:    port,
					}
					if username != "" && password != "" {
						socks.Username = username
						socks.Password = password
					}
					if manage.Manager.AddNode(node.NewNodeByData(socks)) {
						c.Println("添加成功")
					}
				}
			}
		},
	})
	shell.AddCmd(nodeCmd)
}
