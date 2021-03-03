package xray

import (
	"Txray/core/config"
	"Txray/log"
	"Txray/tools"
	"path/filepath"
)

// 检查xray程序和资源文件 geoip.dat、geosite.dat是否存在且在同一目录
func CheckFile() bool {
	xrayPath := GetPath()
	if xrayPath == "" {
		log.Error("在", config.GetCoreDir(), " 下没有找到xray程序")
		log.Error("请在 https://github.com/XTLS/Xray-core/releases 下载最新版本")
		log.Error("并将解压后的文件夹或所有文件移动到 ", config.GetCoreDir(), " 下")
		return false
	} else {
		path := filepath.Dir(GetPath())
		if tools.IsFile(tools.PathJoin(path, "geoip.dat")) && tools.IsFile(tools.PathJoin(path, "geosite.dat")) {
			return true
		} else {
			log.Error("在 ", path, " 下没有找到xray程序的资源文件 geoip.dat 和 geosite.dat")
			log.Error("请在 https://github.com/XTLS/Xray-core/releases 下载最新版本")
			log.Error("并将缺失的文件移动到 ", path, " 下")
			return false
		}
	}
}

// 查找xray程序所在绝对路径
func GetPath() string {
	path := config.GetCoreDir()
	files, _ := tools.FindFileByName(path, "xray", ".exe")
	if len(files) == 0 {
		return ""
	}
	return files[0]
}
