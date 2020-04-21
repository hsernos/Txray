package tool

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"strconv"
	"strings"
)

// GetTable 将数据转换成表格形式
func GetTable(headers []string, datas ...[]string)  {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAlignment(1)
	for _, v := range datas {
		table.Append(v)
	}
	table.Render()
}
// IndexDeal 类似切片，返回索引列表
func IndexDeal(key string, max int) []int {

	if key == "all" {
		result := make([]int, 0, max)
		for i := 0; i < max; i++ {
			if i >= 0 {
				result = append(result, i)
			}
		}
		return result
	} else if strings.Contains(key, ",") {
		result := make([]int, 0, 30)
		for _, x := range strings.Split(key, ",") {
			for _, y := range IndexDeal(x, max) {
				result = append(result, y)
			}
		}
		result = RemoveRep(result)
		sort.Ints(result)
		return result
	} else if strings.Contains(key, "-") {
		l := strings.Split(key, "-")
		if len(l) == 2 && IsInt(l[0]) && IsInt(l[1]) {
			a, b := StrToInt(l[0]), StrToInt(l[1])
			if a > b {
				a, b = b, a
			}
			result := make([]int, 0, b-a+1)
			for i := a; i <= b; i++ {
				if i < max && i >= 0 {
					result = append(result, i)
				}
			}
			return result
		}
	} else if IsInt(key) {
		num := StrToInt(key)
		if num < max && num >= 0 {
			return []int{int(num)}
		}
	}
	return []int{}
}

// IsInt 是否为整数
func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// IsUint 是否为整数
func IsUint(s string) bool {
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

// StrToUint string --> uint
func StrToUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}

// StrToInt string --> int
func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// IntToStr int --> string
func IntToStr(i int) string {
	return strconv.Itoa(i)
}
// UintToStr uint --> string
func UintToStr(i uint) string {
	return strconv.Itoa(int(i))
}


// BoolToStr uint --> string
func BoolToStr(b bool) string {
	return strconv.FormatBool(b)
}