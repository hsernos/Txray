package core

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

// 初始化节点数据
func (n *node) initNodeData() {
	if n.data == nil {
		n.data = protocols.ParseLink(n.Link)
	}
}

// 初始化全部节点数据
func (c *Core) initNodesData() {
	for _, n := range c.Nodes {
		n.initNodeData()
	}
}

// ---------------------------------------------------------------------------------------------

// 打印节点信息
func (c *Core) Show(index int) {
	if c.HasNode(index) {
		info := c.Nodes[index-1].data.GetInfo()
		format.ShowTopBottomSepLine('=', strings.Split(info, "\n")...)
	} else {
		log.Warn("索引超出范围")

	}
}

// 获取选定节点索引
func (c *Core) GetNodeIndex() uint {
	return c.Index + 1
}

// 获取选定节点协议
func (c *Core) GetNodeMode() string {
	return c.Nodes[c.Index].data.GetProtocolMode()
}

// 获取选定节点链接
func (c *Core) GetNodeLink() string {
	return c.Nodes[c.Index].Link
}

// 存在index对应的节点（从1开始）
func (c *Core) HasNode(index int) bool {
	return index > 0 && index <= len(c.Nodes)
}

// 获取Tcp排序索引
func (c *Core) GetTcpSortIndex(isDes bool) []int {
	list := make(tools.Sorts, 0, len(c.Nodes))
	for i, n := range c.Nodes {
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
func (c *Core) GetNodes(key string) [][]string {
	length := len(c.Nodes)
	var indexList []int
	isTcp := false
	if key == "test" || key == "t" || key == "tcp" {
		indexList = c.GetTcpSortIndex(true)
		isTcp = true
	} else {
		indexList = format.IndexDeal(key, length)
	}
	result := make([][]string, 0, len(indexList))
	for i, x := range indexList {
		n := c.Nodes[x]
		if isTcp {
			result = append(result, []string{
				tools.IntToStr(x + 1),
				tools.IntToStr(len(indexList) - i),
				n.data.GetProtocolMode(),
				n.data.GetName(),
				n.data.GetAddr(),
				strconv.Itoa(n.data.GetPort()),
				n.TestResult,
			})
		} else {
			result = append(result, []string{
				tools.IntToStr(x + 1),
				n.data.GetProtocolMode(),
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
func (c *Core) FindNodes(key string) [][]string {
	result := make([][]string, 0)
	for i, n := range c.Nodes {
		if strings.Index(n.data.GetName(), key) >= 0 {
			result = append(result, []string{
				tools.IntToStr(i + 1),
				n.data.GetProtocolMode(),
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
func (c *Core) ExportNodes(key string) []string {
	indexList := format.IndexDeal(key, len(c.Nodes))
	result := make([]string, 0, len(indexList))
	for _, x := range indexList {
		n := c.Nodes[x]
		result = append(result, n.Link)
	}
	return result
}

// 测试节点的tcp延迟（多线程）
func (c *Core) PingNodes(key string) {
	indexList := format.IndexDeal(key, len(c.Nodes))
	if indexList == nil || len(indexList) == 0 {
		return
	}
	defer c.Save()
	for _, n := range c.Nodes {
		n.TestResult = ""
	}
	chs := make([]chan float32, len(indexList))
	for i, x := range indexList {
		n := c.Nodes[x]
		chs[i] = make(chan float32)
		go tools.Go_Tcping(n.data.GetAddr(), n.data.GetPort(), 4, chs[i])
	}
	var min float32 = 30000
	index := -1
	for i, ch := range chs {
		n := c.Nodes[indexList[i]]
		d := <-ch
		n.TestResult = fmt.Sprintf("%.4vms", d)
		if d < min && d != -1 {
			min = d
			index = indexList[i]
		}
	}
	if index != -1 {
		c.Index = uint(index)
	}
}

// ---------------------------------------------------------------------------------------------

// 从本地订阅文件批量添加节点
func (c *Core) AddNodeBySubFile(absPath string) {
	if tools.IsFile(absPath) {
		subText := tools.ReadFile(absPath)[0]
		log.Info("解析文件中...")
		links := protocols.Sub2links(subText)
		log.Info("解析文件完成，解析VMess链接如下: ")
		log.Info("=======================================================================")
		for index, x := range links {
			log.Info(fmt.Sprintf("%3d. ", index+1), x)
		}
		log.Info("=======================================================================")
		c.AddNodeByLinks(links...)
	} else {
		log.Warn("该文件不存在")
	}
}

// 从本地链接文件批量添加节点
func (c *Core) AddNodeByFile(absPath string) {
	if tools.IsFile(absPath) {
		log.Info("读取文件中...")
		links := tools.ReadFile(absPath)
		log.Info("读取文件完成，解析链接如下")
		for index, x := range links {
			log.Info(fmt.Sprintf("%3d. ", index+1), x)
		}
		c.AddNodeByLinks(links...)
	} else {
		log.Warn("该文件不存在")
	}
}

// 从订阅文本更新节点
func (c *Core) AddNodeBySubText(subtext string) {
	links := protocols.Sub2links(subtext)
	c.AddNodeByLinks(links...)
}

// 根据链接添加节点
func (c *Core) AddNodeByLinks(links ...string) {
	defer c.Save()
	var i = 0
	for _, link := range links {
		data := protocols.ParseLink(link)
		if data != nil {
			n := new(node)
			n.data = data
			n.Link = link
			c.Nodes = append(c.Nodes, n)
			i++
		} else {
			log.Warnf("添加失败: [%s]", link)
		}
	}
	log.Infof("成功更新 [%d] 个节点, 失败 [%d] 个", i, len(links)-i)
}

// 添加一个VMess节点
func (c *Core) AddVMessNode(remarks, addr, port, id, aid, network, types, host, path, tls string) {
	if !tools.IsNetPort(port) {
		log.Errorf("%s: 不是一个端口", port)
		return
	}
	if !tools.IsInt(aid) {
		log.Errorf("%s: 不是一个数字", aid)
		return
	}
	defer c.Save()
	// 初始化VMess节点数据
	vmessData := new(protocols.VMess)
	vmessData.V = "2"
	vmessData.Ps = remarks
	vmessData.Add = addr
	vmessData.Port = port
	vmessData.Id = id
	vmessData.Aid = aid
	vmessData.Net = network
	vmessData.Type = types
	vmessData.Host = host
	vmessData.Path = path
	vmessData.Tls = tls
	// 生成节点
	n := new(node)
	n.data = vmessData
	n.Link = vmessData.GetLink()
	// 添加到集合中
	c.Nodes = append(c.Nodes, n)
}

// 添加一个VLess节点
func (c *Core) AddVLESSNode() {
	//TODO
}

// 添加一个Trojan节点
func (c *Core) AddTrojanNode(remarks, addr, port, password string) {
	if !tools.IsNetPort(port) {
		log.Errorf("%s: 不是一个端口", port)
		return
	}
	defer c.Save()
	// 初始化VMess节点数据
	trojanData := new(protocols.Trojan)
	trojanData.Remarks = remarks
	trojanData.Address = addr
	trojanData.Port = port
	trojanData.Password = password
	// 生成节点
	n := new(node)
	n.data = trojanData
	n.Link = trojanData.GetLink()
	// 添加到集合中
	c.Nodes = append(c.Nodes, n)
}

// 添加一个ShadowSocks节点
func (c *Core) AddShadowSocksNode(remarks, addr, port, password, method string) {
	if !tools.IsNetPort(port) {
		log.Errorf("%s: 不是一个端口", port)
		return
	}
	defer c.Save()
	ss := new(protocols.ShadowSocks)
	ss.Remarks = remarks
	ss.Address = addr
	ss.Port = port
	ss.Password = password
	ss.Method = method
	// 生成节点
	n := new(node)
	n.data = ss
	n.Link = ss.GetLink()
	// 添加到集合中
	c.Nodes = append(c.Nodes, n)
}

// ---------------------------------------------------------------------------------------------

// 删除节点信息
func (c *Core) DelNodes(key string) {
	length := len(c.Nodes)
	indexList := format.OtherIndex(key, length)
	if len(indexList) == length {
		return
	}
	defer c.Save()
	isDelete := true // 设置当前选择节点是否被删除标记
	newNodes := make([]*node, 0, len(indexList))
	for i, index := range indexList {
		if index == int(c.Index) {
			c.Index = uint(i)
			isDelete = false
		}
		newNodes = append(newNodes, c.Nodes[index])
	}
	c.Nodes = newNodes
	if isDelete {
		c.Index = 0
	}
	log.Info("删除了 [", length-len(c.Nodes), "] 条")
}
