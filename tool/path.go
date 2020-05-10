package tool

import (
	"os"
	"os/user"
	"runtime"
	"strings"
)

// 获取用户家目录
func Home() (string, error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	} else {
		return "", err
	}
}

// 拼接目录
func Join(elem ...string) string {
	return strings.Join(elem, string(os.PathSeparator))
}

// 检查某个程序是否在PATH环境变量下
func CheckPATH(exec string) string {
	path := os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		for _, x := range strings.Split(path, ";") {
			if PathExists(Join(x, exec) + ".exe") {
				return x
			}
		}
	} else {
		for _, x := range strings.Split(path, ":") {
			if PathExists(Join(x, exec)) {
				return x
			}
		}
	}
	return ""
}
