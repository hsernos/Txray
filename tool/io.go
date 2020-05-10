package tool

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	log "v3ray/logger"
)

// WriteJSON 将对象写入json文件
func WriteJSON(v interface{}, path string) {

	file, e := os.Create(path)
	if e != nil {
		log.Error("文件创建失败", e.Error())
	}
	defer file.Close()
	jsonEncode := json.NewEncoder(file)
	jsonEncode.SetIndent("", "\t")
	err := jsonEncode.Encode(v)
	if err != nil {
		log.Error("编码错误", err.Error())
	}
}

// ReadJSON 读取json文件
func ReadJSON(path string, v interface{}) {

	file, e := os.Open(path)
	if e != nil {
		log.Error("文件打开失败", e.Error())
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
