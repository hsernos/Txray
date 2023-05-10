package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

func Socks() int {
	return viper.GetInt(key.Socks)
}

func SetSocks(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("socks端口取值 1~65535")
	}
	viper.Set(key.Socks, port)
	return viper.WriteConfig()
}

func Http() int {
	return viper.GetInt(key.Http)
}

func SetHttp(port int) error {
	if port < 0 || port > 65535 {
		return errors.New("http端口取值 0~65535")
	}
	viper.Set(key.Http, port)
	return viper.WriteConfig()
}

func UDP() bool {
	return viper.GetBool(key.UDP)
}

func SetUDP(status bool) error {
	viper.Set(key.UDP, status)
	return viper.WriteConfig()
}

func Sniffing() bool {
	return viper.GetBool(key.Sniffing)
}

func SetSniffing(status bool) error {
	viper.Set(key.Sniffing, status)
	return viper.WriteConfig()
}

func FromLanConn() bool {
	return viper.GetBool(key.FromLanConn)
}

func SetFromLanConn(status bool) error {
	viper.Set(key.FromLanConn, status)
	return viper.WriteConfig()
}

func Mux() bool {
	return viper.GetBool(key.Mux)
}

func SetMux(status bool) error {
	viper.Set(key.Mux, status)
	return viper.WriteConfig()
}

func Pid() int {
	return viper.GetInt(key.PID)
}

func SetPid(pid int) error {
	viper.Set(key.PID, pid)
	return viper.WriteConfig()
}

func AllowInsecure() bool {
	return viper.GetBool(key.AllowInsecure)
}

func SetAllowInsecure(status bool) error {
	viper.Set(key.AllowInsecure, status)
	return viper.WriteConfig()
}