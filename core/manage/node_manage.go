// core/manage/node_manage.go 负责节点相关的增删查改等管理操作
package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/log"
)

// NodeLen 节点数量
func (m *Manage) NodeLen() int {
	return len(m.NodeList)
}

// GetNodeByIndex 获取节点
func (m *Manage) GetNode(index int) *node.Node {
	return m.getNode(index - 1)
}

// getNode 获取节点
func (m *Manage) getNode(i int) *node.Node {
	if i < m.NodeLen() && i >= 0 {
		return m.NodeList[i]
	}
	return nil
}

// addNode 添加节点到管理列表
// 返回值：bool，表示添加是否成功
func (m *Manage) addNode(n *node.Node) bool {
	if n == nil {
		return false
	}
	// 检查节点是否可以被过滤
	if f := m.IsCanFilter(n); f != nil {
		log.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
		return false
	}
	// 序列化节点数据
	n.Serialize2Data()
	// 将节点添加到节点列表
	m.NodeList = append(m.NodeList, n)
	return true
}

// AddNode 添加节点，并保存管理器状态
func (m *Manage) AddNode(n *node.Node) bool {
	ok := false
	// 尝试添加节点
	if ok = m.addNode(n); ok {
		m.Save() // 保存管理器状态
	}
	return ok
}

// NodeForEach 遍历节点列表，执行给定的函数
func (m *Manage) NodeForEach(funC func(int, *node.Node)) {
	for i, n := range m.NodeList {
		funC(i+1, n)
	}
}

// Tcping 并行对节点列表中的节点进行TCP Ping测试
// 测试完成后，节点会根据测试结果进行排序
func (m *Manage) Tcping() {
	m.NodeForEach(func(i int, n *node.Node) {
		node.WG.Add(1)
		go n.Tcping()
	})
	node.WG.Wait() // 等待所有TCP Ping测试完成
	defer m.Save()
	// 根据测试结果对节点进行排序
	m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
		return n1.TestResult < n2.TestResult
	})
	// 设置排序后的第一个节点为选中节点
	m.SetSelectedIndex(1)
}

// NodeSort 对节点列表进行排序
// less：排序规则，返回true表示n1排在n2前面
func (m *Manage) NodeSort(less func(*node.Node, *node.Node) bool) {
	if m.NodeLen() <= 1 {
		return
	}
	// 插入排序算法
	for i := 1; i < m.NodeLen(); i++ {
		preIndex := i - 1
		current := m.getNode(i)
		// 根据排序规则将节点插入到正确位置
		for preIndex >= 0 && !less(m.getNode(preIndex), current) {
			m.NodeList[preIndex+1] = m.NodeList[preIndex]
			preIndex -= 1
		}
		m.NodeList[preIndex+1] = current
	}
}

// Sort 根据给定的模式对节点列表进行排序
// mode：排序模式，具体含义如下：
// 0 - 反转节点顺序
// 1 - 按照协议模式排序
// 2 - 按照节点名称排序
// 3 - 按照节点地址排序
// 4 - 按照节点端口排序
// 5 - 按照测试结果排序
func (m *Manage) Sort(mode int) {
	selectedNode := m.SelectedNode()
	switch mode {
	case 0:
		// 反转节点顺序
		for i := 0; i < m.NodeLen()/2; i++ {
			j := m.NodeLen() - i - 1
			m.NodeList[i], m.NodeList[j] = m.NodeList[j], m.NodeList[i]
		}
	case 1:
		// 按照协议模式排序
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetProtocolMode() < n2.Protocol.GetProtocolMode()
		})
	case 2:
		// 按照节点名称排序
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetName() < n2.GetName()
		})
	case 3:
		// 按照节点地址排序
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetAddr() < n2.GetAddr()
		})
	case 4:
		// 按照节点端口排序
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetPort() < n2.GetPort()
		})
	case 5:
		// 按照测试结果排序
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.TestResult < n2.TestResult
		})
	default:
		return
	}
	defer m.Save()
	// 设置排序后的节点为选中节点
	m.SetSelectedIndexByNode(selectedNode)
}

// DelNode 根据给定的关键字删除节点
// key：关键字，可以是节点的名称、地址或其他属性
func (m *Manage) DelNode(key string) {
	indexList := core.IndexList(key, m.NodeLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	selectedNode := m.SelectedNode()
	newNodeList := make([]*node.Node, 0)
	// 遍历节点列表，删除匹配的节点
	m.NodeForEach(func(i int, n *node.Node) {
		if HasIn(i, indexList) {
			// 移动到回收站
			m.MoveToRecycle(n)
		} else {
			newNodeList = append(newNodeList, n)
		}
	})
	m.NodeList = newNodeList
	// 设置删除后仍然存在的节点为选中节点
	m.SetSelectedIndexByNode(selectedNode)
}

// DelNodeById 根据节点ID删除节点
// id：节点的唯一标识符
func (m *Manage) DelNodeById(id string) {
	defer m.Save()
	selectedNode := m.SelectedNode()
	newNodeList := make([]*node.Node, 0)
	// 遍历节点列表，删除匹配的节点
	m.NodeForEach(func(i int, n *node.Node) {
		if n.SubID == id {
			// 移动到回收站
			m.MoveToRecycle(n)
		} else {
			newNodeList = append(newNodeList, n)
		}
	})
	m.NodeList = newNodeList
	// 设置删除后仍然存在的节点为选中节点
	m.SetSelectedIndexByNode(selectedNode)
}

// GetNodeLink 根据给定的关键字获取节点链接
// key：关键字，可以是节点的名称、地址或其他属性
// 返回值：节点链接列表
func (m *Manage) GetNodeLink(key string) []string {
	links := make([]string, 0)
	// 获取匹配关键字的节点索引列表
	for _, index := range core.IndexList(key, m.NodeLen()) {
		links = append(links, m.GetNode(index).GetLink())
	}
	return links
}
