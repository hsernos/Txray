// core/manage/index_manage.go 负责节点索引相关的管理操作
package manage

import "Txray/core/node"

func (m *Manage) SelectedNode() *node.Node {
	return m.GetNode(m.SelectedIndex())
}

// SelectedIndex 返回当前选中的节点索引
func (m *Manage) SelectedIndex() int {
	return m.Index
}

// SetSelectedIndex 设置当前选中的节点索引
func (m *Manage) SetSelectedIndex(index int) int {
	m.Index = index
	return m.Index
}

// SetSelectedIndexByNode 根据节点指针设置当前选中的节点索引
func (m *Manage) SetSelectedIndexByNode(n *node.Node) int {
	m.SetSelectedIndex(1)
	if n != nil {
		m.NodeForEach(func(i int, n1 *node.Node) {
			if n == n1 {
				m.SetSelectedIndex(i)
				return
			}
		})
	}
	return m.Index
}
