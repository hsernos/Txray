// core/setting/base.go 负责基础设置项的定义与操作
package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

// Mixed 返回当前混合代理端口
func Mixed() int {
	return viper.GetInt(key.Mixed)
}

// SetMixed 设置混合代理端口
// port: 代理端口，取值范围为 1~65535
func SetMixed(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("混合代理端口取值 1~65535")
	}
	viper.Set(key.Mixed, port)
	return viper.WriteConfig()
}

// Socks 返回当前socks代理端口
func Socks() int {
	return viper.GetInt(key.Socks)
}

// SetSocks 设置socks代理端口
// port: 代理端口，取值范围为 1~65535
func SetSocks(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("socks端口取值 1~65535")
	}
	viper.Set(key.Socks, port)
	return viper.WriteConfig()
}

// Http 返回当前http代理端口
func Http() int {
	return viper.GetInt(key.Http)
}

// SetHttp 设置http代理端口
// port: 代理端口，取值范围为 0~65535
func SetHttp(port int) error {
	if port < 0 || port > 65535 {
		return errors.New("http端口取值 0~65535")
	}
	viper.Set(key.Http, port)
	return viper.WriteConfig()
}

// UDP 返回当前UDP状态
func UDP() bool {
	return viper.GetBool(key.UDP)
}

// SetUDP 设置UDP状态
// status: true开启UDP支持，false关闭
func SetUDP(status bool) error {
	viper.Set(key.UDP, status)
	return viper.WriteConfig()
}

// Sniffing 返回当前是否开启流量嗅探
func Sniffing() bool {
	return viper.GetBool(key.Sniffing)
}

// SetSniffing 设置流量嗅探状态
// status: true开启嗅探，false关闭
func SetSniffing(status bool) error {
	viper.Set(key.Sniffing, status)
	return viper.WriteConfig()
}

// FromLanConn 返回当前是否允许局域网连接
func FromLanConn() bool {
	return viper.GetBool(key.FromLanConn)
}

// SetFromLanConn 设置是否允许局域网连接
// status: true允许，false不允许
func SetFromLanConn(status bool) error {
	viper.Set(key.FromLanConn, status)
	return viper.WriteConfig()
}

// Mux 返回当前是否开启多路复用
func Mux() bool {
	return viper.GetBool(key.Mux)
}

// SetMux 设置多路复用状态
// status: true开启多路复用，false关闭
func SetMux(status bool) error {
	viper.Set(key.Mux, status)
	return viper.WriteConfig()
}

// Pid 返回当前进程ID
func Pid() int {
	return viper.GetInt(key.PID)
}

// SetPid 设置进程ID
// pid: 进程ID
func SetPid(pid int) error {
	viper.Set(key.PID, pid)
	return viper.WriteConfig()
}

// AllowInsecure 返回当前是否允许不安全连接
func AllowInsecure() bool {
	return viper.GetBool(key.AllowInsecure)
}

// SetAllowInsecure 设置是否允许不安全连接
// status: true允许不安全连接，false不允许
func SetAllowInsecure(status bool) error {
	viper.Set(key.AllowInsecure, status)
	return viper.WriteConfig()
}

// VersionMin 返回当前版本最小值
func VersionMin() string {
	return viper.GetString(key.VersionMin)
}

// SetVersionMin 设置版本最小值
// min: 版本最小值，为空字符串表示不限制
func SetVersionMin(min string) error {
	viper.Set(key.VersionMin, min)
	return viper.WriteConfig()
}

// VersionMax 返回当前版本最大值
func VersionMax() string {
	return viper.GetString(key.VersionMax)
}

// SetVersionMax 设置版本最大值
// max: 版本最大值，为空字符串表示不限制
func SetVersionMax(max string) error {
	viper.Set(key.VersionMax, max)
	return viper.WriteConfig()
}
