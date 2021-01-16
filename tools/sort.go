package tools

type Sort struct {
	Index int
	Value float64
}

type Sorts []Sort

//Len()
func (s Sorts) Len() int {
	return len(s)
}

//Less()
func (s Sorts) Less(i, j int) bool {
	return s[i].Value < s[j].Value
}

//Swap()
func (s Sorts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
