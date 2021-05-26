package tools

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// 判断路径是否正确，是否存在文件或文件夹
func IsPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 判断路径是否正确且为文件夹
func IsDir(path string) bool {
	i, err := os.Stat(path)
	if err == nil {
		return i.IsDir()
	}
	return false
}

// 判断路径是否正确且为文件
func IsFile(path string) bool {
	i, err := os.Stat(path)
	if err == nil {
		return !i.IsDir()
	}
	return false
}

// 将路径分隔符装换成Linux
func ToLinuxPathSeparator(basePath string) string {
	return strings.ReplaceAll(basePath, "\\", "/")
}

// 获取当前进程的可执行文件的路径名
func GetRunPath() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
	//return filepath.Dir("/Users/hsernos/Desktop/")
}

// 获取用户家目录
func GetHome() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir, nil
	} else {
		return "", err
	}
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

// 遍历目录，获取所有的文件和目录
func GetFilesAndDirs(path string) ([]string, []string, error) {
	path = strings.TrimRight(path, string(os.PathSeparator))
	paths, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}
	dirs := make([]string, 0)
	files := make([]string, 0)
	for _, p := range paths {
		absPath := path + string(os.PathSeparator) + p.Name()
		if p.IsDir() {
			dirs = append(dirs, absPath)
			fs, ds, err := GetFilesAndDirs(absPath)
			if err != nil {
				return nil, nil, err
			}
			dirs = append(dirs, ds...)
			files = append(files, fs...)
		} else {
			files = append(files, absPath)
		}
	}
	return files, dirs, nil
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

// 在PATH环境变量中查找程序的绝对路径
func FindEXEPathInEnv(exec string) string {
	path := os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		for _, x := range strings.Split(path, ";") {
			if IsFile(PathJoin(x, exec)+".exe") || IsFile(PathJoin(x, exec)+".EXE") {
				return x
			}
		}
	} else {
		for _, x := range strings.Split(path, ":") {
			if IsFile(PathJoin(x, exec)) {
				return x
			}
		}
	}
	return ""
}
