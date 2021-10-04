package manage

import "Txray/core/node"

func (m *Manage) SelectedNode() *node.Node {
	return m.GetNode(m.SelectedIndex())
}

func (m *Manage) SelectedIndex() int {
	return m.Index
}

func (m *Manage) SetSelectedIndex(index int) int {
	m.Index = index
	return m.Index
}

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
