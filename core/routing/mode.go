// core/routing/mode.go 负责路由模式的定义与相关操作
package routing

import (
	"regexp"
	"strings"
)

// Mode 表示路由规则的模式，可能的值有 IP 和 Domain
type Mode string

const (
	ModeIP     Mode = "IP"     // IP 表示 IP 地址规则
	ModeDomain Mode = "Domain" // Domain 表示域名规则
)

// 判断是IP规则还是域名规则
// IP|Domain
// GetRuleMode 根据给定的规则字符串返回其对应的模式类型
// 参数:
//   - str: 规则字符串
// 返回值:
//   - 规则的模式类型，可能是 ModeIP 或 ModeDomain
func GetRuleMode(str string) Mode {
	if strings.HasPrefix(str, "geoip:") {
		return ModeIP
	}
	if strings.Contains(str, "ip.dat:") {
		return ModeIP
	}
	pattern := `(?:^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(?:/(?:[1-9]|[1-2][0-9]|3[0-2]){1})?$`
	re, _ := regexp.Compile(pattern)
	if re.MatchString(str) {
		return ModeIP
	}
	return ModeDomain
}
