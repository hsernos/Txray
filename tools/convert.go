package tools

import "strconv"

// 是否为整数
func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// 是否为非负整数
func IsUint(s string) bool {
	_, err := strconv.ParseUint(s, 10, 64)
	return err == nil
}

// 是否为网络端口，port：0~65535
func IsNetPort(s string) bool {
	if IsUint(s) {
		return StrToUint(s) <= 65535
	}
	return false
}

// string --> uint
func StrToUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}

// string --> int
func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// int --> string
func IntToStr(i int) string {
	return strconv.Itoa(i)
}

// uint --> string
func UintToStr(i uint) string {
	return strconv.Itoa(int(i))
}

// BoolToStr uint --> string
func BoolToStr(b bool) string {
	return strconv.FormatBool(b)
}
