package manage

import (
	"Txray/core"
	"Txray/core/node"
	"Txray/core/sub"
	"Txray/log"
	"encoding/json"
	"os"
)

type Manage struct {
	Subs            []*sub.Subscirbe   `json:"subs"`
	Index           int                `json:"index"`
	NodeList        []*node.Node       `json:"nodes"`
	Filter          []*node.NodeFilter `json:"filter"`
	RecycleNodeList []*node.Node       `json:"-"`
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

func NewManage() *Manage {
	return &Manage{
		NodeList: make([]*node.Node, 0),
		Subs:     make([]*sub.Subscirbe, 0),
		Index:    0,
		Filter:   make([]*node.NodeFilter, 0),
	}
}

// Save 保存数据
func (m *Manage) Save() {
	err := core.WriteJSON(m, core.DataFile)
	if err != nil {
		log.Error(err)
	}
}
