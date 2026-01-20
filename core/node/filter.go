// core/node/filter.go 负责节点过滤相关的定义与操作
package node

import (
	"regexp"
	"strconv"
	"strings"
)

// NodeFilter 定义了节点过滤器的结构体
type NodeFilter struct {
	Mode  FilterMode `json:"mode"`  // 过滤模式
	Re    string     `json:"re"`    // 正则表达式
	IsUse bool       `json:"is_use"`// 是否启用
}

// FilterMode 定义了过滤模式的类型
type FilterMode string

const (
	NameFilter     FilterMode = "name" // 名称过滤
	AddrFilter     FilterMode = "addr" // 地址过滤
	PortFilter     FilterMode = "port" // 端口过滤
	ProtocolFilter FilterMode = "proto"// 协议过滤
)

// NewNodeFilter 根据给定的键值创建一个新的节点过滤器
func NewNodeFilter(key string) *NodeFilter {
	if strings.HasPrefix(key, "name:") {
		return &NodeFilter{Mode: NameFilter, Re: key[5:], IsUse: true}
	} else if strings.HasPrefix(key, "addr:") {
		return &NodeFilter{Mode: AddrFilter, Re: key[5:], IsUse: true}
	} else if strings.HasPrefix(key, "port:") {
		return &NodeFilter{Mode: PortFilter, Re: key[5:], IsUse: true}
	} else if strings.HasPrefix(key, "proto:") {
		return &NodeFilter{Mode: ProtocolFilter, Re: key[6:], IsUse: true}
	} else {
		return &NodeFilter{Mode: NameFilter, Re: key, IsUse: true}
	}
}

// IsMatch 检查节点是否匹配过滤器
func (nf *NodeFilter) IsMatch(n *Node) bool {
	regexp, _ := regexp.Compile(nf.Re)
	if n != nil {
		return regexp.MatchString(nf.getData(n))
	}
	return false
}

// String 返回过滤器的字符串表示
func (nf *NodeFilter) String() string {
	switch nf.Mode {
	case AddrFilter:
		return "addr:" + nf.Re
	case PortFilter:
		return "port:" + nf.Re
	case ProtocolFilter:
		return "proto:" + nf.Re
	case NameFilter:
		return "name:" + nf.Re
	default:
		return "name:" + nf.Re
	}
}

// getData 获取节点的相关数据
func (nf *NodeFilter) getData(n *Node) string {
	if n == nil {
		return ""
	}
	switch nf.Mode {
	case AddrFilter:
		return n.GetAddr()
	case PortFilter:
		return strconv.Itoa(n.GetPort())
	case ProtocolFilter:
		return string(n.GetProtocolMode())
	case NameFilter:
		return n.GetName()
	default:
		return n.GetName()
	}
}
