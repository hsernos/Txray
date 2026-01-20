// cmd/help/help.go 负责帮助文档内容的定义与输出
package help

import (
	_ "embed"
)

//go:embed help.txt
var Help string // 帮助文档内容

//go:embed setting.txt
var Setting string // 设置相关文档内容

//go:embed node.txt
var Node string // 节点相关文档内容

//go:embed sub.txt
var Sub string // 订阅相关文档内容

//go:embed filter.txt
var Filter string // 过滤器相关文档内容

//go:embed recycle.txt
var Recycle string // 回收站相关文档内容

//go:embed routing.txt
var Routing string // 路由相关文档内容

//go:embed alias.txt
var Alias string // 别名相关文档内容
