package util

func ContainsInt(list []int, ele int) bool {
	for _, element := range list {
		if element == ele {
			return true
		}
	}
	return false
}
