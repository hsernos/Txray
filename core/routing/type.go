// core/routing/type.go 负责路由相关类型、常量的定义与实现
package routing

// Type 表示路由的类型
type Type string

const (
	// TypeProxy 表示代理类型的路由
	TypeProxy Type = "Proxy"
	// TypeDirect 表示直连类型的路由
	TypeDirect Type = "Direct"
	// TypeBlock 表示阻止类型的路由
	TypeBlock Type = "Block"
)
