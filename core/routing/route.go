// core/routing/route.go 负责路由规则的定义与相关操作
package routing

import (
	"Txray/core"
	"Txray/core/manage"
	"Txray/log"
	"strconv"
)

// 添加规则
// AddRule 函数用于添加路由规则
// rt: 规则类型，支持 TypeBlock、TypeDirect 和 TypeProxy
// list: 规则列表，支持多个规则以切片形式传入
// 返回值: 成功添加的规则数量
func AddRule(rt Type, list ...string) int {
	defer route.save()
	count := 0
	for _, rule := range list {
		if rule != "" {
			r := &routing{
				Data: rule,
				Mode: GetRuleMode(rule),
			}
			count += 1
			switch rt {
			case TypeBlock:
				route.Block = append(route.Block, r)
			case TypeDirect:
				route.Direct = append(route.Direct, r)
			case TypeProxy:
				route.Proxy = append(route.Proxy, r)
			}
		}
	}
	return count
}

// 根据规则类型和关键字获取匹配的规则
// GetRule 函数用于获取指定类型和关键字的路由规则
// rt: 规则类型，支持 TypeBlock、TypeDirect 和 TypeProxy
// key: 关键字，用于匹配规则
// 返回值: 匹配的规则列表，每条规则以切片形式返回，包含规则序号、模式和数据
func GetRule(rt Type, key string) [][]string {
	var rules []*routing
	switch rt {
	case TypeDirect:
		rules = route.Direct
	case TypeProxy:
		rules = route.Proxy
	case TypeBlock:
		rules = route.Block
	}
	indexList := core.IndexList(key, len(rules))
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		r := rules[x-1]
		result = append(result, []string{
			strconv.Itoa(x),
			string(r.Mode),
			r.Data,
		})
	}
	return result
}

// 对路由数据进行分组
// GetRulesGroupData 函数用于获取指定类型路由规则的 IP 和域名分组数据
// rt: 规则类型，支持 TypeBlock、TypeDirect 和 TypeProxy
// 返回值: 两个切片，分别包含 IP 列表和域名列表
func GetRulesGroupData(rt Type) ([]string, []string) {
	ips := make([]string, 0)
	domains := make([]string, 0)
	var rules []*routing
	switch rt {
	case TypeDirect:
		rules = route.Direct
	case TypeProxy:
		rules = route.Proxy
	case TypeBlock:
		rules = route.Block
	}
	for _, x := range rules {
		if x.Mode == "Domain" {
			domains = append(domains, x.Data)
		} else {
			ips = append(ips, x.Data)
		}
	}
	return ips, domains
}

// 删除规则
// DelRule 函数用于删除指定类型和关键字的路由规则
// rt: 规则类型，支持 TypeBlock、TypeDirect 和 TypeProxy
// key: 关键字，用于匹配规则
func DelRule(rt Type, key string) {
	var rules []*routing
	switch rt {
	case TypeDirect:
		rules = route.Direct
	case TypeProxy:
		rules = route.Proxy
	case TypeBlock:
		rules = route.Block
	}
	indexList := core.IndexList(key, len(rules))
	if len(indexList) == 0 {
		return
	}
	defer route.save()
	result := make([]*routing, 0)
	for i, rule := range rules {
		if !manage.HasIn(i+1, indexList) {
			result = append(result, rule)
		}
	}
	switch rt {
	case TypeDirect:
		route.Direct = result
	case TypeProxy:
		route.Proxy = result
	case TypeBlock:
		route.Block = result
	}
	log.Info("删除了 [", len(indexList), "] 条规则")
}

// RuleLen 函数用于获取指定类型的规则数量
// rt: 规则类型，支持 TypeBlock、TypeDirect 和 TypeProxy
// 返回值: 规则数量
func RuleLen(rt Type) int {
	switch rt {
	case TypeDirect:
		return len(route.Direct)
	case TypeProxy:
		return len(route.Proxy)
	case TypeBlock:
		return len(route.Block)
	default:
		return 0
	}
}
