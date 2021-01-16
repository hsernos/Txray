package core

import "Txray/log"

// 本地socks监听端口
func (c *Core) GetSocksPort() uint {
	return c.Settings.Socks
}

// 本地http监听端口
func (c *Core) GetHttpPort() uint {
	return c.Settings.Http
}

// 设置本地监听端口
func (c *Core) SetSocksPort(port uint) {
	defer c.Save()
	c.Settings.Socks = port
	log.Info("设置本地socks5监听端口为 [", port, "]")
}

// 设置本地监听端口
func (c *Core) SetHttpPort(port uint) {
	defer c.Save()
	c.Settings.Http = port
	log.Info("设置本地http监听端口为 [", port, "]")
}

// 设置是否允许udp转发
func (c *Core) SetUDP(udp bool) {
	defer c.Save()
	c.Settings.UDP = udp
	if udp {
		log.Info("开启UDP转发")
	} else {
		log.Info("关闭UDP转发")
	}

}

// 设置是否允许流量监听
func (c *Core) SetSniffing(sniffing bool) {
	defer c.Save()
	c.Settings.Sniffing = sniffing
	if sniffing {
		log.Info("开启流量监听")
	} else {
		log.Info("关闭流量监听")
	}
}

// 设置是否允许多路复用
func (c *Core) SetMux(mux bool) {
	defer c.Save()
	c.Settings.Mux = mux
	if mux {
		log.Info("开启多路复用")
	} else {
		log.Info("关闭多路复用")
	}
}

// 设置是否允许局域网的连接
func (c *Core) SetLANConn(conn bool) {
	defer c.Save()
	c.Settings.AllowLANConn = conn
	if conn {
		log.Info("允许来着局域网的连接")
	} else {
		log.Info("禁止来着局域网的连接")
	}
}

// 设置是否绕过局域网和大陆
func (c *Core) SetBypassLanAndContinent(bypass bool) {
	defer c.Save()
	c.Settings.BypassLanAndContinent = bypass
	log.Info("设置是否绕过局域网和大陆为: ", bypass)
	if bypass {
		log.Info("开启绕过局域网和大陆")
	} else {
		log.Info("关闭绕过局域网和大陆")
	}
}

// 设置路由策略
func (c *Core) SetDomainStrategy(mode int) {
	defer c.Save()
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
