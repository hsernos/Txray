package cmd

// 参数解析
func FlagsParse(args []string, keys map[string]string) map[string]string {
	resultMap := make(map[string]string)
	key := "data"
	for _, x := range args {
		if len(x) >= 2 {
			if x[:2] == "--" {
				key = x[2:]
				resultMap[key] = ""
			} else if x[:1] == "-" {
				d, ok := keys[x[1:]]
				if ok {
					key = d
				} else {
					key = x[1:]
				}
				resultMap[key] = ""
			} else {
				resultMap[key] = x
			}
		} else if len(x) > 0 {
			resultMap[key] = x
		}
	}
	return resultMap
}
