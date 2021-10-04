package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/log"
)

func (m *Manage) AddFilter(key string) {
	m.Filter = append(m.Filter, node.NewNodeFilter(key))
	m.Save()
}

func (m *Manage) RunFilter(key string) {
	defer m.Save()
	selectedNode := m.SelectedNode()
	newNodeList := make([]*node.Node, 0)
	if key == "" {
		m.NodeForEach(func(i int, n *node.Node) {
			if f := m.IsCanFilter(n); f != nil {
				log.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
			} else {
				newNodeList = append(newNodeList, n)
			}
		})

	} else if f := node.NewNodeFilter(key); f != nil {
		m.NodeForEach(func(i int, n *node.Node) {
			if f.IsMatch(n) {
				log.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
			} else {
				newNodeList = append(newNodeList, n)
			}
		})
	}
	m.NodeList = newNodeList
	m.SetSelectedIndexByNode(selectedNode)
}

func (m *Manage) FilterForEach(funC func(int, *node.NodeFilter)) {
	for i, f := range m.Filter {
		funC(i+1, f)
	}
}

func (m *Manage) getFilter(i int) *node.NodeFilter {
	if i >= 0 && i < m.FilterLen() {
		return m.Filter[i]
	}
	return nil
}

func (m *Manage) GetFilter(index int) *node.NodeFilter {
	return m.getFilter(index - 1)
}

func (m *Manage) DelFilter(key string) {
	indexList := core.IndexList(key, m.FilterLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newFilterList := make([]*node.NodeFilter, 0)
	m.FilterForEach(func(i int, filter *node.NodeFilter) {
		if !HasIn(i, indexList) {
			newFilterList = append(newFilterList, filter)
		}
	})
	m.Filter = newFilterList
}

func (m *Manage) SetFilter(key string, isOpen bool) {
	indexList := core.IndexList(key, m.FilterLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	for _, index := range indexList {
		if filter := m.GetFilter(index); filter != nil {
			filter.IsUse = isOpen
		}
	}
}

func (m *Manage) FilterLen() int {
	return len(m.Filter)
}

func (m *Manage) IsCanFilter(n *node.Node) *node.NodeFilter {
	if n == nil {
		return nil
	}
	var f *node.NodeFilter = nil
	m.FilterForEach(func(i int, filter *node.NodeFilter) {
		if filter.IsUse && filter.IsMatch(n) {
			f = filter
		}
	})
	return f
}
