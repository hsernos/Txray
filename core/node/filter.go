package node

import (
	"regexp"
	"strconv"
	"strings"
)

type NodeFilter struct {
	Mode  FilterMode `json:"mode"`
	Re    string     `json:"re"`
	IsUse bool       `json:"is_use"`
}

type FilterMode string

const (
	NameFilter     FilterMode = "name"
	AddrFilter     FilterMode = "addr"
	PortFilter     FilterMode = "port"
	ProtocolFilter FilterMode = "proto"
)

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

func (nf *NodeFilter) IsMatch(n *Node) bool {
	regexp, _ := regexp.Compile(nf.Re)
	if n != nil {
		return regexp.MatchString(nf.getData(n))
	}
	return false
}

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
