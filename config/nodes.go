package config

import (
	log "Tv2ray/logger"
	"Tv2ray/tools"
	"Tv2ray/vmess"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 获取tcp延迟最小的num个节点
func (c *Config) GetNodeTestSort(num int) []int {
	var indexs []int
	list := make(tools.Sorts, 0)
	for i, node := range c.Nodes {
		if node.TestResult != "" && node.TestResult != "-1ms" {
			f, _ := strconv.ParseFloat(node.TestResult[:len(node.TestResult)-2], 32)
			list = append(list, tools.Sort{i, float32(f)})
		}
	}
	sort.Sort(sort.Reverse(list))
	for _, x := range list {
		indexs = append(indexs, x.Index)
	}
	if num > 0 && num <= len(indexs) {
		return indexs[len(indexs)-num:]
	}
	return indexs
}

// 获取node数据
func (c *Config) GetNodes(key string) [][]string {
	l := len(c.Nodes)
	var indexs []int
	if key == "test" || key == "t" || key == "tcping" {
		indexs = c.GetNodeTestSort(-1)
	} else {
		indexs = tools.IndexDeal(key, l)
	}
	result := make([][]string, 0, len(indexs))
	for _, x := range indexs {
		node := c.Nodes[x]
		result = append(result, []string{
			tools.IntToStr(x),
			node.Remarks,
			node.Address,
			tools.UintToStr(node.Port),
			node.TestResult,
		})
	}
	return result
}

func (c *Config) GetNode(index uint) *node {
	if int(index) >= len(c.Nodes) {
		return nil
	} else {
		return c.Nodes[index]
	}
}

// 查找node数据
func (c *Config) FindNodes(key string) [][]string {
	result := make([][]string, 0)
	for i, node := range c.Nodes {
		if strings.Index(node.Remarks, key) >= 0 {
			result = append(result, []string{
				tools.IntToStr(i),
				node.Remarks,
				node.Address,
				tools.UintToStr(node.Port),
				node.TestResult,
			})
		}
	}
	return result
}

// 获取选定节点索引
func (c *Config) GetNodeIndex() uint {
	return c.Index
}

// 导出node数据
func (c *Config) ExportNodes(key string) []string {
	l := len(c.Nodes)
	indexs := tools.IndexDeal(key, l)
	result := make([]string, 0, len(indexs))
	for _, x := range indexs {
		node := c.Nodes[x]
		result = append(result, nodeToVmessobj(node).ToLink())
	}
	return result
}

// 获取节点代理访问外网的延迟
func (c *Config) TestNode(url string) (string, string) {
	start := time.Now()
	res, e := tools.GetBySocks5Proxy(url, "127.0.0.1", c.Settings.Port, 10)
	elapsed := time.Since(start)
	if e != nil {
		log.Warn(e)
		return "-1ms", "Error"
	}
	return fmt.Sprintf("%4.0fms", float32(elapsed.Nanoseconds())/1e6), res.Status
}

// 测试节点的tcp延迟（多线程）
func (c *Config) PingNodes(key string) {
	l := len(c.Nodes)
	indexs := tools.IndexDeal(key, l)
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
		go tools.Go_Tcping(node.Address, int(node.Port), 4, chs[i])
	}
	var min float32 = 30000
	index := -1
	for i, ch := range chs {
		node := c.Nodes[indexs[i]]
		d := <-ch
		node.TestResult = fmt.Sprintf("%.4vms", d)
		if d < min && d != -1 {
			min = d
			index = indexs[i]
		}
	}
	if index != -1 {
		c.Index = uint(index)
	}
}

// 根据vmess链接文件批量添加节点
func (c *Config) AddNodeByFile(path string) {
	if tools.IsFile(path) {
		log.Info("读取文件中...")
		links := tools.ReadFile(path)
		log.Info("读取文件完成，解析vmess链接如下")
		for index, x := range links {
			log.Info(fmt.Sprintf("%3d ", index), x)
		}
		c.AddNodeByVmessLinks(links)
	} else {
		log.Warn("该文件不存在")
	}
}

// 根据订阅文件批量添加节点
func (c *Config) AddNodeBySubFile(path string) {
	if tools.IsFile(path) {
		subText := tools.ReadFile(path)[0]
		log.Info("解析文件中...")
		links := vmess.Sub2links(subText)
		log.Info("解析文件完成，解析vmess链接如下: ")
		log.Info("=======================================================================")
		for _, x := range links {
			fmt.Println(x)
		}
		log.Info("=======================================================================")
		c.AddNodeByVmessLinks(links)
	} else {
		log.Warn("该文件不存在")
	}
}

// 根据vmess链接添加节点
func (c *Config) AddNodeByVmessLinks(links []string) {
	defer c.SaveJSON()
	objs := vmess.Links2vmessObjs(links)
	for _, obj := range objs {
		c.Nodes = append(c.Nodes, vmessObjToNode(obj, ""))
	}
	log.Info("更新了 [", len(objs), "] 个节点")
}

// 添加节点
func (c *Config) AddNode(remarks, address, port, id, aid, security, network, types, host, path, tls string) {
	if !tools.IsUint(port) {
		log.Warn("端口必须为正整数")
		return
	}
	if !tools.IsInt(aid) {
		log.Warn("aid必须为整数")
		return
	}
	defer c.SaveJSON()
	n := node{}
	n.Remarks = remarks
	n.Address = address
	n.Port = tools.StrToUint(port)
	n.ID = id
	n.AlterID = tools.StrToInt(aid)
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

// 删除节点信息
func (c *Config) DelNodes(key string) {
	l := len(c.Nodes)
	indexs := tools.IndexDeal(key, l)
	if len(indexs) == 0 {
		return
	}
	defer c.SaveJSON()
	result := make([]*node, 0, l-len(indexs))
	j := 0
	for i, y := range c.Nodes {
		if j < len(indexs) {
			if i == indexs[j] {
				if indexs[j] == int(c.Index) {
					c.Index = 0
				} else if indexs[j] < int(c.Index) {
					c.Index -= 1
				}
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
