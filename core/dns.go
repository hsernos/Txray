package core

import "Txray/log"

// 设置DNS监听端口
func (c *Core) SetDNSPort(port uint) {
	defer c.Save()
	c.DNS.Port = port
	log.Info("设置本地DNS监听端口为 [", port, "]")
}

// 设置境外DNS
func (c *Core) SetOutlandDNS(dns string) {
	defer c.Save()
	c.DNS.Outland = dns
	log.Info("设置境外DNS为 [", dns, "]")
}

// 设置境内DNS
func (c *Core) SetInlandDNS(dns string) {
	defer c.Save()
	c.DNS.Inland = dns
	log.Info("设置境内DNS为 [", dns, "]")
}

// 设置备用DNS端口（多个用英文逗号分隔）
func (c *Core) SetBackupDNS(dns string) {
	defer c.Save()
	c.DNS.Backup = dns
	log.Info("设置备用DNS为 [", dns, "]")
}
