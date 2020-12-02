package config

import (
	log "Tv2ray/logger"
	"Tv2ray/tools"
)

// AddDNS 添加一条dns记录
func (c *Config) AddDNS(dns string) {

	defer c.SaveJSON()
	c.DNS = append(c.DNS, dns)
	log.Info("添加一条DNS [", dns, "]")
}

// GetDNS 获取dns记录
func (c *Config) GetDNS() [][]string {
	result := make([][]string, 0, len(c.DNS))
	for i, x := range c.DNS {
		result = append(result, []string{tools.IntToStr(i), x})
	}
	return result
}

// DelDNS 删除DNS
func (c *Config) DelDNS(key string) {
	l := len(c.DNS)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]string, 0, l-len(indexs))
	j := 0
	for i, y := range c.DNS {
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
	c.DNS = result
	log.Info("删除了 [", l-len(c.DNS), "] 条")
}
