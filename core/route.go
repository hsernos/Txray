package core

import (
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"regexp"
	"strings"
)

// 添加一条直连规则
func (c *Core) AddDirect(list ...string) {
	defer c.Save()
	num := 0
	for _, data := range list {
		if data != "" {
			mode := IsIpORDomain(data)
			if mode == "Unknown" {
				log.Warnf("%s: 不是IP或Domain规则", data)
			} else {
				num++
				r := new(routing)
				r.Data = data
				r.Mode = mode
				c.Direct = append(c.Direct, r)
				if len(list) == 1 {
					log.Infof("添加一条直连规则 [%s] %s", mode, data)
				}
			}
		}
	}
	if len(list) > 1 {
		log.Infof("共添加了 %d 条直连规则", num)
	}
}

// 添加一条代理规则
func (c *Core) AddProxy(list ...string) {
	defer c.Save()
	num := 0
	for _, data := range list {
		if data != "" {
			mode := IsIpORDomain(data)
			if mode == "Unknown" {
				log.Warnf("%s: 不是IP或Domain规则", data)
			} else {
				num++
				r := new(routing)
				r.Data = data
				r.Mode = mode
				c.Proxy = append(c.Proxy, r)
				if len(list) == 1 {
					log.Infof("添加一条代理规则 [%s] %s", mode, data)
				}
			}
		}
	}
	if len(list) > 1 {
		log.Infof("共添加了 %d 条代理规则", num)
	}
}

// 添加一条禁止规则
func (c *Core) AddBlock(list ...string) {
	defer c.Save()
	num := 0
	for _, data := range list {
		if data != "" {
			mode := IsIpORDomain(data)
			if mode == "Unknown" {
				log.Warnf("%s: 不是IP或Domain规则", data)
			} else {
				num++
				r := new(routing)
				r.Data = data
				r.Mode = mode
				c.Block = append(c.Block, r)
				if len(list) == 1 {
					log.Infof("添加一条禁止规则 [%s] %s", mode, data)
				}
			}
		}
	}
	if len(list) > 1 {
		log.Infof("共添加了 %d 条禁止规则", num)
	}
}

// 获取直连规则
func (c *Core) GetDirect(key string) [][]string {
	indexList := format.IndexDeal(key, len(c.Direct))
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		r := c.Direct[x]
		result = append(result, []string{
			tools.IntToStr(x + 1),
			r.Mode,
			r.Data,
		})
	}
	return result
}

// 对路由数据进行分组
func GetRouteGroupData(routes *[]*routing) ([]string, []string) {
	ips := make([]string, 0)
	domains := make([]string, 0)
	for _, x := range *routes {
		if x.Mode == "Domain" {
			domains = append(domains, x.Data)
		} else {
			ips = append(ips, x.Data)
		}
	}
	return ips, domains
}

// 获取代理规则
func (c *Core) GetProxy(key string) [][]string {
	indexList := format.IndexDeal(key, len(c.Proxy))
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		r := c.Proxy[x]
		result = append(result, []string{
			tools.IntToStr(x + 1),
			r.Mode,
			r.Data,
		})
	}
	return result
}

// 获取禁止规则
func (c *Core) GetBlock(key string) [][]string {
	indexList := format.IndexDeal(key, len(c.Block))
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		r := c.Block[x]
		result = append(result, []string{
			tools.IntToStr(x + 1),
			r.Mode,
			r.Data,
		})
	}
	return result
}

// 删除直连规则
func (c *Core) DelDirect(key string) {
	length := len(c.Direct)
	indexList := format.IndexDeal(key, length)
	if len(indexList) == 0 {
		return
	}
	defer c.Save()
	result := make([]*routing, 0, length-len(indexList))
	j := 0
	for i, y := range c.Direct {
		if j < len(indexList) {
			if i == indexList[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Direct = result
	log.Info("删除了 [", length-len(result), "] 条规则")
}

// 删除代理规则
func (c *Core) DelProxy(key string) {
	length := len(c.Proxy)
	indexList := format.IndexDeal(key, length)
	if len(indexList) == 0 {
		return
	}
	defer c.Save()
	result := make([]*routing, 0, length-len(indexList))
	j := 0
	for i, y := range c.Proxy {
		if j < len(indexList) {
			if i == indexList[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Proxy = result
	log.Info("删除了 [", length-len(result), "] 条规则")
}

// 删除禁止规则
func (c *Core) DelBlock(key string) {
	length := len(c.Block)
	indexList := format.IndexDeal(key, length)
	if len(indexList) == 0 {
		return
	}
	defer c.Save()
	result := make([]*routing, 0, length-len(indexList))
	j := 0
	for i, y := range c.Block {
		if j < len(indexList) {
			if i == indexList[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Block = result
	log.Info("删除了 [", length-len(result), "] 条规则")
}

// 判断是IP规则还是域名规则
// IP|Unknown|Domain
func IsIpORDomain(str string) string {
	if strings.HasPrefix(str, "regexp:") {
		return "Domain"
	}
	if strings.HasPrefix(str, "domain:") {
		return "Domain"
	}
	if strings.HasPrefix(str, "fill:") {
		return "Domain"
	}
	if strings.HasPrefix(str, "geosite:") {
		return "Domain"
	}
	if strings.HasPrefix(str, "geosite:") {
		return "Domain"
	}
	if strings.HasPrefix(str, "geoip:") {
		return "IP"
	}
	pattern := "^ext:[a-zA-Z0-9_]*?(ip|site).dat:"
	re, _ := regexp.Compile(pattern)
	result := re.FindStringSubmatch(str)
	if len(result) == 2 {
		if result[1] == "ip" {
			return "IP"
		} else {
			return "Domain"
		}
	}
	pattern = `(^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(/([1-9]|[1-2][0-9]|3[0-2]){1})?$`
	re, _ = regexp.Compile(pattern)
	result = re.FindStringSubmatch(str)
	if len(result) != 0 {
		return "IP"
	}
	pattern = `^([a-zA-Z0-9][-a-zA-Z0-9]{0,62})(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`
	re, _ = regexp.Compile(pattern)
	result = re.FindStringSubmatch(str)
	if len(result) != 0 && len(str) < 256 {
		return "Domain"
	}
	return "Unknown"
}
