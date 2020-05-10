package tool

import "strconv"

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
