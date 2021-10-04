package xray

import (
	"Txray/core"
	"Txray/log"
	"os"
	"path/filepath"
	"strings"
)

// 检查xray程序和资源文件 geoip.dat、geosite.dat是否存在且在同一目录
func CheckFile() bool {
	xrayPath := XrayPath()
	if xrayPath == "" {
		log.Error("在", core.GetCoreDir(), " 下没有找到xray程序")
		log.Error("请在 https://github.com/XTLS/Xray-core/releases 下载最新版本")
		log.Error("并将解压后的文件夹或所有文件移动到 ", core.GetCoreDir(), " 下")
		return false
	} else {
		path := filepath.Dir(xrayPath)
		_, err1 := os.Stat(PathJoin(path, "geoip.dat"))
		_, err2 := os.Stat(PathJoin(path, "geosite.dat"))
		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			log.Error("在 ", path, " 下没有找到xray程序的资源文件 geoip.dat 或v geosite.dat")
			log.Error("请在 https://github.com/XTLS/Xray-core/releases 下载最新版本")
			log.Error("并将缺失的文件移动到 ", path, " 下")
			return false
		} else {
			return true
		}
	}
}

func XrayPath() string {
	files, _ := FindFileByName(core.GetCoreDir(), "xray", ".exe")
	if len(files) == 0 {
		return ""
	}
	return files[0]
}

// 遍历目录，查找文件
func FindFileByName(root, name, ext string) ([]string, error) {
	root = strings.TrimRight(root, string(os.PathSeparator))
	paths, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	objList := make([]string, 0)
	for _, p := range paths {
		absPath := root + string(os.PathSeparator) + p.Name()
		if p.IsDir() {
			o, err := FindFileByName(absPath, name, ext)
			if err != nil {
				return nil, err
			}
			objList = append(objList, o...)
		} else {
			if p.Name() == name || p.Name() == name+ext {
				objList = append(objList, absPath)
			}
		}
	}
	return objList, nil
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
