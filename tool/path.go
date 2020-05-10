package tool

import (
	"os"
	"os/user"
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
func CheckPATH(exec string) bool {
	path := os.Getenv("PATH")
	if strings.IndexAny(path, ";") >= 0 {
		for _, x := range strings.Split(path, ";") {
			if PathExists(Join(x, exec) + ".exe") {
				return true
			}
		}
	} else {
		for _, x := range strings.Split(path, ":") {
			if PathExists(Join(x, exec)) {
				return true
			}
		}
	}
	return false
}
