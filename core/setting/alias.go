package setting

import (
	"Txray/core"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type Alias struct {
	Name string
	Cmd  string
}

func NewAlias(name, cmd string) *Alias {
	return &Alias{
		Name: name,
		Cmd:  cmd,
	}
}

func (a *Alias) GetCmd() [][]string {
	result := make([][]string, 0)
	for _, line := range strings.Split(a.Cmd, "|") {
		cmd := strings.Fields(line)
		if cmd != nil {
			result = append(result, cmd)
		}
	}
	return result
}

func AliasList() []*Alias {
	result := make([]*Alias, 0)
	for k, v := range viper.GetStringMapString("alias") {
		result = append(result, NewAlias(k, v))
	}
	if len(result) <= 1 {
		return result
	}
	for i := 1; i < len(result); i++ {
		preIndex := i - 1
		current := result[i]
		for preIndex >= 0 && result[preIndex].Name > current.Name {
			result[preIndex+1] = result[preIndex]
			preIndex -= 1
		}
		result[preIndex+1] = current
	}
	return result
}

func AddAlias(key, cmd string) {
	if ok, _ := regexp.MatchString("^[^ ]*$", key); ok {
		v := viper.Get("alias")
		if v == nil {
			viper.SetDefault("alias."+key, cmd)
		} else {
			v.(map[string]interface{})[key] = cmd
		}
		viper.WriteConfig()
	}
}

func DelAlias(key string) []string {
	result := make([]string, 0)
	aliasList := AliasList()
	indexList := core.IndexList(key, len(aliasList))
	if len(indexList) == 0 {
		delAlias(key)
		result = append(result, key)
	} else {
		for _, index := range indexList {
			delAlias(aliasList[index-1].Name)
			result = append(result, aliasList[index-1].Name)
		}
	}
	return result
}

func delAlias(key string) {
	delete(viper.Get("alias").(map[string]interface{}), key)
	viper.WriteConfig()
}
