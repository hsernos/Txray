// core/manage/recycle_manage.go 负责回收站相关的节点管理操作
package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/log"
)

// MoveToRecycle 将节点移动到回收站
func (m *Manage) MoveToRecycle(n *node.Node) {
	if n != nil {
		m.RecycleNodeList = append(m.RecycleNodeList, n)
	}
}

// MoveFormRecycle 从回收站恢复节点
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

// RecycleLen 获取回收站节点数量
func (m *Manage) RecycleLen() int {
	return len(m.RecycleNodeList)
}

// getRecycleNode 获取指定索引的回收站节点
func (m *Manage) getRecycleNode(i int) *node.Node {
	if i >= 0 && i < m.RecycleLen() {
		return m.RecycleNodeList[i]
	}
	return nil
}

// GetRecycleNode 获取指定索引的回收站节点（对外暴露接口）
func (m *Manage) GetRecycleNode(index int) *node.Node {
	return m.getRecycleNode(index - 1)
}

// RecycleForEach 遍历回收站中的每个节点
func (m *Manage) RecycleForEach(funC func(int, *node.Node)) {
	for i, n := range m.RecycleNodeList {
		funC(i+1, n)
	}
}

// ClearRecycle 清空回收站
func (m *Manage) ClearRecycle() {
	m.RecycleNodeList = make([]*node.Node, 0)
}
