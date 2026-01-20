// core/manage/filter_manage.go 负责节点过滤相关的管理操作
package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/log"
)

// AddFilter 为管理器添加一个过滤规则
// key: 过滤规则的关键字
func (m *Manage) AddFilter(key string) {
	m.Filter = append(m.Filter, node.NewNodeFilter(key))
	m.Save()
}

// RunFilter 执行过滤操作
// key: 过滤规则的关键字，支持部分匹配
func (m *Manage) RunFilter(key string) {
	defer m.Save()
	selectedNode := m.SelectedNode()
	newNodeList := make([]*node.Node, 0)
	// 如果没有指定过滤规则，则应用所有已启用的过滤规则
	if key == "" {
		m.NodeForEach(func(i int, n *node.Node) {
			// 检查节点是否可以被过滤
			if f := m.IsCanFilter(n); f != nil {
				log.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
			} else {
				newNodeList = append(newNodeList, n)
			}
		})

	// 如果指定了过滤规则，则仅应用该规则
	} else if f := node.NewNodeFilter(key); f != nil {
		m.NodeForEach(func(i int, n *node.Node) {
			// 检查节点是否匹配过滤规则
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

// FilterForEach 遍历所有过滤规则
// funC: 对每个过滤规则执行的函数
func (m *Manage) FilterForEach(funC func(int, *node.NodeFilter)) {
	for i, f := range m.Filter {
		funC(i+1, f)
	}
}

// getFilter 获取指定索引的过滤规则
// i: 过滤规则的索引
// 返回值: 对应的过滤规则指针，如果索引无效则返回nil
func (m *Manage) getFilter(i int) *node.NodeFilter {
	if i >= 0 && i < m.FilterLen() {
		return m.Filter[i]
	}
	return nil
}

// GetFilter 获取指定索引的过滤规则（对外暴露的接口）
// index: 过滤规则的索引（从1开始）
// 返回值: 对应的过滤规则指针，如果索引无效则返回nil
func (m *Manage) GetFilter(index int) *node.NodeFilter {
	return m.getFilter(index - 1)
}

// DelFilter 删除指定的过滤规则
// key: 过滤规则的关键字
func (m *Manage) DelFilter(key string) {
	indexList := core.IndexList(key, m.FilterLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newFilterList := make([]*node.NodeFilter, 0)
	// 遍历当前所有过滤规则，删除匹配的规则
	m.FilterForEach(func(i int, filter *node.NodeFilter) {
		if !HasIn(i, indexList) {
			newFilterList = append(newFilterList, filter)
		}
	})
	m.Filter = newFilterList
}

// SetFilter 设置过滤规则的启用状态
// key: 过滤规则的关键字
// isOpen: 是否启用过滤规则
func (m *Manage) SetFilter(key string, isOpen bool) {
	indexList := core.IndexList(key, m.FilterLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	// 设置对应的过滤规则为启用或禁用
	for _, index := range indexList {
		if filter := m.GetFilter(index); filter != nil {
			filter.IsUse = isOpen
		}
	}
}

// FilterLen 获取当前过滤规则的数量
// 返回值: 过滤规则的数量
func (m *Manage) FilterLen() int {
	return len(m.Filter)
}

// IsCanFilter 检查节点是否可以被过滤
// n: 待检查的节点
// 返回值: 如果节点可以被过滤，返回对应的过滤规则；否则返回nil
func (m *Manage) IsCanFilter(n *node.Node) *node.NodeFilter {
	if n == nil {
		return nil
	}
	var f *node.NodeFilter = nil
	// 遍历所有过滤规则，检查节点是否匹配
	m.FilterForEach(func(i int, filter *node.NodeFilter) {
		if filter.IsUse && filter.IsMatch(n) {
			f = filter
		}
	})
	return f
}
