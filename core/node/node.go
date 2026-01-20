// core/node/node.go 负责节点对象的定义与相关操作
package node

import (
	"Txray/core/protocols"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Node 节点结构体，包含协议、订阅ID、数据和测试结果
type Node struct {
	protocols.Protocol `json:"-"`
	SubID              string  `json:"sub_id"`
	Data               string  `json:"data"`
	TestResult         float64 `json:"-"`
}

// TestResultStr 返回格式化的测试结果字符串
func (n *Node) TestResultStr() string {
	if n.TestResult == 0 {
		return ""
	} else if n.TestResult >= 99998 {
		return "-1ms"
	} else {
		return fmt.Sprintf("%.4vms", n.TestResult)
	}
}

// NewNode 创建一个新的节点
func NewNode(link, id string) *Node {
	if data := protocols.ParseLink(link); data != nil {
		return &Node{Protocol: data, SubID: id}
	}
	return nil
}

// NewNodeByData 根据协议数据创建节点
func NewNodeByData(protocol protocols.Protocol) *Node {
	return &Node{Protocol: protocol}
}

// ParseData 反序列化Data字段，解析出协议
func (n *Node) ParseData() {
	n.Protocol = protocols.ParseLink(n.Data)
}

// Serialize2Data 序列化节点数据为Data字段
func (n *Node) Serialize2Data() {
	n.Data = n.GetLink()
}

var WG sync.WaitGroup

// Tcping 测试节点的TCP连通性，并计算平均响应时间
func (n *Node) Tcping() {
	count := 3
	var sum float64
	var timeout time.Duration = 3 * time.Second
	isTimeout := false
	for range count {
		start := time.Now()
		d := net.Dialer{Timeout: timeout}
		conn, err := d.Dial("tcp", net.JoinHostPort(n.GetAddr(), fmt.Sprintf("%d", n.GetPort())))
		if err != nil {
			isTimeout = true
			break
		}
		conn.Close()
		elapsed := time.Since(start)
		sum += float64(elapsed.Nanoseconds()) / 1e6
	}
	if isTimeout {
		n.TestResult = 99999
	} else {
		n.TestResult = sum / float64(count)
	}
	WG.Done()
}

// Show 显示节点信息
func (n *Node) Show() {
	ShowTopBottomSepLine('=', strings.Split(n.GetInfo(), "\n")...)
}
