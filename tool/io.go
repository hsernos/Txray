package tool

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)



// WriteJSON 将对象写入json文件
func WriteJSON(v interface{}, path string) {
	
	file, e := os.Create(path)
	if e != nil {
		log.Fatal("文件创建失败", e.Error())
	}
	defer file.Close()
	jsonEncode := json.NewEncoder(file)
	jsonEncode.SetIndent("", "\t")
	err := jsonEncode.Encode(v)
	if err != nil {
		log.Fatal("编码错误", err.Error())
	}
}

// ReadJSON 读取json文件
func ReadJSON(path string, v interface{}) {
	
	file, e := os.Open(path)
	if e != nil {
		log.Fatal("文件打开失败", e.Error())
	}
	defer file.Close()
	decode := json.NewDecoder(file)
	decode.Decode(v)
}

// PathExists 判断文件或文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// ReadFile 读取文件
func ReadFile(path string) []string {
	data, _ := ioutil.ReadFile(path)
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

// 检查某个程序是否在PATH环境变量下
func CheckPATH(exec string) bool {
	for _, x := range strings.Split(Env("PATH"),":") {
		if PathExists(Join(x,exec)) || PathExists(Join(x,exec+".exe")) {
			return true
		}
	}
	return false
}