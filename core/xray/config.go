package xray

import (
	"Txray/core/config"
	"Txray/core/node"
	"Txray/core/protocols"
	"Txray/core/routing"
	"Txray/core/setting"
	"Txray/log"
	"Txray/tools"
)

// 生成xray-core配置文件
func GenConfig(index int) string {
	path := tools.PathJoin(config.GetConfigDir(), "config.json")
	var conf = map[string]interface{}{
		"log":       logConfig(),
		"inbounds":  inboundsConfig(),
		"outbounds": outboundConfig(index),
		"policy":    policyConfig(),
		"dns":       dnsConfig(),
		"routing":   routingConfig(),
	}
	err := tools.WriteJSON(conf, path)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return path
}

// 日志
func logConfig() interface{} {
	path := tools.PathJoin(config.GetConfigDir(), "xray_access.log")
	return map[string]string{
		"access":   path,
		"loglevel": "warning",
	}
}

// 入站
func inboundsConfig() interface{} {
	listen := "127.0.0.1"
	base := setting.Base()
	if base.AllowLANConn {
		listen = "0.0.0.0"
	}
	data := []interface{}{
		map[string]interface{}{
			"tag":      "proxy",
			"port":     base.Socks,
			"listen":   listen,
			"protocol": "socks",
			"sniffing": map[string]interface{}{
				"enabled": base.Sniffing,
				"destOverride": []string{
					"http",
					"tls",
				},
			},
			"settings": map[string]interface{}{
				"auth":      "noauth",
				"udp":       base.UDP,
				"userLevel": 0,
			},
		},
	}
	if base.Http > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "http",
			"port":     base.Http,
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
				"address":   setting.OutlandDNS(),
				"network":   "tcp,udp",
				"port":      53,
			},
		})
	}
	return data
}

// 出站
func outboundConfig(index int) interface{} {
	out := make([]interface{}, 0)
	n := node.GetNode(index)
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
					"port":     tools.StrToInt(ss.Port),
					"password": ss.Password,
					"method":   ss.Method,
					"level":    0,
				},
			},
		},
		"streamSettings": map[string]interface{}{
			"network": "tcp",
		},
		"mux": map[string]interface{}{
			"enabled": setting.Base().Mux,
		},
	}
}

// Trojan
func trojanOutbound(trojan *protocols.Trojan) interface{} {
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "trojan",
		"settings": map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  trojan.Address,
					"port":     tools.StrToInt(trojan.Port),
					"password": trojan.Password,
					"level":    0,
				},
			},
		},
		"streamSettings": map[string]interface{}{
			"network":  "tcp",
			"security": "tls",
			"tlsSettings": map[string]interface{}{
				"allowInsecure": false,
				"serverName":    "",
			},
		},
		"mux": map[string]interface{}{
			"enabled": setting.Base().Mux,
		},
	}
}

// VMess
func vMessOutbound(vmess *protocols.VMess) interface{} {
	var tlsSettings interface{}
	var tcpSettings interface{}
	var kcpSettings interface{}
	var wsSettings interface{}
	var httpSettings interface{}
	var quicSettings interface{}
	if vmess.Tls == "tls" {
		tlsSettings = map[string]interface{}{
			"allowInsecure": true,
			"serverName":    vmess.Host,
		}
	}
	switch vmess.Net {
	case "tcp":
		tcpSettings = nil
	case "kcp":
		kcp := setting.KCP()
		kcpSettings = map[string]interface{}{
			"mtu":              kcp.Mtu,
			"tti":              kcp.Tti,
			"uplinkCapacity":   kcp.UplinkCapacity,
			"downlinkCapacity": kcp.DownlinkCapacity,
			"congestion":       kcp.Congestion,
			"readBufferSize":   kcp.ReadBufferSize,
			"writeBufferSize":  kcp.WriteBufferSize,
			"header": map[string]interface{}{
				"type": "none",
			},
		}
	case "ws":
		wsSettings = map[string]interface{}{
			"connectionReuse": true,
			"path":            vmess.Path,
			"headers": map[string]interface{}{
				"Host": vmess.Host,
			},
		}
	case "h2":
		httpSettings = map[string]interface{}{
			"path": vmess.Path,
			"host": []interface{}{
				vmess.Host,
			},
		}
	case "quic":
		quicSettings = nil
	}
	return map[string]interface{}{
		"tag":      "proxy",
		"protocol": "vmess",
		"settings": map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": vmess.Add,
					"port":    tools.StrToInt(vmess.Port),
					"users": []interface{}{
						map[string]interface{}{
							"id":       vmess.Id,
							"alterId":  tools.StrToInt(vmess.Aid),
							"security": "auto",
							"level":    0,
						},
					},
				},
			},
		},
		"streamSettings": map[string]interface{}{
			"network":      vmess.Net,
			"security":     vmess.Tls,
			"tlsSettings":  tlsSettings,
			"tcpSettings":  tcpSettings,
			"kcpSettings":  kcpSettings,
			"wsSettings":   wsSettings,
			"httpSettings": httpSettings,
			"quicSettings": quicSettings,
		},
		"mux": map[string]interface{}{
			"enabled": setting.Base().Mux,
		},
	}
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
	servers := make([]interface{}, 0, 0)
	if setting.InlandDNS() != "" {
		servers = append(servers, map[string]interface{}{
			"address": setting.InlandDNS(),
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if setting.OutlandDNS() != "" {
		servers = append(servers, map[string]interface{}{
			"address": setting.OutlandDNS(),
			"port":    53,
			"domains": []interface{}{
				"geosite:geolocation-!cn",
				"geosite:speedtest",
			},
		})
	}
	if len(setting.BackupDNS()) != 0 && setting.BackupDNS()[0] != "" {
		for _, bk := range setting.BackupDNS() {
			servers = append(servers, bk)
		}
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
	rules := make([]interface{}, 0, 0)
	if setting.DNSPort() != 0 {
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"inboundTag": []interface{}{
				"dns-in",
			},
			"outboundTag": "dns-out",
		})
	}
	if setting.OutlandDNS() != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "proxy",
			"ip": []string{
				setting.OutlandDNS(),
			},
		})
	}
	if setting.InlandDNS() != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "direct",
			"ip": []string{
				setting.InlandDNS(),
			},
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
	base := setting.Base()
	if base.BypassLanAndContinent {
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
		"domainStrategy": base.DomainStrategy,
		"rules":          rules,
	}
}
