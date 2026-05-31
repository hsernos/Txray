package xray

import (
	"Txray/core"
	"Txray/core/protocols"
	"Txray/core/protocols/field"
	"Txray/core/routing"
	"Txray/core/setting"
	"Txray/log"
	"path/filepath"
	"strings"
)

// 生成xray-core配置文件
func GenConfig(node protocols.Protocol) string {
	path := filepath.Join(core.GetConfigDir(), "config.json")
	var conf = map[string]interface{}{
		"log":       logConfig(),
		"inbounds":  inboundsConfig(),
		"outbounds": outboundConfig(node),
		"policy":    policyConfig(),
		"dns":       dnsConfig(),
		"routing":   routingConfig(),
	}
	err := core.WriteJSON(conf, path)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return path
}

// 日志
func logConfig() interface{} {
	path := core.LogFile
	return map[string]string{
		"access":   path,
		"loglevel": "warning",
	}
}

// 入站
func inboundsConfig() interface{} {
	listen := "127.0.0.1"
	if setting.FromLanConn() {
		listen = "0.0.0.0"
	}
	data := []interface{}{
		map[string]interface{}{
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
		},
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
	case protocols.ModeHysteria2:
		v := n.(*protocols.Hysteria2)
		out = append(out, hysteria2Outbound(v))
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
		streamSettings["tlsSettings"] = map[string]interface{}{
			"allowInsecure": setting.AllowInsecure(),
			"serverName":    trojan.Sni(),
		}
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
			"network": "raw",
		},
		"mux": map[string]interface{}{
			"enabled": false,
		},
	}
}

// VMess
func vMessOutbound(vmess *protocols.VMess) interface{} {
	mux := setting.Mux()
	network := vmess.Net
	switch network {
	case "tcp":
		network = "raw"
	case "kcp":
		network = "mkcp"
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": vmess.Tls,
	}
	if vmess.Tls == "tls" {
		streamSettings["tlsSettings"] = genTlsSetting(vmess.Sni, vmess.Alpn, "", setting.AllowInsecure())
	}
	switch network {
	case "tcp", "raw":
		streamSettings["rawSettings"] = genRawSetting(vmess.Type, vmess.Host, vmess.Path)
	case "xhttp":
		streamSettings["xhttpSettings"] = genXhttpSetting(vmess.Host, vmess.Path, vmess.Type, "")
	case "kcp", "mkcp":
		streamSettings["kcpSettings"] = genMkcpSetting("1350")
	case "ws":
		streamSettings["wsSettings"] = genWsSetting(vmess.Host, vmess.Path)
	case "h2":
		streamSettings["httpSettings"] = genHttpSetting(vmess.Host, vmess.Path)
	case "grpc":
		streamSettings["grpcSettings"] = genGrpcSetting(vmess.Type, vmess.Path, vmess.Host)
	case "httpupgrade":
		streamSettings["httpupgradeSettings"] = genHttpupgradeSetting(vmess.Host, vmess.Path)
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

// VMessAEAD
func vMessAEADOutbound(vmess *protocols.VMessAEAD) interface{} {
	mux := setting.Mux()
	security := vmess.GetValue(field.TlsSecurity)
	network := vmess.GetValue(field.NetworkType)
	switch network {
	case "tcp":
		network = "raw"
	case "kcp":
		network = "mkcp"
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		streamSettings["tlsSettings"] = genTlsSetting(
			vmess.GetHostValue(field.SNI),
			vmess.GetValue(field.Alpn),
			vmess.GetValue(field.FingerPrint),
			setting.AllowInsecure(),
		)
	case "reality":
		streamSettings["realitySettings"] = genRealitySetting(
			vmess.GetHostValue(field.SNI),
			vmess.GetValue(field.FingerPrint),
			vmess.GetValue(field.PublicKey),
			vmess.GetValue(field.ShortId),
			vmess.GetValue(field.SpiderX),
			vmess.GetValue(field.Mldsa65Verify),
		)
		mux = false
	}
	switch network {
	case "tcp", "raw":
		streamSettings["rawSettings"] = genRawSetting(
			vmess.GetValue(field.RawHeaderType),
			vmess.GetValue(field.RawHost),
			vmess.GetValue(field.RawPath),
		)
	case "xhttp":
		streamSettings["xhttpSettings"] = genXhttpSetting(
			vmess.GetHostValue(field.XhttpHost),
			vmess.GetValue(field.XhttpPath),
			vmess.GetValue(field.XhttpMode),
			vmess.GetValue(field.XhttpExtra),
		)
	case "kcp", "mkcp":
		streamSettings["kcpSettings"] = genMkcpSetting(vmess.GetValue(field.KcpMtu))
	case "grpc":
		streamSettings["grpcSettings"] = genGrpcSetting(
			vmess.GetValue(field.GrpcMode),
			vmess.GetValue(field.GrpcServiceName),
			vmess.GetValue(field.GrpcAuthority),
		)
	case "ws":
		streamSettings["wsSettings"] = genWsSetting(vmess.GetValue(field.WsHost), vmess.GetValue(field.WsPath))
	case "h2":
		streamSettings["httpSettings"] = genHttpSetting(vmess.GetValue(field.H2Host), vmess.GetValue(field.H2Path))
	case "httpupgrade":
		streamSettings["httpupgradeSettings"] = genHttpupgradeSetting(
			vmess.GetValue(field.HttpUpgradeHost),
			vmess.GetValue(field.HttpUpgradePath),
		)
	}
	if vmess.GetValue(field.Finalmask) != "" {
		streamSettings["finalmask"] = genFinalmask(vmess.GetValue(field.Finalmask))
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"address":  vmess.Address,
			"port":     vmess.Port,
			"id":       vmess.ID,
			"security": vmess.GetValue(field.VMessEncryption),
			"level":    0,
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// VLESS
func vLessOutbound(vless *protocols.VLess) interface{} {
	mux := setting.Mux()
	security := vless.GetValue(field.TlsSecurity)
	network := vless.GetValue(field.NetworkType)
	switch network {
	case "tcp":
		network = "raw"
	case "kcp":
		network = "mkcp"
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		streamSettings["tlsSettings"] = genTlsSetting(
			vless.GetHostValue(field.SNI),
			vless.GetValue(field.Alpn),
			vless.GetValue(field.FingerPrint),
			setting.AllowInsecure(),
		)
	case "reality":
		streamSettings["realitySettings"] = genRealitySetting(
			vless.GetHostValue(field.SNI),
			vless.GetValue(field.FingerPrint),
			vless.GetValue(field.PublicKey),
			vless.GetValue(field.ShortId),
			vless.GetValue(field.SpiderX),
			vless.GetValue(field.Mldsa65Verify),
		)
		mux = false
	}
	switch network {
	case "tcp", "raw":
		streamSettings["rawSettings"] = genRawSetting(
			vless.GetValue(field.RawHeaderType),
			vless.GetValue(field.RawHost),
			vless.GetValue(field.RawPath),
		)
	case "xhttp":
		streamSettings["xhttpSettings"] = genXhttpSetting(
			vless.GetHostValue(field.XhttpHost),
			vless.GetValue(field.XhttpPath),
			vless.GetValue(field.XhttpMode),
			vless.GetValue(field.XhttpExtra),
		)
	case "kcp", "mkcp":
		streamSettings["kcpSettings"] = genMkcpSetting(vless.GetValue(field.KcpMtu))
	case "grpc":
		streamSettings["grpcSettings"] = genGrpcSetting(
			vless.GetValue(field.GrpcMode),
			vless.GetValue(field.GrpcServiceName),
			vless.GetValue(field.GrpcAuthority),
		)
	case "ws":
		streamSettings["wsSettings"] = genWsSetting(vless.GetValue(field.WsHost), vless.GetValue(field.WsPath))
	case "h2":
		streamSettings["httpSettings"] = genHttpSetting(vless.GetValue(field.H2Host), vless.GetValue(field.H2Path))
	case "httpupgrade":
		streamSettings["httpupgradeSettings"] = genHttpupgradeSetting(
			vless.GetValue(field.HttpUpgradeHost),
			vless.GetValue(field.HttpUpgradePath),
		)
	}
	if vless.GetValue(field.Finalmask) != "" {
		streamSettings["finalmask"] = genFinalmask(vless.GetValue(field.Finalmask))
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "vless",
		"settings": map[string]interface{}{
			"address":    vless.Address,
			"port":       vless.Port,
			"id":         vless.ID,
			"flow":       vless.GetValue(field.Flow),
			"encryption": vless.GetValue(field.VLessEncryption),
			"level":      0,
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}

// Hysteria2
func hysteria2Outbound(hysteria2 *protocols.Hysteria2) interface{} {
	mux := setting.Mux()
	network := "hysteria"
	security := "tls"
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		streamSettings["tlsSettings"] = genTlsSetting(
			hysteria2.GetHostValue(field.SNI),
			hysteria2.GetValue(field.Alpn),
			hysteria2.GetValue(field.FingerPrint),
			setting.AllowInsecure(),
		)
	}
	// hysteriaSettings
	streamSettings["hysteriaSettings"] = map[string]interface{}{
		"version": 2,
		"auth":    hysteria2.Password,
	}
	if hysteria2.GetValue(field.Mport) != "" {
		// finalmask
		streamSettings["finalmask"] = map[string]interface{}{
			"quicParams": map[string]interface{}{
				"congestion": "brutal",
				"brutalUp":   "100mbps",
				"brutalDown": "100mbps",
				"udpHop": map[string]interface{}{
					"ports":    hysteria2.GetValue(field.Mport),
					"interval": strings.TrimSuffix(hysteria2.GetValue(field.Hopinterval), "s"),
				},
			},
		}
	}
	if hysteria2.GetValue(field.Finalmask) != "" {
		streamSettings["finalmask"] = genFinalmask(hysteria2.GetValue(field.Finalmask))
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "hysteria",
		"settings": map[string]interface{}{
			"address": hysteria2.Address,
			"port":    hysteria2.Port,
			"version": 2,
		},
		"streamSettings": streamSettings,
		"mux": map[string]interface{}{
			"enabled": mux,
		},
	}
}
