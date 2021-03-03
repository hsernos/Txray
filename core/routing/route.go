package routing

import (
	"Txray/log"
	"Txray/tools"
	"Txray/tools/format"
)

// 添加规则
func AddRule(rt Type, list ...string) {
	defer route.save()
	num := 0
	for _, rule := range list {
		if rule != "" {
			var mode = GetRuleMode(rule)
			switch mode {
			case ModeUnknown:
				log.Warnf("%s: 不是IP或Domain规则", rule)
			default:
				num++
				r := new(routing)
				r.Data = rule
				r.Mode = mode
				switch rt {
				case TypeBlock:
					route.Block = append(route.Block, r)
				case TypeDirect:
					route.Direct = append(route.Direct, r)
				case TypeProxy:
					route.Proxy = append(route.Proxy, r)
				}
				if len(list) == 1 {
					log.Infof("%s: 添加一条%s规则", rule, mode)
				}
			}

		}
	}
	if num > 1 {
		log.Infof("共添加了 %d 条规则", num)
	}
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
	indexList := format.IndexDeal(key, len(rules))
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		r := rules[x]
		result = append(result, []string{
			tools.IntToStr(x + 1),
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
	length := len(rules)
	indexList := format.OtherIndex(key, length)
	if len(indexList) == length {
		return
	}
	defer route.save()
	result := make([]*routing, 0, len(indexList))
	for _, index := range indexList {
		result = append(result, rules[index])
	}
	switch rt {
	case TypeDirect:
		route.Direct = result
	case TypeProxy:
		route.Proxy = result
	case TypeBlock:
		route.Block = result
	}
	log.Info("删除了 [", length-len(result), "] 条规则")
}
