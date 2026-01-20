// core/routing/struct.go 负责路由相关结构体的定义与实现
package routing

import (
	"Txray/core"
	"Txray/log"
	"encoding/json"
	"os"
)

// routing 结构体用于定义单条路由规则
type routing struct {
	Data string `json:"data"` // 路由规则的具体内容
	Mode Mode   `json:"mode"` // 路由模式
}

// Routing 结构体用于定义完整的路由配置
type Routing struct {
	Proxy  []*routing `json:"proxy"`  // 代理规则
	Direct []*routing `json:"direct"` // 直连规则
	Block  []*routing `json:"block"`  // 阻止规则
}

var route *Routing = NewRouting()

// NewRouting 函数用于创建一个新的 Routing 实例
func NewRouting() *Routing {
	return &Routing{
		Proxy:  make([]*routing, 0),
		Direct: make([]*routing, 0),
		Block:  make([]*routing, 0),
	}
}

// init 函数用于初始化路由配置，加载配置文件
func init() {
	if _, err := os.Stat(core.RoutingFile); os.IsNotExist(err) {
		route.save() // 如果配置文件不存在，则创建一个新的
	} else {
		file, _ := os.Open(core.RoutingFile)
		defer file.Close()
		err := json.NewDecoder(file).Decode(route) // 解析配置文件
		if err != nil {
			log.Error(err)
		}
	}
}

// save 方法用于将 Routing 实例保存到配置文件
func (r *Routing) save() {
	err := core.WriteJSON(r, core.RoutingFile) // 将结构体编码为 JSON 并写入文件
	if err != nil {
		log.Error(err)
	}
}
