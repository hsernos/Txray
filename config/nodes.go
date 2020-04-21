package config

import (
	ping "github.com/vearne/go-ping"
	"time"
	log "v3ray/logger"
	"v3ray/tool"
)

// GetNodes 获取node数据
func (c *Config) GetNodes(key string) [][]string {
	l := len(c.Nodes)
	indexs := tool.IndexDeal(key, l)
	result := make([][]string, 0, len(indexs))
	for _, x := range indexs {
		node := c.Nodes[x]
		result = append(result, []string{
			tool.IntToStr(x),
			node.Remarks,
			node.Address,
			tool.UintToStr(node.Port),
			node.Security,
			node.Network,
			node.StreamSecurity,
			node.TestResult,
		})
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
	indexs := tool.IndexDeal(key, l)
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
	indexs := tool.IndexDeal(key, l)
	if indexs == nil || len(indexs) == 0 {
		return
	}
	for _, node := range c.Nodes{
		node.TestResult = ""
	}
	ipSlice := []string{}
	for _, x := range indexs {
		node := c.Nodes[x]
		ipSlice = append(ipSlice, node.Address)
	}
	bp, err := ping.NewBatchPinger(ipSlice, 3, time.Second*1, time.Second*6)
    if err != nil {
        log.Error(err)
	}
	bp.OnFinish = func(stSlice []*ping.Statistics) {
        for i, st := range stSlice{
			node := c.Nodes[indexs[i]]
			node.TestResult = st.AvgRtt.String()
        }

    }
	bp.Run()
}

// AddNodeByFile 根据vmess链接文件批量添加节点
func (c *Config) AddNodeByFile(path string) {
	if tool.PathExists(path) {
		links := tool.ReadFile(path)
		c.AddNodeByVmessLinks(links)
	} else {
		log.Warn("该文件不存在")
	}
}

// AddNodeByVmessLinks 根据vmess链接添加节点
func (c *Config) AddNodeByVmessLinks(links []string) {
	defer c.SaveJSON()
	objs := tool.VmessListToObj(links)
	for _, obj := range objs {
		c.Nodes = append(c.Nodes, vmessObjToNode(obj, ""))
	}
	log.Info("更新了 [", len(objs), "] 个节点")
}

// AddNode 添加节点
func (c *Config) AddNode(remarks, address, port, id, aid, security, network, types, host, path, tls string) {
	if !tool.IsUint(port) {
		log.Warn("端口必须为正整数")
		return
	}
	if !tool.IsInt(aid) {
		log.Warn("aid必须为整数")
		return
	}
	defer c.SaveJSON()
	n := node{}
	n.Remarks = remarks
	n.Address = address
	n.Port = tool.StrToUint(port)
	n.ID = id
	n.AlterID = tool.StrToInt(aid)
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
	indexs := tool.IndexDeal(key, l)
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