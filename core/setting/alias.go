// core/setting/alias.go 负责别名相关的设置与操作
package setting

import (
	"Txray/core"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

// Alias 结构体表示一个别名，包括名称和命令
type Alias struct {
	Name string // 别名名称
	Cmd  string // 别名对应的命令
}

// NewAlias 创建一个新的 Alias 实例
func NewAlias(name, cmd string) *Alias {
	return &Alias{
		Name: name,
		Cmd:  cmd,
	}
}

// GetCmd 返回别名对应命令的字符串切片
func (a *Alias) GetCmd() [][]string {
	result := make([][]string, 0)
	// 以 "|" 分割别名命令，并将每个命令转换为字符串切片
	for _, line := range strings.Split(a.Cmd, "|") {
		cmd := strings.Fields(line)
		if cmd != nil {
			result = append(result, cmd)
		}
	}
	return result
}

// AliasList 返回所有别名的列表
func AliasList() []*Alias {
	result := make([]*Alias, 0)
	// 获取配置中所有别名
	for k, v := range viper.GetStringMapString("alias") {
		result = append(result, NewAlias(k, v))
	}
	// 如果别名数量小于等于1，直接返回
	if len(result) <= 1 {
		return result
	}
	// 对别名按名称进行排序
	for i := 1; i < len(result); i++ {
		preIndex := i - 1
		current := result[i]
		// 插入排序法
		for preIndex >= 0 && result[preIndex].Name > current.Name {
			result[preIndex+1] = result[preIndex]
			preIndex -= 1
		}
		result[preIndex+1] = current
	}
	return result
}

// AddAlias 添加一个新的别名
func AddAlias(key, cmd string) {
	// 检查别名是否合法
	if ok, _ := regexp.MatchString("^[^ ]*$", key); ok {
		v := viper.Get("alias")
		// 如果别名不存在，则设置默认值
		if v == nil {
			viper.SetDefault("alias."+key, cmd)
		} else {
			// 否则更新现有别名
			v.(map[string]interface{})[key] = cmd
		}
		// 写入配置文件
		viper.WriteConfig()
	}
}

// DelAlias 删除一个别名
func DelAlias(key string) []string {
	result := make([]string, 0)
	aliasList := AliasList()
	// 获取要删除的别名在列表中的索引
	indexList := core.IndexList(key, len(aliasList))
	if len(indexList) == 0 {
		// 如果未找到索引，直接删除
		delAlias(key)
		result = append(result, key)
	} else {
		// 根据索引删除对应的别名
		for _, index := range indexList {
			delAlias(aliasList[index-1].Name)
			result = append(result, aliasList[index-1].Name)
		}
	}
	return result
}

// delAlias 实际执行别名删除操作
func delAlias(key string) {
	delete(viper.Get("alias").(map[string]interface{}), key)
	viper.WriteConfig()
}
