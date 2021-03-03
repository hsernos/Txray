package config

import (
	"Txray/tools"
	"os"
)

// 获取配置文件所在目录
func GetConfigDir() string {
	dir := os.Getenv("TXRAY_HOME")
	if dir != "" && tools.IsDir(dir) {
		return dir
	}
	return tools.GetRunPath()
}

// 获取配置文件所在目录
func GetCoreDir() string {
	dir := os.Getenv("CORE_HOME")
	if dir != "" && tools.IsDir(dir) {
		return dir
	}
	return tools.GetRunPath()
}

var (
	Routing     = tools.PathJoin(GetConfigDir(), "routing.json")
	NodesAndSub = tools.PathJoin(GetConfigDir(), "node.json")
	Setting     = tools.PathJoin(GetConfigDir(), "setting.json")
)
