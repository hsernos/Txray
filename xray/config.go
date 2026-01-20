// xray/config.go 负责生成 xray-core 的配置文件，包含日志、入站、出站、策略、DNS、路由等配置项
package xray

import (
	"Txray/core"                 // 配置目录、文件写入
	"Txray/core/protocols"       // 协议定义
	"Txray/core/protocols/field" // 字段定义
	"Txray/core/routing"         // 路由配置
	"Txray/core/setting"         // 设置项
	"Txray/log"                  // 日志
	"path/filepath"              // 路径处理
	"strings"                    // 字符串处理
)

// GenConfig 生成 xray-core 配置文件，返回配置文件路径
// node: 协议节点对象
func GenConfig(node protocols.Protocol) string {
	path := filepath.Join(core.GetConfigDir(), "config.json") // 配置文件路径
	var conf = map[string]interface{}{
		"version":   versionConfig(),      // 版本配置
		"log":       logConfig(),          // 日志配置
		"inbounds":  inboundsConfig(),     // 入站配置
		"outbounds": outboundConfig(node), // 出站配置
		"policy":    policyConfig(),       // 策略配置
		"dns":       dnsConfig(),          // DNS 配置
		"routing":   routingConfig(),      // 路由配置
	}
	err := core.WriteJSON(conf, path) // 写入 JSON 文件
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return path
}

// versionConfig 生成版本配置
func versionConfig() interface{} {
	version := map[string]interface{}{
		"min": setting.VersionMin(),
		"max": setting.VersionMax(),
	}
	return version
}

// logConfig 生成日志配置
func logConfig() interface{} {
	path := core.LogFile
	return map[string]string{
		"access":   path,      // 访问日志路径
		"loglevel": "warning", // 日志级别
	}
}

// inboundsConfig 生成入站配置
func inboundsConfig() interface{} {
	listen := "127.0.0.1" // 默认仅本地监听
	if setting.FromLanConn() {
		listen = "0.0.0.0" // 允许局域网访问
	}
	data := []interface{}{
		// 默认处理混合入站流量
		map[string]interface{}{
			"tag":      "mixed",
			"port":     setting.Mixed(),
			"listen":   listen,
			"protocol": "mixed",
			"sniffing": map[string]interface{}{
				"enabled": setting.Sniffing(),
				"destOverride": []string{
					"http",
					"tls",
					"quic",
					"fakedns",
					"fakedns+others",
				},
			},
			"settings": map[string]interface{}{
				"auth":      "noauth",
				"udp":       setting.UDP(),
				"userLevel": 0,
			},
		},
	}
	if setting.Socks() > 0 {
		data = append(data, map[string]interface{}{ // 添加 Socks5 入站
			"tag":      "proxy",
			"port":     setting.Socks(),
			"listen":   listen,
			"protocol": "socks",
			"sniffing": map[string]interface{}{
				"enabled": setting.Sniffing(),
				"destOverride": []string{
					"http",
					"tls",
				},
			},
			"settings": map[string]interface{}{
				"auth":      "noauth",
				"udp":       setting.UDP(),
				"userLevel": 0,
			},
		})
	}
	if setting.Http() > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "http",
			"port":     setting.Http(),
			"listen":   listen,
			"protocol": "http",
			"settings": map[string]interface{}{
				"userLevel": 0,
			},
		})
	}
	if setting.DNSPort() > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "dns-in",
			"port":     setting.DNSPort(),
			"listen":   listen,
			"protocol": "dokodemo-door",
			"settings": map[string]interface{}{
				"userLevel": 0,
				"address":   setting.DNSForeign(),
				"network":   "tcp,udp",
				"port":      53,
			},
		})
	}
	return data
}

// 本地策略
func policyConfig() interface{} {
	return map[string]interface{}{
		"levels": map[string]interface{}{
			"0": map[string]interface{}{
				"handshake":    4,
				"connIdle":     300,
				"uplinkOnly":   1,
				"downlinkOnly": 1,
				"bufferSize":   10240,
			},
		},
		"system": map[string]interface{}{
			"statsInboundUplink":   true,
			"statsInboundDownlink": true,
		},
	}
}

// DNS
func dnsConfig() interface{} {
	servers := make([]interface{}, 0)
	if setting.DNSDomestic() != "" {
		servers = append(servers, map[string]interface{}{
			"address": setting.DNSDomestic(),
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if setting.DNSBackup() != "" {
		servers = append(servers, map[string]interface{}{
			"address": setting.DNSBackup(),
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if setting.DNSForeign() != "" {
		servers = append(servers, map[string]interface{}{
			"address": setting.DNSForeign(),
			"port":    53,
			"domains": []interface{}{
				"geosite:geolocation-!cn",
				"geosite:speedtest",
			},
		})
	}
	return map[string]interface{}{
		"hosts": map[string]interface{}{
			"domain:googleapis.cn": "googleapis.com",
		},
		"servers": servers,
	}
}

// 路由
func routingConfig() interface{} {
	rules := make([]interface{}, 0)
	if setting.DNSPort() != 0 {
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"inboundTag": []interface{}{
				"dns-in",
			},
			"outboundTag": "dns-out",
		})
	}
	if setting.DNSForeign() != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "proxy",
			"ip": []string{
				setting.DNSForeign(),
			},
		})
	}
	if setting.DNSDomestic() != "" || setting.DNSBackup() != "" {
		var ip []string
		if setting.DNSDomestic() != "" {
			ip = append(ip, setting.DNSDomestic())
		}
		if setting.DNSBackup() != "" {
			ip = append(ip, setting.DNSBackup())
		}
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "direct",
			"ip":          ip,
		})
	}
	ips, domains := routing.GetRulesGroupData(routing.TypeBlock)
	if len(ips) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "block",
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "block",
			"domain":      domains,
		})
	}
	ips, domains = routing.GetRulesGroupData(routing.TypeDirect)
	if len(ips) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"domain":      domains,
		})
	}
	ips, domains = routing.GetRulesGroupData(routing.TypeProxy)
	if len(ips) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "proxy",
			"ip":          ips,
		})
	}
	if len(domains) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "proxy",
			"domain":      domains,
		})
	}

	if setting.RoutingBypass() {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"ip": []string{
				"geoip:private",
				"geoip:cn",
			},
		})
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"outboundTag": "direct",
			"domain": []string{
				"geosite:cn",
			},
		})
	}
	return map[string]interface{}{
		"domainStrategy": setting.RoutingStrategy(),
		"rules":          rules,
	}
}

// 出站
func outboundConfig(n protocols.Protocol) interface{} {
	out := make([]interface{}, 0)
	switch n.GetProtocolMode() {
	case protocols.ModeTrojan:
		t := n.(*protocols.Trojan)
		out = append(out, trojanOutbound(t))
	case protocols.ModeShadowSocks:
		ss := n.(*protocols.ShadowSocks)
		out = append(out, shadowsocksOutbound(ss))
	case protocols.ModeVMess:
		v := n.(*protocols.VMess)
		out = append(out, vMessOutbound(v))
	case protocols.ModeSocks:
		v := n.(*protocols.Socks)
		out = append(out, socksOutbound(v))
	case protocols.ModeVLESS:
		v := n.(*protocols.VLess)
		out = append(out, vLessOutbound(v))
	case protocols.ModeVMessAEAD:
		v := n.(*protocols.VMessAEAD)
		out = append(out, vMessAEADOutbound(v))
	}
	out = append(out, map[string]interface{}{
		"tag":      "direct",
		"protocol": "freedom",
		"settings": map[string]interface{}{},
	})
	out = append(out, map[string]interface{}{
		"tag":      "block",
		"protocol": "blackhole",
		"settings": map[string]interface{}{
			"response": map[string]interface{}{
				"type": "http",
			},
		},
	})
	out = append(out, map[string]interface{}{
		"tag":      "dns-out",
		"protocol": "dns",
	})
	return out
}

// Shadowsocks
func shadowsocksOutbound(ss *protocols.ShadowSocks) interface{} {
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "shadowsocks",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  ss.Address,
					"port":     ss.Port,
					"password": ss.Password,
					"method":   ss.Method,
					"level":    0,
				},
			},
		},
		"streamSettings": map[string]interface{}{
			"network": "tcp",
		},
	}
}

// Trojan
func trojanOutbound(trojan *protocols.Trojan) interface{} {
	streamSettings := map[string]interface{}{
		"network":  "tcp",
		"security": "tls",
	}
	if trojan.Sni() != "" {
		tlsSettings := map[string]interface{}{
			"allowInsecure": setting.AllowInsecure(),
			"serverName":    trojan.Sni(),
		}
		if trojan.EchConfigList() != "" {
			tlsSettings["echConfigList"] = trojan.EchConfigList()
		}
		if trojan.EchForceQuery() != "" {
			tlsSettings["echForceQuery"] = trojan.EchForceQuery()
		}
		if trojan.Has("pinnedPeerCertSha256") && trojan.Get("pinnedPeerCertSha256") != "" {
			tlsSettings["pinnedPeerCertSha256"] = trojan.Get("pinnedPeerCertSha256")
		}
		streamSettings["tlsSettings"] = tlsSettings
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "trojan",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  trojan.Address,
					"port":     trojan.Port,
					"password": trojan.Password,
					"level":    0,
				},
			},
		},
		"streamSettings": streamSettings,
	}
}

// VMess
func vMessOutbound(vmess *protocols.VMess) interface{} {
	mux := setting.Mux()
	streamSettings := map[string]interface{}{
		"network":  vmess.Net,
		"security": vmess.Tls,
	}
	if vmess.Tls == "tls" {
		tlsSettings := map[string]interface{}{
			"allowInsecure": setting.AllowInsecure(),
		}
		if vmess.Sni != "" {
			tlsSettings["serverName"] = vmess.Sni
		}
		if vmess.Alpn != "" {
			tlsSettings["alpn"] = strings.Split(vmess.Alpn, ",")
		}
		if vmess.EchConfigList != "" {
			tlsSettings["echConfigList"] = vmess.EchConfigList
		}
		if vmess.EchForceQuery != "" {
			tlsSettings["echForceQuery"] = vmess.EchForceQuery
		}
		if vmess.PCS != "" {
			tlsSettings["pinnedPeerCertSha256"] = vmess.PCS
		}
		streamSettings["tlsSettings"] = tlsSettings
	}
	switch vmess.Net {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
		if vmess.Type != "none" {
			kcpSettings["seed"] = vmess.Path
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vmess.Path,
			"headers": map[string]interface{}{
				"Host": vmess.Host,
			},
		}
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vmess.Host, ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vmess.Path,
			"host": host,
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vmess.Host,
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
		if vmess.Host != "none" {
			quicSettings["key"] = vmess.Path
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vmess.Path,
			"multiMode":   vmess.Type == "multi",
		}
	case "splithttp":
		streamSettings["splithttpSettings"] = map[string]interface{}{
			"host":  vmess.GetValue(field.SpHost),
			"path":  vmess.GetValue(field.SpPath),
			"mode":  vmess.GetValue(field.SpMode),
			"extra": vmess.GetExtraValue(field.SpExtra),
		}
	case "xhttp":
		streamSettings["xhttpSettings"] = map[string]interface{}{
			"host":  vmess.GetValue(field.XhHost),
			"path":  vmess.GetValue(field.XhPath),
			"mode":  vmess.GetValue(field.XhMode),
			"extra": vmess.GetExtraValue(field.XhExtra),
		}
	case "xhttpupgrade":
		streamSettings["xhttpUpgradeSettings"] = map[string]interface{}{
			"host": vmess.GetValue(field.XhHost),
			"path": vmess.GetValue(field.XhPath),
		}
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vmess.Add,
					"port":    vmess.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":       vmess.Id,
							"alterId":  vmess.Aid,
							"security": vmess.Scy,
							"level":    0,
						},
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// socks
func socksOutbound(socks *protocols.Socks) interface{} {
	user := map[string]interface{}{
		"address": socks.Address,
		"port":    socks.Port,
	}
	if socks.Username != "" || socks.Password != "" {
		user["users"] = []interface{}{
			map[string]interface{}{
				"user": socks.Username,
				"pass": socks.Password,
			},
		}
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "socks",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				user,
			},
		},
		"streamSettings": map[string]interface{}{
			"network": "tcp",
			"tcpSettings": map[string]interface{}{
				"header": map[string]interface{}{
					"type": "none",
				},
			},
		},
		"mux": map[string]interface{}{
			"enabled": false,
		},
	}
}

// VLESS
func vLessOutbound(vless *protocols.VLess) interface{} {
	mux := setting.Mux()
	security := vless.GetValue(field.Security)
	network := vless.GetValue(field.NetworkType)
	user := map[string]interface{}{
		"id":         vless.ID,
		"flow":       vless.GetValue(field.Flow),
		"encryption": vless.GetValue(field.VLessEncryption),
		"level":      0,
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		tlsSettings := map[string]interface{}{
			"allowInsecure": setting.AllowInsecure(),
		}
		sni := vless.GetHostValue(field.SNI)
		alpn := vless.GetValue(field.Alpn)
		echConfigList := vless.GetValue(field.EchConfigList)
		echForceQuery := vless.GetValue(field.EchForceQuery)
		pcs := vless.GetValue(field.PCS)
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if alpn != "" {
			tlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		if echConfigList != "" {
			tlsSettings["echConfigList"] = echConfigList
		}
		if echForceQuery != "" {
			tlsSettings["echForceQuery"] = echForceQuery
		}
		if pcs != "" {
			tlsSettings["pinnedPeerCertSha256"] = pcs
		}
		streamSettings["tlsSettings"] = tlsSettings
	case "xtls":
		xtlsSettings := map[string]interface{}{
			"allowInsecure": setting.AllowInsecure(),
		}
		sni := vless.GetHostValue(field.SNI)
		alpn := vless.GetValue(field.Alpn)
		echConfigList := vless.GetValue(field.EchConfigList)
		echForceQuery := vless.GetValue(field.EchForceQuery)
		pcs := vless.GetValue(field.PCS)
		if sni != "" {
			xtlsSettings["serverName"] = sni
		}
		if alpn != "" {
			xtlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		if echConfigList != "" {
			xtlsSettings["echConfigList"] = echConfigList
		}
		if echForceQuery != "" {
			xtlsSettings["echForceQuery"] = echForceQuery
		}
		if pcs != "" {
			xtlsSettings["pinnedPeerCertSha256"] = pcs
		}
		streamSettings["xtlsSettings"] = xtlsSettings
		mux = false
	case "reality":
		realitySettings := map[string]interface{}{
			//REALITY 服务端、客户端配置填入 "show": true 即可输出 X25519MLKEM768、ML-DSA-65 相关日志，以确定它们被用到。
			"show":          false,
			"fingerprint":   vless.GetValue(field.TLSFingerPrint),
			"serverName":    vless.GetHostValue(field.SNI),
			"publicKey":     vless.GetValue(field.RealityPublicKey),
			"shortId":       vless.GetValue(field.RealityShortId),
			"spiderX":       vless.GetValue(field.RealitySpiderX),
			"mldsa65Verify": vless.GetValue(field.RealityMldsa65Verify),
		}
		streamSettings["realitySettings"] = realitySettings
		mux = false
	}
	switch network {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vless.GetValue(field.TCPHeaderType),
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vless.GetValue(field.MkcpHeaderType),
			},
		}
		if vless.Has(field.Seed.Key) {
			kcpSettings["seed"] = vless.GetValue(field.Seed)
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vless.GetHostValue(field.H2Host), ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vless.GetValue(field.H2Path),
			"host": host,
		}
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vless.GetValue(field.WsPath),
			"headers": map[string]interface{}{
				"Host": vless.GetValue(field.WsHost),
			},
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vless.GetValue(field.QuicSecurity),
			"header": map[string]interface{}{
				"type": vless.GetValue(field.QuicHeaderType),
			},
		}
		if vless.GetValue(field.QuicSecurity) != "none" {
			quicSettings["key"] = vless.GetValue(field.QuicKey)
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vless.GetValue(field.GrpcServiceName),
			"multiMode":   vless.GetValue(field.GrpcMode) == "multi",
		}
	case "splithttp":
		streamSettings["splithttpSettings"] = map[string]interface{}{
			"host":  vless.GetValue(field.SpHost),
			"path":  vless.GetValue(field.SpPath),
			"mode":  vless.GetValue(field.SpMode),
			"extra": vless.GetExtraValue(field.SpExtra),
		}
	case "xhttp":
		streamSettings["xhttpSettings"] = map[string]interface{}{
			"host":  vless.GetValue(field.XhHost),
			"path":  vless.GetValue(field.XhPath),
			"mode":  vless.GetValue(field.XhMode),
			"extra": vless.GetExtraValue(field.XhExtra),
		}
	case "xhttpupgrade":
		streamSettings["xhttpUpgradeSettings"] = map[string]interface{}{
			"host": vless.GetValue(field.XhHost),
			"path": vless.GetValue(field.XhPath),
		}
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "vless",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vless.Address,
					"port":    vless.Port,
					"users": []interface{}{
						user,
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// VMessAEAD
func vMessAEADOutbound(vmess *protocols.VMessAEAD) interface{} {
	mux := setting.Mux()
	security := vmess.GetValue(field.Security)
	network := vmess.GetValue(field.NetworkType)
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		tlsSettings := map[string]interface{}{
			"allowInsecure": setting.AllowInsecure(),
		}
		sni := vmess.GetHostValue(field.SNI)
		alpn := vmess.GetValue(field.Alpn)
		echConfigList := vmess.GetValue(field.EchConfigList)
		echForceQuery := vmess.GetValue(field.EchForceQuery)
		pcs := vmess.GetValue(field.PCS)
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if alpn != "" {
			tlsSettings["alpn"] = strings.Split(alpn, ",")
		}
		if echConfigList != "" {
			tlsSettings["echConfigList"] = echConfigList
		}
		if echForceQuery != "" {
			tlsSettings["echForceQuery"] = echForceQuery
		}
		if pcs != "" {
			tlsSettings["pinnedPeerCertSha256"] = pcs
		}
		streamSettings["tlsSettings"] = tlsSettings
	case "reality":
		realitySettings := map[string]interface{}{
			"show":          false,
			"fingerprint":   vmess.GetValue(field.TLSFingerPrint),
			"serverName":    vmess.GetHostValue(field.SNI),
			"publicKey":     vmess.GetValue(field.RealityPublicKey),
			"shortId":       vmess.GetValue(field.RealityShortId),
			"spiderX":       vmess.GetValue(field.RealitySpiderX),
			"mldsa65Verify": vmess.GetValue(field.RealityMldsa65Verify),
		}
		streamSettings["realitySettings"] = realitySettings
		mux = false
	}
	switch network {
	case "tcp":
		streamSettings["tcpSettings"] = map[string]interface{}{
			"header": map[string]interface{}{
				"type": vmess.GetValue(field.TCPHeaderType),
			},
		}
	case "kcp":
		kcpSettings := map[string]interface{}{
			"mtu":              1350,
			"tti":              50,
			"uplinkCapacity":   12,
			"downlinkCapacity": 100,
			"congestion":       false,
			"readBufferSize":   2,
			"writeBufferSize":  2,
			"header": map[string]interface{}{
				"type": vmess.GetValue(field.MkcpHeaderType),
			},
		}
		if vmess.Has(field.Seed.Key) {
			kcpSettings["seed"] = vmess.GetValue(field.Seed)
		}
		streamSettings["kcpSettings"] = kcpSettings
	case "h2":
		mux = false
		host := make([]string, 0)
		for _, line := range strings.Split(vmess.GetHostValue(field.H2Host), ",") {
			line = strings.TrimSpace(line)
			if line != "" {
				host = append(host, line)
			}
		}
		streamSettings["httpSettings"] = map[string]interface{}{
			"path": vmess.GetValue(field.H2Path),
			"host": host,
		}
	case "ws":
		streamSettings["wsSettings"] = map[string]interface{}{
			"path": vmess.GetValue(field.WsPath),
			"headers": map[string]interface{}{
				"Host": vmess.GetValue(field.WsHost),
			},
		}
	case "quic":
		quicSettings := map[string]interface{}{
			"security": vmess.GetValue(field.QuicSecurity),
			"header": map[string]interface{}{
				"type": vmess.GetValue(field.QuicHeaderType),
			},
		}
		if vmess.GetValue(field.QuicSecurity) != "none" {
			quicSettings["key"] = vmess.GetValue(field.QuicKey)
		}
		streamSettings["quicSettings"] = quicSettings
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vmess.GetValue(field.GrpcServiceName),
			"multiMode":   vmess.GetValue(field.GrpcMode) == "multi",
		}
	case "splithttp":
		streamSettings["splithttpSettings"] = map[string]interface{}{
			"host":  vmess.GetValue(field.SpHost),
			"path":  vmess.GetValue(field.SpPath),
			"mode":  vmess.GetValue(field.SpMode),
			"extra": vmess.GetExtraValue(field.SpExtra),
		}
	case "xhttp":
		streamSettings["xhttpSettings"] = map[string]interface{}{
			"host":  vmess.GetValue(field.XhHost),
			"path":  vmess.GetValue(field.XhPath),
			"mode":  vmess.GetValue(field.XhMode),
			"extra": vmess.GetExtraValue(field.XhExtra),
		}
	case "xhttpupgrade":
		streamSettings["xhttpUpgradeSettings"] = map[string]interface{}{
			"host": vmess.GetValue(field.XhHost),
			"path": vmess.GetValue(field.XhPath),
		}
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vmess.Address,
					"port":    vmess.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":       vmess.ID,
							"security": vmess.GetValue(field.VMessEncryption),
							"level":    0,
						},
					},
				},
			},
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}
