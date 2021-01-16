package format

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

// 生成num个相同字符组成的字符串
func RepeatChar(ch byte, num int) string {
	var buf bytes.Buffer
	for i := 1; i < num; i++ {
		buf.WriteByte(ch)
	}
	return buf.String()
}

// 生成同str一样宽度相同字符组成的字符串
func SameLenRepeatChar(ch byte, str string) string {
	width := runewidth.StringWidth(str)
	return RepeatChar(ch, width+1)
}

// 所有字符串中最大宽度
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
func ShowTopBottomSepLine(ch byte, str ...string) {
	width := MaxWidth(str...) + 1
	fmt.Println(RepeatChar(ch, width))
	fmt.Println(strings.Join(str, "\n"))
	fmt.Println(RepeatChar(ch, width))
}
