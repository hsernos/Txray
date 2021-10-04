package cmd

import (
	"Txray/cmd/help"
	"Txray/core"
	"Txray/core/manage"
	"Txray/core/node"
	"Txray/core/protocols"
	"Txray/core/sub"
	"Txray/log"
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func InitNodeShell(shell *ishell.Shell) {
	nodeCmd := &ishell.Cmd{
		Name: "node",
		Func: func(c *ishell.Context) {
			argMap := FlagsParse(c.Args, map[string]string{
				"d": "desc",
			})
			key := "all"
			if _, ok := argMap["data"]; ok {
				key, _ = argMap["data"]
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
				key, _ = argMap["data"]
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
				data, _ := ioutil.ReadFile(fileArg)
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
				c.ShowPrompt(true)
				protocolMode := modeList[i]
				c.Println()
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
					data := make(map[string]string, 0)
					networkList := []string{
						"tcp",
						"kcp",
						"ws",
						"h2",
						"quic",
						"grpc",
					}
					network_index := c.MultiChoice(networkList, "传输协议（network）?")
					network := networkList[network_index]
					if network != "tcp" {
						data["type"] = network
					}
					switch network {
					case "tcp":
						typeList := []string{
							"none",
							"http",
						}
						type_index := c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[type_index] != "none" {
							data["headerType"] = typeList[type_index]
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
						type_index := c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[type_index] != "none" {
							data["headerType"] = typeList[type_index]
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
						type_index := c.MultiChoice(typeList, "伪装类型（headerType）?")
						data["headerType"] = typeList[type_index]
						quicSecurityList := []string{
							"none",
							"aes-128-gcm",
							"chacha20-poly1305",
						}
						quicSecurity_index := c.MultiChoice(quicSecurityList, "QUIC 的加密方式（quicSecurity）?")
						quicSecurity := quicSecurityList[quicSecurity_index]
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
						grpcMode_index := c.MultiChoice(grpcModeList, "gRPC 的传输模式（mode）?")
						mode := grpcModeList[grpcMode_index]
						if mode != "gun" {
							data["mode"] = mode
						}
					}
					securityList := []string{
						"none",
						"tls",
					}
					security_index := c.MultiChoice(securityList, "底层传输安全（security）?")
					security := securityList[security_index]
					if security != "none" {
						data["security"] = security
					}
					c.Print("SNI（sni）: ")
					sni := c.ReadLine()
					if sni != "" {
						data["sni"] = sni
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
					data := make(map[string]string, 0)
					networkList := []string{
						"tcp",
						"kcp",
						"ws",
						"h2",
						"quic",
						"grpc",
					}
					network_index := c.MultiChoice(networkList, "传输协议（network）?")
					network := networkList[network_index]
					if network != "tcp" {
						data["type"] = network
					}
					switch network {
					case "tcp":
						typeList := []string{
							"none",
							"http",
						}
						type_index := c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[type_index] != "none" {
							data["headerType"] = typeList[type_index]
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
						type_index := c.MultiChoice(typeList, "伪装类型（headerType） ?")
						if typeList[type_index] != "none" {
							data["headerType"] = typeList[type_index]
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
						type_index := c.MultiChoice(typeList, "伪装类型（headerType）?")
						data["headerType"] = typeList[type_index]
						quicSecurityList := []string{
							"none",
							"aes-128-gcm",
							"chacha20-poly1305",
						}
						quicSecurity_index := c.MultiChoice(quicSecurityList, "QUIC 的加密方式（quicSecurity）?")
						quicSecurity := quicSecurityList[quicSecurity_index]
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
						grpcMode_index := c.MultiChoice(grpcModeList, "gRPC 的传输模式（mode）?")
						mode := grpcModeList[grpcMode_index]
						if mode != "gun" {
							data["mode"] = mode
						}
					}

					securityList := []string{
						"none",
						"tls",
						"xtls",
					}
					security_index := c.MultiChoice(securityList, "底层传输安全（security）?")
					security := securityList[security_index]
					if security != "none" {
						data["security"] = security
					}
					c.Print("SNI（sni）: ")
					sni := c.ReadLine()
					if sni != "" {
						data["sni"] = sni
					}
					switch security {
					case "xtls":
						flowList := []string{
							"xtls-rprx-origin",
							"xtls-rprx-origin-udp443",
							"xtls-rprx-direct",
							"xtls-rprx-direct-udp443",
							"xtls-rprx-splice",
							"xtls-rprx-splice-udp443",
						}
						flow_index := c.MultiChoice(flowList, "流控（flow）?")
						flow := flowList[flow_index]
						data["flow"] = flow
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
					c.Print("额外ID（alterID）: ")
					alterID, err := strconv.Atoi(c.ReadLine())
					if err != nil || port < 1 || port > 65535 {
						log.Warn("额外ID为数字，且取值为1~65535")
						return
					}
					networkList := []string{
						"tcp",
						"kcp",
						"ws",
						"h2",
						"quic",
						"grpc",
					}
					network := c.MultiChoice(networkList, "传输协议（network）?")

					typeList := []string{
						"none",
						"http",
					}
					switch networkList[network] {
					case "tcp":
						typeList = []string{
							"none",
							"http",
						}
					case "kcp", "quic":
						typeList = []string{
							"none",
							"srtp",
							"utp",
							"wechat-video",
							"dtls",
							"wireguard",
						}
					case "ws", "h2":
						typeList = []string{
							"none",
						}
					case "grpc":
						typeList = []string{
							"gun",
							"multi",
						}
					}
					types := 0
					if networkList[network] != "ws" && networkList[network] != "h2" {
						types = c.MultiChoice(typeList, "伪装类型（type） ?")
					}
					c.Print("伪装域名 host（ws host/h2 host(多个以,隔开)/QUIC 加密方式）: ")
					host := c.ReadLine()
					c.Print("path（ws path/h2 path/QUIC 加密秘钥/Kcp Seed/grpc ServerName）: ")
					path := c.ReadLine()
					tlsList := []string{
						"tls",
						"",
					}
					tls := c.MultiChoice(tlsList, "底层安全传输 ?")
					c.Println("========================")
					vmess := &protocols.VMess{
						Ps:   remarks,
						Add:  address,
						Port: port,
						Id:   id,
						Aid:  alterID,
						Net:  networkList[network],
						Type: typeList[types],
						Host: host,
						Path: path,
						Tls:  tlsList[tls],
					}
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
					c.Println("========================")
					trojan := &protocols.Trojan{
						Remarks:  remarks,
						Password: password,
						Address:  addr,
						Port:     port,
					}
					if sni != "" {
						trojan.Values = url.Values{
							"sni": []string{sni},
						}
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
