package tools

import (
	"encoding/json"
	"os"
	"strings"
)

// 将对象写入json文件
func WriteJSON(v interface{}, path string) error {

	file, e := os.Create(path)
	if e != nil {
		return e
	}
	defer file.Close()
	jsonEncode := json.NewEncoder(file)
	jsonEncode.SetIndent("", "\t")
	return jsonEncode.Encode(v)
}

// 读取json文件
func ReadJSON(path string, v interface{}) error {
	file, e := os.Open(path)
	if e != nil {
		return e
	}
	defer file.Close()
	decode := json.NewDecoder(file)
	return decode.Decode(v)
}

// ReadFile 读取文件
func ReadFile(path string) []string {
	data, _ := os.ReadFile(path)
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}
