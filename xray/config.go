package xray

import (
	"Txray/core"
	"Txray/core/protocols"
	"Txray/core/protocols/field"
	"Txray/core/routing"
	"Txray/core/setting"
	"Txray/log"
	"strings"
)

// 生成xray-core配置文件
func GenConfig(node protocols.Protocol) string {
	path := PathJoin(core.GetConfigDir(), "config.json")
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
	path := PathJoin(core.GetConfigDir(), "xray_access.log")
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

//Shadowsocks
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
			"allowInsecure": false,
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

// VMess
func vMessOutbound(vmess *protocols.VMess) interface{} {
	mux := setting.Mux()
	streamSettings := map[string]interface{}{
		"network":  vmess.Net,
		"security": vmess.Tls,
	}
	if vmess.Tls == "tls" {
		streamSettings["tlsSettings"] = map[string]interface{}{
			"allowInsecure": false,
			"serverName":    vmess.Host,
		}
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
		streamSettings["quicSettings"] = map[string]interface{}{
			"security": vmess.Host,
			"key":      vmess.Path,
			"header": map[string]interface{}{
				"type": vmess.Type,
			},
		}
	case "grpc":
		streamSettings["grpcSettings"] = map[string]interface{}{
			"serviceName": vmess.Path,
			"multiMode":   vmess.Type == "multi",
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
							"security": "auto",
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
		user["users"] = map[string]interface{}{
			"user": socks.Username,
			"pass": socks.Password,
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
		"encryption": vless.GetValue(field.VLessEncryption),
		"level":      0,
	}
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	switch security {
	case "tls":
		streamSettings["tlsSettings"] = map[string]interface{}{
			"allowInsecure": false,
			"serverName":    vless.GetHostValue(field.SNI),
		}
	case "xtls":
		streamSettings["xtlsSettings"] = map[string]interface{}{
			"allowInsecure": false,
			"serverName":    vless.GetHostValue(field.SNI),
		}
		user["flow"] = vless.GetValue(field.Flow)
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
		streamSettings["tlsSettings"] = map[string]interface{}{
			"allowInsecure": true,
			"serverName":    vmess.GetHostValue(field.SNI),
		}
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
