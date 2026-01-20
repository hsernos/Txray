// core/key.go 负责全局 key 常量的定义
package core

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Range 返回一个从 srart 到 end 的整数切片，包含 srart 和 end
func Range(srart, end int) []int {
	result := make([]int, 0)
	if srart <= end {
		for srart <= end {
			result = append(result, srart)
			srart++
		}
	} else {
		for srart >= end {
			result = append(result, srart)
			srart--
		}
	}
	return result
}

// IndexList 根据给定的 key 和最大值 max 返回一个整数切片
// key 可以是 "all" 或者以逗号分隔的数字和范围
// 如果 key 是 "all"，则返回从 1 到 max 的所有整数
// 如果 key 是数字，则返回该数字
// 如果 key 是范围（如 "1-5"），则返回范围内的所有整数
func IndexList(key string, max int) []int {
	if max == 0 {
		return []int{}
	}
	if key == "all" {
		return Range(1, max)
	}
	result := make([]int, 0)
	for _, item := range strings.Split(key, ",") {
		item = strings.Trim(item, " ")
		re1 := "^[1-9][0-9]*$"
		re2 := "(^[0-9]*)-([0-9]*$)"
		if re, _ := regexp.Compile(re1); re.MatchString(item) {
			i, _ := strconv.Atoi(item)
			if i > 0 && i <= max {
				result = append(result, i)
			}
			continue
		}
		if re, _ := regexp.Compile(re2); re.MatchString(item) {
			start := 1
			end := max
			s := re.FindStringSubmatch(item)
			if s[1] != "" {
				start, _ = strconv.Atoi(s[1])
			}
			if s[2] != "" {
				end, _ = strconv.Atoi(s[2])
			}
			if start > end {
				start, end = end, start
			}
			if start > max || end < 1 {
				continue
			}
			if start < 1 {
				start = 1
			}
			if end > max {
				end = max
			}
			result = append(result, Range(start, end)...)
			continue
		}
	}
	result = RemoveRepByMap(result)
	sort.Ints(result)
	return result
}

// RemoveRepByMap 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// Reverse 反转整数切片
func Reverse(slc []int) []int {
	for i := 0; i < len(slc)/2; i++ {
		j := len(slc) - i - 1
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}
