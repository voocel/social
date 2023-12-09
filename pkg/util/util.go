package util

func InSlice(s []int32, v int32) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}
