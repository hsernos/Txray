package routing

import (
	"Txray/core"
	"Txray/core/manage"
	"Txray/log"
	"strconv"
)

// 添加规则
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
