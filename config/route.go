package config

import "Tv2ray/tools"
import log "Tv2ray/logger"

// 添加一条直连IP规则
func (c *Config) AddDirectIP(data string) {
	defer c.SaveJSON()
	c.Direct.IP = append(c.Direct.IP, data)
	log.Info("添加一条直连规则 [", data, "]")
}

// 添加一条直连Domain规则
func (c *Config) AddDirectDomain(data string) {
	defer c.SaveJSON()
	c.Direct.Domain = append(c.Direct.Domain, data)
	log.Info("添加一条直连规则 [", data, "]")
}

// 添加一条代理IP规则
func (c *Config) AddProxyIP(data string) {
	defer c.SaveJSON()
	c.Proxy.IP = append(c.Proxy.IP, data)
	log.Info("添加一条代理规则 [", data, "]")
}

// 添加一条代理Domain规则
func (c *Config) AddProxyDomain(data string) {
	defer c.SaveJSON()
	c.Proxy.Domain = append(c.Proxy.Domain, data)
	log.Info("添加一条代理规则 [", data, "]")
}

// 添加一条禁止IP规则
func (c *Config) AddBlockIP(data string) {

	defer c.SaveJSON()
	c.Block.IP = append(c.Block.IP, data)
	log.Info("添加一条禁止规则 [", data, "]")
}

// 添加一条禁止Domain规则
func (c *Config) AddBlockDomain(data string) {
	defer c.SaveJSON()
	c.Block.Domain = append(c.Block.Domain, data)
	log.Info("添加一条禁止规则 [", data, "]")
}

// 获取直连IP规则
func (c *Config) GetDirectIP() [][]string {
	result := make([][]string, 0, len(c.Direct.IP))
	for i, x := range c.Direct.IP {
		result = append(result, []string{tools.IntToStr(i + 1), x})
	}
	return result
}

// 获取直连Dmain规则
func (c *Config) GetDirectDomain() [][]string {
	result := make([][]string, 0, len(c.Direct.Domain))
	for i, x := range c.Direct.Domain {
		result = append(result, []string{tools.IntToStr(i + 1), x})
	}
	return result
}

// 获取代理IP规则
func (c *Config) GetProxyIP() [][]string {
	result := make([][]string, 0, len(c.Proxy.IP))
	for i, x := range c.Proxy.IP {
		result = append(result, []string{tools.IntToStr(i + 1), x})
	}
	return result
}

// 获取代理Domian规则
func (c *Config) GetProxyDomain() [][]string {
	result := make([][]string, 0, len(c.Proxy.Domain))
	for i, x := range c.Proxy.Domain {
		result = append(result, []string{tools.IntToStr(i + 1), x})
	}
	return result
}

// 获取禁止IP规则
func (c *Config) GetBlockIP() [][]string {
	result := make([][]string, 0, len(c.Block.IP))
	for i, x := range c.Block.IP {
		result = append(result, []string{tools.IntToStr(i + 1), x})
	}
	return result
}

// 获取禁止Domain规则
func (c *Config) GetBlockDomain() [][]string {
	result := make([][]string, 0, len(c.Block.Domain))
	for i, x := range c.Block.Domain {
		result = append(result, []string{tools.IntToStr(i + 1), x})
	}
	return result
}

// 删除直连IP规则
func (c *Config) DelDirectIP(key string) {

	l := len(c.Direct.IP)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.Direct.IP {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Direct.IP = result
	log.Info("删除了 [", l-len(c.Direct.IP), "] 条")
}

// 删除直连Domain规则
func (c *Config) DelDirectDomain(key string) {
	l := len(c.Direct.Domain)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.Direct.Domain {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Direct.Domain = result
	log.Info("删除了 [", l-len(c.Direct.Domain), "] 条")
}

// 删除代理IP规则
func (c *Config) DelProxyIP(key string) {
	l := len(c.Proxy.IP)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.Proxy.IP {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Proxy.IP = result
	log.Info("删除了 [", l-len(c.Proxy.IP), "] 条")
}

// 删除代理Domain规则
func (c *Config) DelProxyDomain(key string) {
	l := len(c.Proxy.Domain)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.Proxy.Domain {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Proxy.Domain = result
	log.Info("删除了 [", l-len(c.Proxy.Domain), "] 条")
}

// 删除禁止IP规则
func (c *Config) DelBlockIP(key string) {
	l := len(c.Block.IP)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.Block.IP {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Block.IP = result
	log.Info("删除了 [", l-len(c.Block.IP), "] 条")
}

// 删除禁止Domain规则
func (c *Config) DelBlockDomain(key string) {
	l := len(c.Block.Domain)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.Block.Domain {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Block.Domain = result
	log.Info("删除了 [", l-len(c.Block.Domain), "] 条")
}
