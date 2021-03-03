package routing

import (
	"regexp"
	"strings"
)

type Mode string

const (
	ModeIP      Mode = "IP"
	ModeDomain       = "Domain"
	ModeUnknown      = "Unknown"
)

// 判断是IP规则还是域名规则
// IP|Unknown|Domain
func GetRuleMode(str string) Mode {
	if strings.HasPrefix(str, "regexp:") {
		return ModeDomain
	}
	if strings.HasPrefix(str, "domain:") {
		return ModeDomain
	}
	if strings.HasPrefix(str, "fill:") {
		return ModeDomain
	}
	if strings.HasPrefix(str, "geosite:") {
		return ModeDomain
	}
	if strings.HasPrefix(str, "geosite:") {
		return ModeDomain
	}
	if strings.HasPrefix(str, "geoip:") {
		return ModeIP
	}
	pattern := "^ext:[a-zA-Z0-9_]*?(ip|site).dat:"
	re, _ := regexp.Compile(pattern)
	result := re.FindStringSubmatch(str)
	if len(result) == 2 {
		if result[1] == "ip" {
			return ModeIP
		} else {
			return ModeDomain
		}
	}
	pattern = `(^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(/([1-9]|[1-2][0-9]|3[0-2]){1})?$`
	re, _ = regexp.Compile(pattern)
	result = re.FindStringSubmatch(str)
	if len(result) != 0 {
		return ModeIP
	}
	pattern = `^([a-zA-Z0-9][-a-zA-Z0-9]{0,62})(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`
	re, _ = regexp.Compile(pattern)
	result = re.FindStringSubmatch(str)
	if len(result) != 0 && len(str) < 256 {
		return ModeDomain
	}
	return ModeUnknown
}
