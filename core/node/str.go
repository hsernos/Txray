// core/node/str.go 负责节点相关的字符串处理工具函数
package node

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

// 生成num个相同字符组成的字符串
// ch：要重复的字符
// num：重复的次数
// 返回由num个字符ch组成的字符串
func RepeatChar(ch byte, num int) string {
	return strings.Repeat(string(ch), num)
}

// 所有字符串中最大宽度
// str：要计算的字符串
// 返回所有字符串中最大宽度
func MaxWidth(str ...string) int {
	max := 0
	for _, s := range str {
		width := runewidth.StringWidth(s)
		if width > max {
			max = width
		}
	}
	return max
}

// 添加上下的分割线
// ch：分割线使用的字符
// str：要显示的字符串
// 根据传入的字符串，显示带有分割线的字符串
func ShowTopBottomSepLine(ch byte, str ...string) {
	width := MaxWidth(str...)
	fmt.Println(RepeatChar(ch, width))
	fmt.Println(strings.Join(str, "\n"))
	fmt.Println(RepeatChar(ch, width))
}
