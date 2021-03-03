package node

import (
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 获取选定节点索引编号
func GetSelectedIndex() int {
	return data.Index
}

// 存在index对应的节点（从1开始）
func HasNode(index int) bool {
	return index > 0 && index <= len(data.Nodes)
}

func NodesSize() int {
	return len(data.Nodes)
}

func SetIndex(index int) {
	data.Index = index
	defer data.save()
}

func GetNode(index int) protocols.Protocol {
	if HasNode(index) {
		return data.Nodes[index-1].data
	}
	log.Warn("所选节点索引编号超出范围")
	return nil
}

// 打印节点信息
func Show(index int) {
	if HasNode(index) {
		info := data.Nodes[index-1].data.GetInfo()
		format.ShowTopBottomSepLine('=', strings.Split(info, "\n")...)
	} else {
		log.Warn("所选节点索引编号超出范围")
	}
}

// 获取Tcp排序索引
func GetTcpSortIndex(isDes bool) []int {
	list := make(tools.Sorts, 0, len(data.Nodes))
	for i, n := range data.Nodes {
		test := n.TestResult
		if test != "" && test != "-1ms" {
			f, _ := strconv.ParseFloat(strings.TrimRight(test, "ms"), 64)
			list = append(list, tools.Sort{Index: i, Value: f})
		}
	}
	if isDes {
		sort.Sort(sort.Reverse(list))
	} else {
		sort.Sort(list)
	}
	var indexList []int
	for _, x := range list {
		indexList = append(indexList, x.Index)
	}
	return indexList
}

// 根据key选择节点
func GetNodes(key string) [][]string {
	length := len(data.Nodes)
	var indexList []int
	isTcp := false
	if key == "test" || key == "t" || key == "tcp" {
		indexList = GetTcpSortIndex(true)
		isTcp = true
	} else {
		indexList = format.IndexDeal(key, length)
	}
	result := make([][]string, 0, len(indexList))
	for i, x := range indexList {
		n := data.Nodes[x]
		if isTcp {
			result = append(result, []string{
				tools.IntToStr(x + 1),
				tools.IntToStr(len(indexList) - i),
				string(n.data.GetProtocolMode()),
				n.data.GetName(),
				n.data.GetAddr(),
				strconv.Itoa(n.data.GetPort()),
				n.TestResult,
			})
		} else {
			result = append(result, []string{
				tools.IntToStr(x + 1),
				string(n.data.GetProtocolMode()),
				n.data.GetName(),
				n.data.GetAddr(),
				strconv.Itoa(n.data.GetPort()),
				n.TestResult,
			})
		}
	}
	return result
}

// 查找节点
func FindNodes(key string) [][]string {
	result := make([][]string, 0)
	for i, n := range data.Nodes {
		if strings.Index(n.data.GetName(), key) >= 0 {
			result = append(result, []string{
				tools.IntToStr(i + 1),
				string(n.data.GetProtocolMode()),
				n.data.GetName(),
				n.data.GetAddr(),
				strconv.Itoa(n.data.GetPort()),
				n.TestResult,
			})
		}
	}
	return result
}

// 导出node数据
func ExportNodes(key string) []string {
	indexList := format.IndexDeal(key, len(data.Nodes))
	result := make([]string, 0, len(indexList))
	for _, x := range indexList {
		n := data.Nodes[x]
		result = append(result, n.Link)
	}
	return result
}

// 测试节点的tcp延迟（多线程）
func PingNodes(key string) {
	indexList := format.IndexDeal(key, len(data.Nodes))
	if indexList == nil || len(indexList) == 0 {
		return
	}
	defer data.save()
	for _, n := range data.Nodes {
		n.TestResult = ""
	}
	chs := make([]chan float32, len(indexList))
	for i, x := range indexList {
		n := data.Nodes[x]
		chs[i] = make(chan float32)
		go tools.Go_Tcping(n.data.GetAddr(), n.data.GetPort(), 4, chs[i])
	}
	var min float32 = 30000
	index := -1
	for i, ch := range chs {
		n := data.Nodes[indexList[i]]
		d := <-ch
		n.TestResult = fmt.Sprintf("%.4vms", d)
		if d < min && d != -1 {
			min = d
			index = indexList[i]
		}
	}
	if index != -1 {
		data.Index = index + 1
	}
}

// 删除节点信息
func DelNodes(key string) {
	length := len(data.Nodes)
	indexList := format.OtherIndex(key, length)
	if len(indexList) == length {
		return
	}
	defer data.save()
	isDelete := true // 设置当前选择节点是否被删除标记
	newNodes := make([]*node, 0, len(indexList))
	for i, index := range indexList {
		if index == data.Index-1 {
			data.Index = i + 1
			isDelete = false
		}
		newNodes = append(newNodes, data.Nodes[index])
	}
	data.Nodes = newNodes
	if isDelete {
		data.Index = 1
	}
	log.Info("删除了 [", length-len(data.Nodes), "] 条")
}
