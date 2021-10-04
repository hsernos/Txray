package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

func TestUrl() string {
	return viper.GetString(key.TestURL)
}

func SetTestUrl(url string) error {
	viper.Set(key.TestURL, url)
	return viper.WriteConfig()
}

func TestTimeout() int {
	return viper.GetInt(key.TestTimeout)
}

func SetTestTimeout(timeout int) error {
	if timeout < 0 {
		return errors.New("取值不能小于0")
	}
	viper.Set(key.TestTimeout, timeout)
	return viper.WriteConfig()
}

func RunBefore() string {
	return viper.GetString(key.RunBefore)
}

func SetRunBefore(cmd string) error {
	viper.Set(key.RunBefore, cmd)
	return viper.WriteConfig()
}
