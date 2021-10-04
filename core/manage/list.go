package manage

func HasIn(index int, indexList []int) bool {
	for _, i := range indexList {
		if i == index {
			return true
		}
	}
	return false
}
