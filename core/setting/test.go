// core/setting/test.go 负责测试相关设置项的定义与操作
package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

// TestUrl 返回测试用例的目标URL
func TestUrl() string {
	return viper.GetString(key.TestURL)
}

// SetTestUrl 设置测试用例的目标URL
func SetTestUrl(url string) error {
	viper.Set(key.TestURL, url)
	return viper.WriteConfig()
}

// TestTimeout 返回测试超时时间，单位为秒
func TestTimeout() int {
	return viper.GetInt(key.TestTimeout)
}

// SetTestTimeout 设置测试超时时间，单位为秒
func SetTestTimeout(timeout int) error {
	if timeout < 0 {
		return errors.New("取值不能小于0")
	}
	viper.Set(key.TestTimeout, timeout)
	return viper.WriteConfig()
}

// TestMinTime 返回测试最小执行时间，单位为秒
func TestMinTime() int {
	return viper.GetInt(key.TestMinTime)
}

// SetTestMinTime 设置测试最小执行时间，单位为秒
func SetTestMinTime(timeout int) error {
	if timeout < 0 {
		return errors.New("取值不能小于0")
	}
	viper.Set(key.TestMinTime, timeout)
	return viper.WriteConfig()
}

// RunBefore 返回测试执行前需要运行的命令
func RunBefore() string {
	return viper.GetString(key.RunBefore)
}

// SetRunBefore 设置测试执行前需要运行的命令
func SetRunBefore(cmd string) error {
	viper.Set(key.RunBefore, cmd)
	return viper.WriteConfig()
}
