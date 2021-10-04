package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/log"
)

func (m *Manage) MoveToRecycle(n *node.Node) {
	if n != nil {
		m.RecycleNodeList = append(m.RecycleNodeList, n)
	}
}

func (m *Manage) MoveFormRecycle(key string) {
	indexList := core.IndexList(key, m.RecycleLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newNodeList := make([]*node.Node, 0)
	m.RecycleForEach(func(i int, n *node.Node) {
		if HasIn(i, indexList) {
			m.addNode(n)
			log.Info("回收站 ==> ", n.GetName())
		} else {
			newNodeList = append(newNodeList, n)
		}
	})
	m.RecycleNodeList = newNodeList
}

func (m *Manage) RecycleLen() int {
	return len(m.RecycleNodeList)
}

func (m *Manage) getRecycleNode(i int) *node.Node {
	if i >= 0 && i < m.RecycleLen() {
		return m.RecycleNodeList[i]
	}
	return nil
}

func (m *Manage) GetRecycleNode(index int) *node.Node {
	return m.getRecycleNode(index - 1)
}

func (m *Manage) RecycleForEach(funC func(int, *node.Node)) {
	for i, n := range m.RecycleNodeList {
		funC(i+1, n)
	}
}

func (m *Manage) ClearRecycle() {
	m.RecycleNodeList = make([]*node.Node, 0)
}
