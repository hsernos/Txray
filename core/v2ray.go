package core

import (
	"Txray/core/protocols"
	"Txray/core/protocols/mode"
	"Txray/log"
	"Txray/tools"
	"path/filepath"
	"strings"
)

// 检查v2ray程序和资源文件 geoip.dat、geosite.dat是否存在且在同一目录
func CheckV2rayFile() bool {
	v2rayPath := GetXrayPath()
	if v2rayPath == "" {
		log.Error("在", tools.GetRunPath(), " 下没有找到v2ray程序")
		log.Error("请在 https://github.com/v2fly/v2ray-core/releases 下载最新版本")
		log.Error("并将解压后的文件夹或所有文件移动到 ", tools.GetRunPath(), " 下")
		return false
	} else {
		path := filepath.Dir(GetXrayPath())
		if tools.IsFile(tools.PathJoin(path, "geoip.dat")) && tools.IsFile(tools.PathJoin(path, "geosite.dat")) {
			return true
		} else {
			log.Error("在 ", path, " 下没有找到v2ray程序的资源文件 geoip.dat 和 geosite.dat")
			log.Error("请在 https://github.com/v2fly/v2ray-core/releases 下载最新版本")
			log.Error("并将缺失的文件移动到 ", path, " 下")
			return false
		}
	}
}

// 查找v2ray程序所在绝对路径
func GetV2rayPath() string {
	path := tools.GetRunPath()
	files, _ := tools.FindFileByName(path, "v2ray", ".exe")
	if len(files) == 0 {
		return ""
	}
	return files[0]
}

// 生成v2ray-core配置文件
func (c *Core) GenV2rayConfig() string {
	path := tools.PathJoin(tools.GetRunPath(), "config.json")
	var conf = map[string]interface{}{
		"log":       c.getV2rayLogConfig(),
		"inbounds":  c.getV2rayInboundsConfig(),
		"outbounds": c.getV2rayOutboundConfig(),
		"policy":    c.getV2rayPolicy(),
		"dns":       c.getV2rayDNSConfig(),
		"routing":   c.getV2rayRoutingConfig(),
	}
	err := tools.WriteJSON(conf, path)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return path
}

// 日志
func (c *Core) getV2rayLogConfig() interface{} {
	path := tools.PathJoin(tools.GetRunPath(), "v2ray_access.log")
	return map[string]string{
		"access":   path,
		"loglevel": "warning",
	}
}

// 入站
func (c *Core) getV2rayInboundsConfig() interface{} {
	listen := "127.0.0.1"
	if c.Settings.AllowLANConn {
		listen = "0.0.0.0"
	}
	data := []interface{}{
		map[string]interface{}{
			"tag":      "proxy",
			"port":     c.Settings.Socks,
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
		},
	}
	if c.Settings.Http > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "http",
			"port":     c.Settings.Http,
			"listen":   listen,
			"protocol": "http",
			"settings": map[string]interface{}{
				"userLevel": 0,
			},
		})
	}
	if c.DNS.Port > 0 {
		data = append(data, map[string]interface{}{
			"tag":      "dns-in",
			"port":     c.DNS.Port,
			"listen":   listen,
			"protocol": "dokodemo-door",
			"settings": map[string]interface{}{
				"userLevel": 0,
				"address":   c.DNS.Outland,
				"network":   "tcp,udp",
				"port":      53,
			},
		})
	}
	return data
}

// 出站
func (c *Core) getV2rayOutboundConfig() interface{} {
	out := make([]interface{}, 0)
	switch c.GetNodeMode() {
	case mode.Trojan:
		out = append(out, c.getV2rayTrojanOutbound())
	case mode.ShadowSocks:
		out = append(out, c.getV2raySSOutbound())
	case mode.VMess:
		out = append(out, c.getV2rayVMessOutbound())
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
func (c *Core) getV2raySSOutbound() interface{} {
	ss := new(protocols.ShadowSocks)
	ss.ParseLink(c.GetNodeLink())
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
			"enabled": c.Settings.Mux,
		},
	}
}

// Trojan
func (c *Core) getV2rayTrojanOutbound() interface{} {
	trojan := new(protocols.Trojan)
	trojan.ParseLink(c.GetNodeLink())
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
				//"allowInsecure": false,
				"serverName": trojan.Address,
			},
		},
		"mux": map[string]interface{}{
			"enabled": c.Settings.Mux,
		},
	}
}

// VMess
func (c *Core) getV2rayVMessOutbound() interface{} {
	vmess := new(protocols.VMess)
	vmess.ParseLink(c.GetNodeLink())
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
			"enabled": c.Settings.Mux,
		},
	}
}

// 本地策略
func (c *Core) getV2rayPolicy() interface{} {
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
func (c *Core) getV2rayDNSConfig() interface{} {
	servers := make([]interface{}, 0, 0)
	if c.DNS.Inland != "" {
		servers = append(servers, map[string]interface{}{
			"address": c.DNS.Inland,
			"port":    53,
			"domains": []interface{}{
				"geosite:cn",
			},
			"expectIPs": []interface{}{
				"geoip:cn",
			},
		})
	}
	if c.DNS.Outland != "" {
		servers = append(servers, map[string]interface{}{
			"address": c.DNS.Outland,
			"port":    53,
			"domains": []interface{}{
				"geosite:geolocation-!cn",
				"geosite:speedtest",
			},
		})
	}
	if c.DNS.Backup != "" {
		for _, bk := range strings.Split(c.DNS.Backup, ",") {
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
func (c *Core) getV2rayRoutingConfig() interface{} {
	rules := make([]interface{}, 0, 0)
	if c.DNS.Port != 0 {
		rules = append(rules, map[string]interface{}{
			"type": "field",
			"inboundTag": []interface{}{
				"dns-in",
			},
			"outboundTag": "dns-out",
		})
	}
	if c.DNS.Outland != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "proxy",
			"ip": []string{
				c.DNS.Outland,
			},
		})
	}
	if c.DNS.Inland != "" {
		rules = append(rules, map[string]interface{}{
			"type":        "field",
			"port":        53,
			"outboundTag": "direct",
			"ip": []string{
				c.DNS.Inland,
			},
		})
	}
	ips, domains := GetRouteGroupData(&c.Block)
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
	ips, domains = GetRouteGroupData(&c.Proxy)
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
	ips, domains = GetRouteGroupData(&c.Direct)
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
	if c.Settings.BypassLanAndContinent {
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
		"domainStrategy": c.Settings.DomainStrategy,
		"rules":          rules,
	}
}
