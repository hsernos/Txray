// core/manage/list.go 负责节点、订阅等列表的管理操作
package manage

// HasIn 函数用于检查一个整数是否存在于整数切片中
// 参数:
//   - index: 要检查的整数
//   - indexList: 整数切片
// 返回值:
//   - 如果 index 存在于 indexList 中，返回 true；否则返回 false
// 实现细节:
//   - 遍历 indexList，逐个与 index 比较
//   - 如果找到相等的，立即返回 true
//   - 如果遍历结束仍未找到，返回 false
func HasIn(index int, indexList []int) bool {
	for _, i := range indexList {
		if i == index {
			return true
		}
	}
	return false
}
