package config

import log "v3ray/logger"

// SetPort 设置本地监听端口
func (c *Config) SetPort(port uint) {
	defer c.SaveJSON()
	c.Settings.Port = port
	log.Info("设置本地监听端口为 [", port, "]")
}

// SetUDP 设置是否允许udp转发
func (c *Config) SetUDP(udp bool) {
	defer c.SaveJSON()
	c.Settings.UDP = udp
	if udp {
		log.Info("开启UDP转发")
	} else {
		log.Info("关闭UDP转发")
	}

}

// SetSniffing 设置是否允许流量监听
func (c *Config) SetSniffing(sniffing bool) {
	defer c.SaveJSON()
	c.Settings.Sniffing = sniffing
	if sniffing {
		log.Info("开启流量监听")
	} else {
		log.Info("关闭流量监听")
	}
}

// SetMux 设置是否允许多路复用
func (c *Config) SetMux(mux bool) {
	defer c.SaveJSON()
	c.Settings.Mux = mux
	if mux {
		log.Info("开启多路复用")
	} else {
		log.Info("关闭多路复用")
	}
}

// SetLANConn 设置是否允许局域网的连接
func (c *Config) SetLANConn(conn bool) {
	defer c.SaveJSON()
	c.Settings.AllowLANConn = conn
	if conn {
		log.Info("允许来着局域网的连接")
	} else {
		log.Info("禁止来着局域网的连接")
	}
}

// SetBypassLanAndContinent 设置是否绕过局域网和大陆
func (c *Config) SetBypassLanAndContinent(bypass bool) {
	defer c.SaveJSON()
	c.Settings.BypassLanAndContinent = bypass
	log.Info("设置是否绕过局域网和大陆为: ", bypass)
	if bypass {
		log.Info("开启绕过局域网和大陆")
	} else {
		log.Info("关闭绕过局域网和大陆")
	}
}

// SetDomainStrategy 设置路由策略
func (c *Config) SetDomainStrategy(mode int) {
	defer c.SaveJSON()
	domainStrategy := c.Settings.DomainStrategy
	switch mode {
	case 1:
		domainStrategy = "AsIs"
	case 2:
		domainStrategy = "IPIfNonMatch"
	case 3:
		domainStrategy = "IPOnDemand"
	}
	c.Settings.DomainStrategy = domainStrategy
	log.Info("设置路由策略为 [", domainStrategy, "]")
}