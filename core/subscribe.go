package core

import (
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

// 添加一条订阅记录
func (c *Core) AddSub(url, remarks string) {
	defer c.Save()
	u1 := uuid.NewV4()
	s := sub{u1.String(), remarks, url, true}
	c.Subs = append(c.Subs, &s)
	log.Info("添加订阅: ", url)
}

// 更新订阅节点
// mode取值为socks5、http、none表示通过socks5、http代理和不使用代理
func (c *Core) AddNodeBySub(key string, mode string, address string, port uint) {
	defer c.Save()
	indexList := make([]int, 0, len(c.Subs))
	// 根据是否启用订阅来更新节点
	if key == "" {
		for i, x := range c.Subs {
			if x.Using {
				indexList = append(indexList, i)
			}
		}
		// 根据key选择订阅来更新节点（忽视是否启用）
	} else {
		indexList = format.IndexDeal(key, len(c.Subs))
	}

	newNodes := make([]*node, 0, len(c.Nodes))
	ids := ""
	for _, i := range indexList {
		subscribe := c.Subs[i]
		subText := ""
		var res *http.Response
		var e error
		switch mode {
		case "socks5":
			res, e = tools.GetBySocks5Proxy(subscribe.URL, address, port, 10)
		case "http":
			res, e = tools.GetByHTTPProxy(subscribe.URL, address, port, 10)
		case "none":
			res, e = tools.GetNoProxy(subscribe.URL, 10)
		}
		if e != nil {
			log.Warn(e)
			break
		}
		subText = tools.ReadDate(res)
		log.Info("访问 [", subscribe.URL, "] -- ", res.Status)
		_ = res.Body.Close()
		links := protocols.Sub2links(subText)
		var num = 0
		for _, link := range links {
			data := protocols.ParseLink(link)
			if data != nil {
				n := new(node)
				n.data = data
				n.Link = link
				n.Subid = subscribe.ID
				newNodes = append(newNodes, n)
				num++
			}
		}
		if num != 0 {
			ids += subscribe.ID + ","
		}
		log.Info("从订阅 [", subscribe.URL, "] 获取了 [", num, "] 个节点")
	}
	log.Info("订阅更新完成，共更新 [", len(newNodes), "] 个节点")
	for i, x := range c.Nodes {
		if x.Subid == "" || !strings.Contains(ids, x.Subid) {
			if i == int(c.Index) {
				c.Index = uint(len(newNodes))
			}
			newNodes = append(newNodes, x)
		} else {
			if i == int(c.Index) {
				c.Index = 0
			}
		}
	}
	c.Nodes = newNodes
}

// 获取订阅信息
func (c *Core) GetSubs(key string) [][]string {
	l := len(c.Subs)
	indexList := format.IndexDeal(key, l)
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
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
func (c *Core) SetSubs(key, using, url, remarks string) {
	l := len(c.Subs)
	indexList := format.IndexDeal(key, l)
	if len(indexList) == 0 {
		return
	}
	defer c.Save()
	for _, x := range indexList {
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
func (c *Core) DelSubs(key string) {
	l := len(c.Subs)
	indexList := format.IndexDeal(key, l)
	if len(indexList) == 0 {
		return
	}
	defer c.Save()
	result := make([]*sub, 0, l-len(indexList))
	j := 0
	for i, y := range c.Subs {
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
	c.Subs = result
	log.Info("共删除了 ", l-len(c.Subs), " 条订阅")
}
