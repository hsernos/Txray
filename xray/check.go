package xray

import (
	"Txray/log"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const CoreName string = "xray"

var XrayPath = ""

func init() {
	checkXrayFile()
	checkResource()
}

// 检查xray程序
func checkXrayFile() {
	// 1.检查环境变量CORE_HOME目录下
	xrayPath := os.Getenv("CORE_HOME")
	if xrayPath != "" {
		if IsExistExe(xrayPath, CoreName) {
			XrayPath = filepath.Join(xrayPath, CoreName)
			return
		}
	}
	// 2.检查当前可执行文件目录下(递归检查)
	path, _ := os.Executable()
	files, _ := FindFileByName(filepath.Dir(path), "xray", ".exe")
	if len(files) != 0 {
		XrayPath = files[0]
		return
	}
	// 3.检查PATH环境变量
	if temp := getExePath(CoreName); temp != "" {
		XrayPath = temp
		return
	}
	// 提示信息
	log.Error("在 ", filepath.Dir(path), " 下没有找到xray程序")
	log.Error("请在 https://github.com/XTLS/Xray-core/releases 下载最新版本")
	log.Error("并将解压后的文件夹或所有文件移动到 ", filepath.Dir(path), " 下")
	os.Exit(0)
}

// 检查xray程序的资源文件
func checkResource() {
	var baseDir []string = make([]string, 0)
	baseDir = append(baseDir, os.Getenv("XRAY_LOCATION_ASSET"))
	baseDir = append(baseDir, os.Getenv("xray.location.asset"))
	baseDir = append(baseDir, filepath.Dir(XrayPath))
	for _, dir := range baseDir {
		if dir != "" {
			if IsExistFile(filepath.Join(dir, "geoip.dat")) && IsExistFile(filepath.Join(dir, "geosite.dat")) {
				return
			}
		}
	}
	log.Error(fmt.Sprintf("在 %s 目录下没有找到资源文件 geoip.dat 和 geosite.dat", filepath.Dir(XrayPath)))
	log.Error("或者配置资源文件的环境变量 XRAY_LOCATION_ASSET")
	os.Exit(0)
}

func IsExistFile(file string) bool {
	fp, err := os.Stat(file)
	return err == nil && !fp.IsDir()
}

// 检查dirPath目录下是否存在filename程序
func IsExistExe(dirPath, filename string) bool {
	if runtime.GOOS == "windows" {
		fp, err := os.Stat(filepath.Join(dirPath, filename+".exe"))
		if err == nil && !fp.IsDir() {
			return true
		}
	}
	fp, err := os.Stat(filepath.Join(dirPath, filename))
	return err == nil && !fp.IsDir()
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

// 查找Path中可执行文件的路径
func getExePath(name string) string {
	data := os.Getenv("PATH")
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}
	for _, x := range strings.Split(data, sep) {
		if strings.TrimSpace(x) != "" {
			if IsExistExe(strings.TrimSpace(x), name) {
				return filepath.Join(strings.TrimSpace(x), name)
			}
		}
	}
	return ""
}
