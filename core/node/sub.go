package node

import (
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"Txray/tools/uuid"
	"net/http"
	"strings"
)

// 添加一条订阅记录
func AddSubscribe(url, remarks string) {
	defer data.save()
	s := subscribe{Remarks: remarks, URL: url, Using: true, ID: uuid.UUID32()}
	data.Subs = append(data.Subs, &s)
	log.Info("添加一条订阅: ", url)
}

// 获取订阅信息
func GetSubscribe(key string) [][]string {
	l := len(data.Subs)
	indexList := format.IndexDeal(key, l)
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		s := data.Subs[x]
		result = append(result, []string{
			tools.IntToStr(x + 1),
			s.Remarks,
			s.URL,
			tools.BoolToStr(s.Using),
		})
	}
	return result
}

// 修改订阅信息
func SetSubs(key, using, url, remarks string) {
	l := len(data.Subs)
	indexList := format.IndexDeal(key, l)
	if len(indexList) == 0 {
		return
	}
	defer data.save()
	for _, x := range indexList {
		s := data.Subs[x]
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
func DelSubs(key string) {
	l := len(data.Subs)
	indexList := format.IndexDeal(key, l)
	if len(indexList) == 0 {
		return
	}
	defer data.save()
	result := make([]*subscribe, 0, l-len(indexList))
	j := 0
	for i, y := range data.Subs {
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
	data.Subs = result
	log.Info("共删除了 ", l-len(data.Subs), " 条订阅")
}

// 更新订阅节点
// mode取值为socks5、http、none表示通过socks5、http代理和不使用代理
func AddNodeBySub(key string, mode string, address string, port uint) {
	defer data.save()
	indexList := make([]int, 0, len(data.Subs))
	// 根据是否启用订阅来更新节点
	if key == "" {
		for i, x := range data.Subs {
			if x.Using {
				indexList = append(indexList, i)
			}
		}
		// 根据key选择订阅来更新节点（忽视是否启用）
	} else {
		indexList = format.IndexDeal(key, len(data.Subs))
	}

	newNodes := make([]*node, 0, len(data.Nodes))
	ids := ""
	for _, i := range indexList {
		sub := data.Subs[i]
		subText := ""
		var res *http.Response
		var e error
		switch mode {
		case "socks5":
			res, e = tools.GetBySocks5Proxy(sub.URL, address, port, 10)
		case "http":
			res, e = tools.GetByHTTPProxy(sub.URL, address, port, 10)
		case "none":
			res, e = tools.GetNoProxy(sub.URL, 10)
		}
		if e != nil {
			log.Warn(e)
			break
		}
		subText = tools.ReadDate(res)
		log.Info("访问 [", sub.URL, "] -- ", res.Status)
		_ = res.Body.Close()
		links := protocols.Sub2links(subText)
		var num = 0
		for _, link := range links {
			data := protocols.ParseLink(link)
			if data != nil {
				n := new(node)
				n.data = data
				n.Link = link
				n.Subid = sub.ID
				newNodes = append(newNodes, n)
				num++
			}
		}
		if num != 0 {
			ids += sub.ID + ","
		}
		log.Info("从订阅 [", sub.URL, "] 获取了 [", num, "] 个节点")
	}
	log.Info("订阅更新完成，共更新 [", len(newNodes), "] 个节点")
	for i, x := range data.Nodes {
		if x.Subid == "" || !strings.Contains(ids, x.Subid) {
			if i == data.Index-1 {
				data.Index = len(newNodes) + 1
			}
			newNodes = append(newNodes, x)
		} else {
			if i == data.Index-1 {
				data.Index = 1
			}
		}
	}
	data.Nodes = newNodes
}
