// core/setting/routing.go 负责路由相关设置项的定义与操作
package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

// RoutingStrategy 获取当前的路由策略
func RoutingStrategy() string {
	return viper.GetString(key.RoutingStrategy)
}

// SetRoutingStrategy 设置路由策略
// mode: 路由策略模式
// 1 - AsIs
// 2 - IPIfNonMatch
// 3 - IPOnDemand
// 返回错误信息（如果有的话）
func SetRoutingStrategy(mode int) error {
	switch mode {
	case 1:
		viper.Set(key.RoutingStrategy, "AsIs")
	case 2:
		viper.Set(key.RoutingStrategy, "IPIfNonMatch")
	case 3:
		viper.Set(key.RoutingStrategy, "IPOnDemand")
	default:
		return errors.New("没有对应的路由策略")
	}
	return viper.WriteConfig()
}

// RoutingBypass 获取当前是否绕过路由
func RoutingBypass() bool {
	return viper.GetBool(key.RoutingBypass)
}

// SetRoutingBypass 设置是否绕过路由
// status: true 表示绕过路由，false 表示不绕过
// 返回错误信息（如果有的话）
func SetRoutingBypass(status bool) error {
	viper.Set(key.RoutingBypass, status)
	return viper.WriteConfig()
}
