package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

func RoutingStrategy() string {
	return viper.GetString(key.RoutingStrategy)
}

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

func RoutingBypass() bool {
	return viper.GetBool(key.RoutingBypass)
}

func SetRoutingBypass(status bool) error {
	viper.Set(key.RoutingBypass, status)
	return viper.WriteConfig()
}
