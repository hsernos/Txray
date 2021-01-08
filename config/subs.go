package config

import (
	log "Tv2ray/logger"
	"Tv2ray/tools"
	"Tv2ray/vmess"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// 添加一条订阅记录
func (c *Config) AddSub(url, remarks string) {
	defer c.SaveJSON()
	u1 := uuid.NewV4()
	s := sub{u1.String(), remarks, url, true}
	c.Subs = append(c.Subs, &s)
	log.Info("添加订阅: ", url)
}

// 通过订阅添加节点
func (c *Config) AddNodeBySub(key string, port uint) {
	defer c.SaveJSON()
	indexs := make([]int, 0, len(c.Subs))
	if key == "" {
		for i, x := range c.Subs {
			if x.Using {
				indexs = append(indexs, i)
			}
		}
	} else {
		indexs = tools.IndexDeal(key, len(c.Subs))
	}

	newNodes := make([]*node, 0, len(c.Nodes))
	ids := ""
	for _, i := range indexs {
		x := c.Subs[i]
		data := ""
		if port > 65535 {
			res, e := tools.GetNoProxy(x.URL, 10)

			if e != nil {
				log.Warn(e)
				break
			}
			data = tools.ReadDate(res)
			log.Info("访问 [", x.URL, "] -- ", res.Status)
			res.Body.Close()
		} else {
			res, e := tools.GetBySocks5Proxy(x.URL, "127.0.0.1", c.Settings.Port, 10)
			if e != nil {
				log.Warn(e)
				break
			}
			data = tools.ReadDate(res)
			log.Info("访问 [", x.URL, "] -- ", res.Status)
			res.Body.Close()
		}
		vmessList := vmess.Sub2links(data)
		Objs := vmess.Links2vmessObjs(vmessList)
		for _, obj := range Objs {
			newNodes = append(newNodes, vmessObjToNode(obj, x.ID))
		}
		if len(Objs) != 0 {
			ids += x.ID + ","
		}
		log.Info("从订阅 [", x.URL, "] 获取了 [", len(Objs), "] 个节点")
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

// 获取订阅信息
func (c *Config) GetSubs(key string) [][]string {
	l := len(c.Subs)
	indexs := tools.IndexDeal(key, l)
	result := make([][]string, 0, len(indexs))
	for _, x := range indexs {
		sub := c.Subs[x]
		result = append(result, []string{
			tools.IntToStr(x + 1),
			sub.Remarks,
			sub.URL,
			tools.BoolToStr(sub.Using),
		})
	}
	return result
}

// 修改订阅信息
func (c *Config) SetSubs(key, using, url, remarks string) {
	l := len(c.Subs)
	indexs := tools.IndexDeal(key, l)
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

// 删除订阅信息
func (c *Config) DelSubs(key string) {
	l := len(c.Subs)
	indexs := tools.IndexDeal(key, l)
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
