package tool

import (
	"os"
	"os/user"
	"path"
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
func Join(elem ...string) string  {
	return path.Join(elem...)
}

// 获取环境变量
func Env(key string) string  {
	return os.Getenv(key)
}

func ConfigPath() string  {
	p := Env("V3RAY")
	if p != "" {
		return p
	} else {
		p,_ = Home()
		if p != "" {
			return Join(p,"v3ray")
		}
		return ""
	}
}
