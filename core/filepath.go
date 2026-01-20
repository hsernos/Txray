// core/filepath.go 负责路径相关的工具函数
package core

import (
	"os"
	"path/filepath"
)

var (
	DataFile    = filepath.Join(GetConfigDir(), "data.json")
	SettingFile = filepath.Join(GetConfigDir(), "setting.toml")
	RoutingFile = filepath.Join(GetConfigDir(), "routing.json")
	LogFile = filepath.Join(GetConfigDir(), "xray_access.log")
)

// 获取配置文件所在目录
// 如果环境变量 TXRAY_HOME 被设置且是一个有效的目录，返回该目录；否则，返回当前可执行文件所在目录
func GetConfigDir() string {
	dir := os.Getenv("TXRAY_HOME")
	if dir != "" && IsDir(dir) {
		return dir
	}
	return GetRunPath()
}

// 获取当前进程的可执行文件的路径名
// 返回当前可执行文件所在的目录
func GetRunPath() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
}

// 判断路径是否正确且为文件夹
// 如果给定路径存在且是一个目录，返回 true；否则，返回 false
func IsDir(path string) bool {
	i, err := os.Stat(path)
	if err == nil {
		return i.IsDir()
	}
	return false
}
