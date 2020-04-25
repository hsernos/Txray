package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	log "v3ray/logger"
	. "v3ray/tool"
)

// GetNodes 获取node数据
func (c *Config) GetNodes(key string) [][]string {
	l := len(c.Nodes)
	var indexs []int
	if key == "tcping" {
		list := make(Sorts, 0)
		for i, node := range c.Nodes {
			if node.TestResult != "" && node.TestResult != "0ms" {
				f, _ := strconv.ParseFloat(node.TestResult[:len(node.TestResult)-2], 32)
				list = append(list, Sort{i, float32(f)})
			}
		}
		sort.Sort(sort.Reverse(list))
		for _, x := range list {
			indexs = append(indexs, x.Index)
		}
	} else {
		indexs = IndexDeal(key, l)
	}
	result := make([][]string, 0, len(indexs))
	for _, x := range indexs {
		node := c.Nodes[x]
		result = append(result, []string{
			IntToStr(x),
			node.Remarks,
			node.Address,
			UintToStr(node.Port),
			node.Security,
			node.TestResult,
		})
	}
	return result
}

// FindNodes 查找node数据
func (c *Config) FindNodes(key string) [][]string {
	result := make([][]string, 0)
	for i, node := range c.Nodes {
		if strings.Index(node.Remarks, key) >= 0 {
			result = append(result, []string{
				IntToStr(i),
				node.Remarks,
				node.Address,
				UintToStr(node.Port),
				node.Security,
				node.TestResult,
			})
		}
	}
	return result
}

// GetNodeIndex 获取选定节点索引
func (c *Config) GetNodeIndex() uint {
	return c.Index
}

// ExportNodes 导出node数据
func (c *Config) ExportNodes(key string) []string {
	l := len(c.Nodes)
	indexs := IndexDeal(key, l)
	result := make([]string, 0, len(indexs))
	for _, x := range indexs {
		node := c.Nodes[x]
		result = append(result, nodeToVmessobj(node).GetVmesslink())
	}
	return result
}

// PingNodes ping node
func (c *Config) PingNodes(key string) {
	l := len(c.Nodes)
	indexs := IndexDeal(key, l)
	if indexs == nil || len(indexs) == 0 {
		return
	}
	for _, node := range c.Nodes {
		node.TestResult = ""
	}
	chs := make([]chan float32, len(indexs))
	for i, x := range indexs {
		node := c.Nodes[x]
		chs[i] = make(chan float32)
		go Go_Tcping(node.Address, int(node.Port), 4, chs[i])

	}
	var min float32 = 30000
	var index int = -1
	for i, ch := range chs {
		node := c.Nodes[indexs[i]]
		d := <-ch
		node.TestResult = fmt.Sprintf("%.4vms", d)
		if d < min && d != 0 {
			min = d
			index = indexs[i]
		}
	}
	if index != -1 {
		c.Index = uint(index)
	}
}

// AddNodeByFile 根据vmess链接文件批量添加节点
func (c *Config) AddNodeByFile(path string) {
	if PathExists(path) {
		links := ReadFile(path)
		c.AddNodeByVmessLinks(links)
	} else {
		log.Warn("该文件不存在")
	}
}

// AddNodeByVmessLinks 根据vmess链接添加节点
func (c *Config) AddNodeByVmessLinks(links []string) {
	defer c.SaveJSON()
	objs := VmessListToObj(links)
	for _, obj := range objs {
		c.Nodes = append(c.Nodes, vmessObjToNode(obj, ""))
	}
	log.Info("更新了 [", len(objs), "] 个节点")
}

// AddNode 添加节点
func (c *Config) AddNode(remarks, address, port, id, aid, security, network, types, host, path, tls string) {
	if !IsUint(port) {
		log.Warn("端口必须为正整数")
		return
	}
	if !IsInt(aid) {
		log.Warn("aid必须为整数")
		return
	}
	defer c.SaveJSON()
	n := node{}
	n.Remarks = remarks
	n.Address = address
	n.Port = StrToUint(port)
	n.ID = id
	n.AlterID = StrToInt(aid)
	n.Security = security
	n.Network = network
	n.HeaderType = types
	n.RequestHost = host
	n.Path = path
	n.StreamSecurity = tls
	n.ConfigVersion = "2"
	c.Nodes = append(c.Nodes, &n)
	log.Info("添加节点成功")
}

// DelNodes 删除节点信息
func (c *Config) DelNodes(key string) {
	l := len(c.Nodes)
	indexs := IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]*node, 0, l-len(indexs))
	j := 0
	for i, y := range c.Nodes {
		if j < len(indexs) {
			if i == indexs[j] {
				j++
			} else {
				result = append(result, y)
			}
		} else {
			result = append(result, y)
		}
	}
	c.Nodes = result
	log.Info("删除了 [", l-len(c.Nodes), "] 条")
}
