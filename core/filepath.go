package core

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	DataFile    = PathJoin(GetConfigDir(), "data.json")
	SettingFile = PathJoin(GetConfigDir(), "setting.toml")
	RoutingFile = PathJoin(GetConfigDir(), "routing.json")
)

// 获取配置文件所在目录
func GetConfigDir() string {
	dir := os.Getenv("TXRAY_HOME")
	if dir != "" && IsDir(dir) {
		return dir
	}
	return GetRunPath()
}

// 获取配置文件所在目录
func GetCoreDir() string {
	dir := os.Getenv("CORE_HOME")
	if dir != "" && IsDir(dir) {
		return dir
	}
	return GetRunPath()
}

// 获取当前进程的可执行文件的路径名
func GetRunPath() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
	//return filepath.Dir("C:\\Users\\hsernos\\Desktop\\")
}

// 判断路径是否正确且为文件夹
func IsDir(path string) bool {
	i, err := os.Stat(path)
	if err == nil {
		return i.IsDir()
	}
	return false
}

// 路径拼接
func PathJoin(elem ...string) string {
	if len(elem) > 0 {
		if strings.HasSuffix(elem[0], string(os.PathSeparator)) {
			elem[0] = strings.TrimRight(elem[0], string(os.PathSeparator))
		}
	}
	return strings.Join(elem, string(os.PathSeparator))
}
