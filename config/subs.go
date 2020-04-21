package config

import (
	uuid "github.com/satori/go.uuid"
	"strings"
	log "v3ray/logger"
	"v3ray/tool"
)

// AddSub 添加一条订阅记录
func (c *Config) AddSub(url, remarks string) {
	defer c.SaveJSON()
	u1 := uuid.NewV4()
	s := sub{u1.String(), remarks, url, true}
	c.Subs = append(c.Subs, &s)
	log.Info("添加订阅: ", url)
}

// AddNodeBySub 通过订阅添加节点
func (c *Config) AddNodeBySub(port uint) {
	defer c.SaveJSON()
	newNodes := make([]*node, 0, len(c.Nodes))
	ids := ""
	for _, x := range c.Subs {
		if x.Using {
			vmessList := tool.SubToVmessList(x.URL, port)
			Objs := tool.VmessListToObj(vmessList)
			for _, obj := range Objs {
				newNodes = append(newNodes, vmessObjToNode(obj, x.ID))
			}
			ids += x.ID + ","
			log.Info("从订阅 [", x.URL, "] 更新了 [", len(Objs), "] 个节点")
		}
	}
	log.Info("订阅更新完成，共更新 [", len(newNodes), "] 个节点")
	for _, x := range c.Nodes {
		if x.Subid == "" || !strings.Contains(ids, x.Subid) {
			newNodes = append(newNodes, x)
		}
	}
	c.Nodes = newNodes
	c.Index = 0
}

// GetSubs 获取订阅信息
func (c *Config) GetSubs(key string) [][]string {
	l := len(c.Subs)
	indexs := tool.IndexDeal(key, l)
	result := make([][]string, 0, len(indexs))
	for _, x := range indexs {
		sub := c.Subs[x]
		result = append(result, []string{
			tool.IntToStr(x),
			sub.Remarks,
			sub.URL,
			tool.BoolToStr(sub.Using),
		})
	}
	return result
}

// SetSubs 修改订阅信息
func (c *Config) SetSubs(key, using, url, remarks string) {
	l := len(c.Subs)
	indexs := tool.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	for _, x := range indexs {
		s := c.Subs[x]
		if using != "" {
			s.Using = using == "true"
		}
		if url != "" {
			s.URL = url
		}
		if remarks != "" {
			s.Remarks = remarks
		}
	}

}

// DelSubs 删除订阅信息
func (c *Config) DelSubs(key string) {
	l := len(c.Subs)
	indexs := tool.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]*sub, 0, l-len(indexs))
	j := 0
	for i, y := range c.Subs {
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
	c.Subs = result
	log.Info("删除了 [", l-len(c.Subs), "] 条")
}
