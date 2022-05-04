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
func GetConfigDir() string {
	dir := os.Getenv("TXRAY_HOME")
	if dir != "" && IsDir(dir) {
		return dir
	}
	return GetRunPath()
}


// 获取当前进程的可执行文件的路径名
func GetRunPath() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
}

// 判断路径是否正确且为文件夹
func IsDir(path string) bool {
	i, err := os.Stat(path)
	if err == nil {
		return i.IsDir()
	}
	return false
}
