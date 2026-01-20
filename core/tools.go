// core/tools.go 提供核心工具函数，如索引处理、字符串处理等
package core

import (
	"encoding/json"
	"os"
)

// WriteJSON 将对象写入json文件
// 参数：
//	v：要写入的对象
//	path：目标文件路径
// 返回值：
//	error：返回错误信息，nil表示成功
// 实现细节：
//	1. 创建或覆盖指定路径的文件
//	2. 使用json.NewEncoder进行json编码
//	3. 设置缩进格式
//	4. 将对象编码后写入文件
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
