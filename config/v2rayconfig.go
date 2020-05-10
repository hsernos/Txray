package config

import "v3ray/tool"

// GenConfig 生成v2ray-core配置文件
func (c *Config) GenConfig() {
	path := tool.Join(c.getConfigPath(), "config.json")
	var conf = map[string]interface{}{
		"log":       c.getLogConfig(),
		"inbounds":  c.getInboundsConfig(),
		"outbounds": c.getOutboundConfig(),
		"policy":    c.getPolicy(),
		"dns":       c.getDNSConfig(),
		"routing":   c.getRoutingConfig(),
	}
	tool.WriteJSON(conf, path)
}

func (c *Config) getLogConfig() interface{} {
	return map[string]string{
		"loglevel": "warning",
	}
}

func (c *Config) getInboundsConfig() interface{} {
	listen := "127.0.0.1"
	if c.Settings.AllowLANConn {
		listen = "0.0.0.0"
	}
	return []interface{}{
		map[string]interface{}{
			"tag":      "proxy",
			"port":     c.Settings.Port,
			"listen":   listen,
			"protocol": "socks",
			"sniffing": map[string]interface{}{
				"enabled": c.Settings.Sniffing,
				"destOverride": []string{
					"http",
					"tls",
				},
			},
			"settings": map[string]interface{}{
				"auth":      "noauth",
				"udp":       c.Settings.UDP,
				"userLevel": 0,
			},
			"streamSettings": nil,
		},
		map[string]interface{}{
			"tag":      "http",
			"port":     c.Settings.Http,
			"listen":   listen,
			"protocol": "http",
			"settings": map[string]interface{}{
				"userLevel": 0,
			},
			"streamSettings": nil,
		},
	}
}

func (c *Config) getOutboundConfig() interface{} {
	if int(c.Index) >= len(c.Nodes) {
		return nil
	}
	n := c.Nodes[c.Index]
	var tlsSettings interface{}
	var tcpSettings interface{}
	var kcpSettings interface{}
	var wsSettings interface{}
	var httpSettings interface{}
	var quicSettings interface{}
	if n.StreamSecurity == "tls" {
		tlsSettings = map[string]interface{}{
			"allowInsecure": true,
			"serverName":    n.RequestHost,
		}
	}

	switch n.Network {
	case "tcp":
		tcpSettings = nil
	case "kcp":
		kcpSettings = map[string]interface{}{
			"mtu":              c.KcpSetting.Mtu,
			"tti":              c.KcpSetting.Tti,
			"uplinkCapacity":   c.KcpSetting.UplinkCapacity,
			"downlinkCapacity": c.KcpSetting.DownlinkCapacity,
			"congestion":       c.KcpSetting.Congestion,
			"readBufferSize":   c.KcpSetting.ReadBufferSize,
			"writeBufferSize":  c.KcpSetting.WriteBufferSize,
			"header": map[string]interface{}{
				"type": "none",
			},
		}
	case "ws":
		wsSettings = map[string]interface{}{
			"connectionReuse": true,
			"path":            n.Path,
			"headers": map[string]interface{}{
				"Host": n.RequestHost,
			},
		}
	case "h2":
		httpSettings = map[string]interface{}{
			"path": n.Path,
			"host": []interface{}{
				n.RequestHost,
			},
		}
	case "quic":
		quicSettings = nil
	}

	return []interface{}{
		map[string]interface{}{
			"tag":      "proxy",
			"protocol": "vmess",
			"settings": map[string]interface{}{
				"vnext": []interface{}{
					map[string]interface{}{
						"address": n.Address,
						"port":    n.Port,
						"users": []interface{}{
							map[string]interface{}{
								"id":       n.ID,
								"alterId":  n.AlterID,
								"security": n.Security,
								"level":    0,
							},
						},
					},
				},
				"servers":  nil,
				"response": nil,
			},
			"streamSettings": map[string]interface{}{
				"network":      n.Network,
				"security":     n.StreamSecurity,
				"tlsSettings":  tlsSettings,
				"tcpSettings":  tcpSettings,
				"kcpSettings":  kcpSettings,
				"wsSettings":   wsSettings,
				"httpSettings": httpSettings,
				"quicSettings": quicSettings,
			},
			"mux": map[string]interface{}{
				"enabled": c.Settings.Mux,
			},
		},
		map[string]interface{}{
			"tag":      "direct",
			"protocol": "freedom",
			"settings": map[string]interface{}{
				"vnext":    nil,
				"servers":  nil,
				"response": nil,
			},
			"streamSettings": nil,
			"mux":            nil,
		},
		map[string]interface{}{
			"tag":      "block",
			"protocol": "blackhole",
			"settings": map[string]interface{}{
				"vnext":   nil,
				"servers": nil,
				"response": map[string]interface{}{
					"type": "http",
				},
			},
			"streamSettings": nil,
			"mux":            nil,
		},
	}
}

func (c *Config) getPolicy() interface{} {
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

func (c *Config) getDNSConfig() interface{} {
	return map[string]interface{}{
		"servers": c.DNS,
	}
}

func (c *Config) getRoutingConfig() interface{} {
	rules := make([]interface{}, 0, 9)
	if c.Settings.BypassLanAndContinent {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "direct",
			"ip": []string{
				"geoip:private",
				"geoip:cn",
			},
			"domain": nil,
		})
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "direct",
			"ip":          nil,
			"domain": []string{
				"geosite:cn",
			},
		})
	}

	if len(c.Direct.IP) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "direct",
			"ip":          c.Direct.IP,
			"domain":      nil,
		})
	}
	if len(c.Direct.Domain) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "direct",
			"ip":          nil,
			"domain":      c.Direct.Domain,
		})
	}
	if len(c.Proxy.IP) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "proxy",
			"ip":          c.Proxy.IP,
			"domain":      nil,
		})
	}
	if len(c.Proxy.Domain) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "proxy",
			"ip":          nil,
			"domain":      c.Proxy.Domain,
		})
	}
	if len(c.Block.IP) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "block",
			"ip":          c.Block.IP,
			"domain":      nil,
		})
	}
	if len(c.Block.Domain) != 0 {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        nil,
			"inboundTag":  nil,
			"outboundTag": "block",
			"ip":          nil,
			"domain":      c.Block.Domain,
		})
	}

	return map[string]interface{}{
		"domainStrategy": c.Settings.DomainStrategy,
		"rules":          rules,
	}
}
