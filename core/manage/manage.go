// core/manage/manage.go 负责节点、订阅、回收站等核心数据的统一管理
package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/core/sub"
	"Txray/log"
	"encoding/json"
	"os"
)

// Manage 结构体用于管理核心数据，包括订阅、节点、过滤器等
type Manage struct {
	Subs            []*sub.Subscirbe   `json:"subs"`            // 订阅列表
	Index           int                `json:"index"`           // 当前索引
	NodeList        []*node.Node       `json:"nodes"`           // 节点列表
	Filter          []*node.NodeFilter `json:"filter"`          // 过滤器列表
	RecycleNodeList []*node.Node       `json:"-"`               // 回收站节点列表
}

var Manager *Manage

// 初始化
func init() {
	Manager = NewManage()
	if _, err := os.Stat(core.DataFile); os.IsNotExist(err) {
		Manager.Save()
	} else {
		file, _ := os.Open(core.DataFile)
		defer file.Close()
		err := json.NewDecoder(file).Decode(Manager)
		if err != nil {
			log.Error(err)
		}
		Manager.NodeForEach(func(i int, n *node.Node) {
			n.ParseData()
		})
	}
}

// NewManage 创建一个新的 Manage 实例
func NewManage() *Manage {
	return &Manage{
		NodeList: make([]*node.Node, 0),
		Subs:     make([]*sub.Subscirbe, 0),
		Index:    0,
		Filter:   make([]*node.NodeFilter, 0),
	}
}

// Save 保存数据到文件
func (m *Manage) Save() {
	err := core.WriteJSON(m, core.DataFile)
	if err != nil {
		log.Error(err)
	}
}
