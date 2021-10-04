package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/core/sub"
	"Txray/log"
	"strings"
)

func (m *Manage) SubForEach(funC func(int, *sub.Subscirbe)) {
	for i, n := range m.Subs {
		funC(i+1, n)
	}
}

func (m *Manage) AddSubscirbe(subscirbe *sub.Subscirbe) {
	if m.HasSub(subscirbe.ID()) {
		log.Warn("该订阅链接已存在")
	} else {
		m.Subs = append(m.Subs, subscirbe)
		m.Save()

	}
}

func (m *Manage) SubLen() int {
	return len(m.Subs)
}

func (m *Manage) getSub(i int) *sub.Subscirbe {
	if i >= 0 && i < m.SubLen() {
		return m.Subs[i]
	}
	return nil
}

func (m *Manage) GetSub(i int) *sub.Subscirbe {
	return m.getSub(i - 1)
}

func (m *Manage) UpdataNode(opt sub.UpdataOption) {
	if opt.Key == "" {
		m.SubForEach(func(i int, subscirbe *sub.Subscirbe) {
			if subscirbe.Using {
				m.updataNode(subscirbe, opt)
			}
		})
	} else {
		for _, index := range core.IndexList(opt.Key, m.SubLen()) {
			m.updataNode(m.getSub(index), opt)
		}
	}
}

func (m *Manage) updataNode(subscirbe *sub.Subscirbe, opt sub.UpdataOption) {
	links := subscirbe.UpdataNode(opt)
	if len(links) == 0 {
		return
	}
	count := 0
	m.DelNodeById(subscirbe.ID())
	for _, link := range links {
		if ok := m.AddNode(node.NewNode(link, subscirbe.ID())); ok {
			count += 1
		}
	}
	log.Infof("从订阅 [%s] 获取了 '%d' 个节点", subscirbe.Url, count)
}

func (m *Manage) HasSub(id string) bool {
	ok := false
	m.SubForEach(func(i int, subscirbe *sub.Subscirbe) {
		if subscirbe.ID() == id {
			ok = true
		}
	})
	return ok
}

func (m *Manage) DelSub(key string) {
	indexList := core.IndexList(key, m.SubLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newSubList := make([]*sub.Subscirbe, 0)
	m.SubForEach(func(i int, subscirbe *sub.Subscirbe) {
		if !HasIn(i, indexList) {
			newSubList = append(newSubList, subscirbe)
		}
	})
	m.Subs = newSubList
}

func (m *Manage) SetSub(key string, using, url, name string) {
	indexList := core.IndexList(key, m.SubLen())
	if len(indexList) == 0 {
		return
	}
	if len(indexList) != 1 && url != "" {
		log.Warn("订阅链接不可以批量更改")
		return
	}
	defer m.Save()
	for _, index := range indexList {
		subscribe := m.GetSub(index)
		switch strings.ToLower(using) {
		case "true", "yes", "y":
			subscribe.Using = true
		case "false", "no", "n":
			subscribe.Using = false
		}
		if url != "" {
			subscribe.Url = url
		}
		if name != "" {
			subscribe.Name = name
		}
	}
}
